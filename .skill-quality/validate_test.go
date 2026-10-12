// Copyright 2026 John Deatherage (@coversheets). SPDX-License-Identifier: MIT
package main

import (
	"bytes"
	"context"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

func writeFixture(t *testing.T, root, name, body string) {
	t.Helper()
	path := filepath.Join(root, name)
	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path, []byte(body), 0644); err != nil {
		t.Fatal(err)
	}
}

func fixture(t *testing.T) string {
	t.Helper()
	root, err := canonical(t.TempDir())
	if err != nil {
		t.Fatal(err)
	}
	if _, err := git(t.Context(), root, "init", "-q"); err != nil {
		t.Fatal(err)
	}
	writeFixture(t, root, "SKILL.md", "---\nname: sample\ndescription: Test a sample.\n---\nRead [reference](guide.md).\n")
	writeFixture(t, root, "guide.md", "Reference\n")
	writeFixture(t, root, ".skill-quality.json", `{"version":1,"skills":["SKILL.md"],"checks":[]}`)
	writeFixture(t, root, "skills.sh.json", `{"groupings":[{"title":"Samples","skills":["sample"]}]}`)
	return root
}

func hasFinding(r report, code string) bool {
	for _, f := range r.Findings {
		if f.Code == code {
			return true
		}
	}
	return false
}

func TestValidation(t *testing.T) {
	cases := []struct{ name, path, body, code string }{
		{"broken link", "SKILL.md", "---\nname: sample\ndescription: Example\n---\n[missing](absent.md)", "missing-link"},
		{"reference link", "SKILL.md", "---\nname: sample\ndescription: Example\n---\n[missing][a]\n\n[a]: absent.md", "missing-link"},
		{"escaping link", "SKILL.md", "---\nname: sample\ndescription: Example\n---\n[escape](../absent.md)", "link-escape"},
		{"duplicate YAML", "SKILL.md", "---\nname: sample\nname: other\ndescription: Example\n---\n", "frontmatter"},
		{"YAML sequence", "SKILL.md", "---\n- sample\n---\n", "frontmatter"},
		{"numeric description", "SKILL.md", "---\nname: sample\ndescription: 42\n---\n", "frontmatter"},
		{"unknown group", "skills.sh.json", `{"groupings":[{"title":"Samples","skills":["absent"]}]}`, "unknown-grouping"},
		{"duplicate group", "skills.sh.json", `{"groupings":[{"title":"Samples","skills":["sample","sample"]}]}`, "duplicate-grouping"},
		{"bad JSON", "skills.sh.json", "{", "page-schema"},
		{"bad shape", "skills.sh.json", "[]", "page-schema"},
		{"null group", "skills.sh.json", `{"groupings":[null]}`, "page-schema"},
		{"null names", "skills.sh.json", `{"groupings":[{"title":"Samples","skills":[null]}]}`, "page-schema"},
		{"unknown page key", "skills.sh.json", `{"extra":true,"groupings":[{"title":"Samples","skills":["sample"]}]}`, "page-schema"},
		{"null title", "skills.sh.json", `{"groupings":[{"title":null,"skills":["sample"]}]}`, "page-schema"},
		{"duplicate JSON", ".skill-quality.json", `{"skills":["SKILL.md"],"skills":[]}`, "configuration"},
		{"duplicate paths", ".skill-quality.json", `{"skills":["SKILL.md","SKILL.md"]}`, "configuration"},
		{"path escape", ".skill-quality.json", `{"skills":["../SKILL.md"]}`, "path-escape"},
		{"check escape", ".skill-quality.json", `{"skills":["SKILL.md"],"checks":[{"name":"bad","command":["true"],"cwd":".."}]}`, "configuration"},
		{"shell string", ".skill-quality.json", `{"skills":["SKILL.md"],"checks":[{"name":"bad","command":"true"}]}`, "configuration"},
		{"null argument", ".skill-quality.json", `{"skills":["SKILL.md"],"checks":[{"name":"bad","command":["echo",null]}]}`, "configuration"},
		{"null checks", ".skill-quality.json", `{"skills":["SKILL.md"],"checks":null}`, "configuration"},
		{"null cwd", ".skill-quality.json", `{"skills":["SKILL.md"],"checks":[{"name":"bad","command":["true"],"cwd":null}]}`, "configuration"},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			root := fixture(t)
			writeFixture(t, root, tc.path, tc.body)
			r := validateRepository(t.Context(), root)
			if r.Status != "fail" || !hasFinding(r, tc.code) {
				t.Fatalf("wanted %s; got %+v", tc.code, r)
			}
		})
	}
}

func TestCleanContentAndCodeExamples(t *testing.T) {
	root := fixture(t)
	r := validateRepository(t.Context(), root)
	if r.Status != "pass" || len(r.NotAssessed) == 0 || r.ContentHash == "" {
		t.Fatalf("%+v", r)
	}
	writeFixture(t, root, "guide.md", "Changed\n")
	if validateRepository(t.Context(), root).ContentHash == r.ContentHash {
		t.Fatal("hash ignored uncommitted edit")
	}
	writeFixture(t, root, "SKILL.md", "---\nname: sample\ndescription: Example\n---\n```md\n[x](absent.md)\n```\n`[x](absent.md)`\n")
	if r := validateRepository(t.Context(), root); r.Status != "pass" {
		t.Fatalf("example links checked: %+v", r)
	}
}

func TestOptionalFrontmatterPolicy(t *testing.T) {
	cases := []struct {
		name, directory, policy, fields, code string
	}{
		{"default permits extensions and different directory", "other", "", "custom: value\n", ""},
		{"allowed extension and matching directory", "sample", `,"frontmatter":{"allowedKeys":["name","description","license","argument-hint"],"requireDirectoryName":true}`, "license: MIT\nargument-hint: '[topic]'\n", ""},
		{"unexpected key", "sample", `,"frontmatter":{"allowedKeys":["name","description"]}`, "license: MIT\n", "frontmatter"},
		{"directory mismatch", "other", `,"frontmatter":{"requireDirectoryName":true}`, "", "skill-directory-name"},
		{"directory check disabled", "other", `,"frontmatter":{"requireDirectoryName":false}`, "", ""},
		{"empty policy", "other", `,"frontmatter":{}`, "custom: value\n", ""},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			root := fixture(t)
			path := "skills/" + tc.directory + "/SKILL.md"
			writeFixture(t, root, path, "---\nname: sample\ndescription: Test a sample.\n"+tc.fields+"---\nBody.\n")
			writeFixture(t, root, ".skill-quality.json", `{"version":1,"skills":["`+path+`"],"checks":[]`+tc.policy+`}`)
			r := validateRepository(t.Context(), root)
			if tc.code == "" {
				if r.Status != "pass" {
					t.Fatalf("unexpected failure: %+v", r.Findings)
				}
			} else if r.Status != "fail" || !hasFinding(r, tc.code) {
				t.Fatalf("wanted %s; got %+v", tc.code, r.Findings)
			}
		})
	}
}

func TestInvalidFrontmatterPolicy(t *testing.T) {
	for _, policy := range []string{
		`null`, `[]`, `true`,
		`{"unknown":true}`, `{"requireDirectoryName":null}`, `{"requireDirectoryName":"yes"}`,
		`{"allowedKeys":null}`, `{"allowedKeys":[]}`, `{"allowedKeys":"name"}`,
		`{"allowedKeys":[null]}`, `{"allowedKeys":[""]}`, `{"allowedKeys":[" "]}`,
		`{"allowedKeys":["name","name"]}`, `{"allowedKeys":[42]}`,
		`{"requireDirectoryName":true,"requireDirectoryName":false}`,
	} {
		t.Run(policy, func(t *testing.T) {
			root := fixture(t)
			writeFixture(t, root, ".skill-quality.json", `{"skills":["SKILL.md"],"frontmatter":`+policy+`}`)
			if r := validateRepository(t.Context(), root); r.Status != "fail" || !hasFinding(r, "configuration") {
				t.Fatalf("invalid policy accepted: %+v", r)
			}
		})
	}
}

func TestFrontmatterPolicyAtRepositoryRoot(t *testing.T) {
	root := fixture(t)
	writeFixture(t, root, ".skill-quality.json", `{"skills":["SKILL.md"],"frontmatter":{"requireDirectoryName":true}}`)
	if r := validateRepository(t.Context(), root); !hasFinding(r, "skill-directory-name") {
		t.Fatalf("root entrypoint skipped directory check: %+v", r)
	}
}

func TestUnsafeSymlinksAndTrackedGeneratedFiles(t *testing.T) {
	root := fixture(t)
	if err := os.Symlink(t.TempDir(), filepath.Join(root, "escape")); err != nil {
		t.Fatal(err)
	}
	if _, err := git(t.Context(), root, "add", "escape"); err != nil {
		t.Fatal(err)
	}
	writeFixture(t, root, ".DS_Store", "generated")
	if _, err := git(t.Context(), root, "add", ".DS_Store"); err != nil {
		t.Fatal(err)
	}
	r := validateRepository(t.Context(), root)
	if !hasFinding(r, "unsafe-symlink") || !hasFinding(r, "generated-file") {
		t.Fatalf("%+v", r)
	}
	writeFixture(t, root, ".skill-quality.json", `{"skills":["SKILL.md"],"checks":[{"name":"escape","command":["true"],"cwd":"escape"}]}`)
	if _, err := readConfig(root); err == nil {
		t.Fatal("symlink cwd accepted")
	}
}

func TestValidateNeverExecutesChecks(t *testing.T) {
	root := fixture(t)
	writeFixture(t, root, ".skill-quality.json", `{"skills":["SKILL.md"],"checks":[{"name":"marker","command":["touch","executed"]}]}`)
	if r := validateRepository(t.Context(), root); r.Status != "pass" {
		t.Fatalf("%+v", r)
	}
	if _, err := os.Stat(filepath.Join(root, "executed")); !os.IsNotExist(err) {
		t.Fatal("validation ran check")
	}
}

func TestRunChecksLiteralArgumentsAndFailure(t *testing.T) {
	root := fixture(t)
	writeFixture(t, root, ".skill-quality.json", `{"skills":["SKILL.md"],"checks":[{"name":"literal","command":["printf","%s","$(touch executed)"]}]}`)
	var out bytes.Buffer
	if err := runChecks(t.Context(), root, &out, &out); err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(out.String(), "$(touch executed)") {
		t.Fatal(out.String())
	}
	if _, err := os.Stat(filepath.Join(root, "executed")); !os.IsNotExist(err) {
		t.Fatal("argument evaluated by shell")
	}
	writeFixture(t, root, ".skill-quality.json", `{"skills":["SKILL.md"],"checks":[{"name":"failure","command":["false"]},{"name":"marker","command":["touch","executed"]}]}`)
	if err := runChecks(t.Context(), root, &out, &out); err == nil {
		t.Fatal("failed check passed")
	}
	if _, err := os.Stat(filepath.Join(root, "executed")); !os.IsNotExist(err) {
		t.Fatal("continued after failure")
	}
}

func TestCommandCancellation(t *testing.T) {
	ctx, cancel := context.WithTimeout(t.Context(), 50*time.Millisecond)
	defer cancel()
	var out bytes.Buffer
	start := time.Now()
	err := execute(ctx, t.TempDir(), []string{"sh", "-c", "sleep 30 & wait"}, &out, &out)
	if err == nil || time.Since(start) > 4*time.Second {
		t.Fatalf("cancellation failed: %v", err)
	}
}

func TestCLI(t *testing.T) {
	for _, args := range [][]string{{"--help"}, {"--version"}} {
		var out, diagnostic bytes.Buffer
		cmd := newRoot(&out, &diagnostic)
		cmd.SetArgs(args)
		if err := cmd.ExecuteContext(t.Context()); err != nil || out.Len() == 0 || diagnostic.Len() != 0 {
			t.Fatalf("help/version: %v", err)
		}
	}
	root := fixture(t)
	output := filepath.Join(t.TempDir(), "report.json")
	var out bytes.Buffer
	cmd := newRoot(&out, &out)
	cmd.SetArgs([]string{"validate", root, "--output", output})
	if err := cmd.ExecuteContext(t.Context()); err != nil {
		t.Fatal(err)
	}
	if _, err := os.Stat(output); err != nil {
		t.Fatal(err)
	}
	cmd = newRoot(&out, &out)
	cmd.SetArgs([]string{"install"})
	if err := cmd.ExecuteContext(t.Context()); err == nil {
		t.Fatal("missing required argument accepted")
	}
	cmd = newRoot(&out, &out)
	cmd.SetArgs([]string{"scan", root})
	if err := cmd.ExecuteContext(t.Context()); err == nil {
		t.Fatal("missing output accepted")
	}
}

func TestNonRepositoryFailsWithEvidence(t *testing.T) {
	r := validateRepository(t.Context(), t.TempDir())
	if r.Status != "fail" || !hasFinding(r, "repository") {
		t.Fatalf("%+v", r)
	}
}

func TestExitCodes(t *testing.T) {
	if exitCode(nil) != 0 || exitCode(errFindings) != 1 || exitCode(exec.ErrNotFound) != 2 {
		t.Fatal("exit code contract")
	}
}

func TestUnknownConfigurationKeysFailClosed(t *testing.T) {
	for _, raw := range []string{
		`{"skills":["SKILL.md"],"chekcs":[{"name":"test","command":["false"]}]}`,
		`{"skills":["SKILL.md"],"Checks":[]}`,
		`{"skills":["SKILL.md"],"checks":[{"name":"test","command":["true"],"cdw":"elsewhere"}]}`,
	} {
		root := fixture(t)
		writeFixture(t, root, ".skill-quality.json", raw)
		if r := validateRepository(t.Context(), root); r.Status != "fail" || !hasFinding(r, "configuration") {
			t.Fatalf("accepted unknown key: %+v", r)
		}
		if err := runChecks(t.Context(), root, &bytes.Buffer{}, &bytes.Buffer{}); err == nil {
			t.Fatal("check accepted unknown key")
		}
	}
}

func TestContentHashRecords(t *testing.T) {
	for _, kind := range []string{"missing", "boundaries", "mode", "symlink"} {
		t.Run(kind, func(t *testing.T) {
			root := fixture(t)
			beforeContent := "missing"
			if kind == "boundaries" {
				beforeContent = "left\x00b\x00right"
			}
			if kind == "symlink" {
				beforeContent = "symlink\x00guide.md"
			}
			writeFixture(t, root, "a", beforeContent)
			if _, err := git(t.Context(), root, "add", "a"); err != nil {
				t.Fatal(err)
			}
			before, err := contentHash(t.Context(), root)
			if err != nil {
				t.Fatal(err)
			}
			switch kind {
			case "missing":
				if err := os.Remove(filepath.Join(root, "a")); err != nil {
					t.Fatal(err)
				}
			case "boundaries":
				writeFixture(t, root, "a", "left")
				writeFixture(t, root, "b", "right")
			case "mode":
				if err := os.Chmod(filepath.Join(root, "a"), 0755); err != nil {
					t.Fatal(err)
				}
			case "symlink":
				if err := os.Remove(filepath.Join(root, "a")); err != nil {
					t.Fatal(err)
				}
				if err := os.Symlink("guide.md", filepath.Join(root, "a")); err != nil {
					t.Fatal(err)
				}
			}
			after, err := contentHash(t.Context(), root)
			if err != nil {
				t.Fatal(err)
			}
			if before == after {
				t.Fatal("different content hashed identically")
			}
			if r := validateRepository(t.Context(), root); r.ContentHashVersion != 2 {
				t.Fatal("hash version absent")
			}
		})
	}
}
