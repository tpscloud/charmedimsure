// Copyright 2026 John Deatherage (@coversheets). SPDX-License-Identifier: MIT
package main

import (
	"archive/tar"
	"bytes"
	"compress/gzip"
	"context"
	"encoding/json"
	"errors"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"slices"
	"strconv"
	"testing"
)

func scannerArchive(t *testing.T, kind byte) []byte {
	t.Helper()
	var b bytes.Buffer
	gz := gzip.NewWriter(&b)
	tarWriter := tar.NewWriter(gz)
	data := []byte("fixture, never executed")
	header := &tar.Header{Name: "gitleaks", Typeflag: kind, Mode: 0700}
	if kind == tar.TypeReg {
		header.Size = int64(len(data))
	}
	if err := tarWriter.WriteHeader(header); err != nil {
		t.Fatal(err)
	}
	if kind == tar.TypeReg {
		if _, err := tarWriter.Write(data); err != nil {
			t.Fatal(err)
		}
	}
	if err := tarWriter.Close(); err != nil {
		t.Fatal(err)
	}
	if err := gz.Close(); err != nil {
		t.Fatal(err)
	}
	return b.Bytes()
}

func TestScannerArchive(t *testing.T) {
	b := scannerArchive(t, tar.TypeReg)
	if _, err := extractScanner(b, "bad"); err == nil {
		t.Fatal("bad digest accepted")
	}
	if got, err := extractScanner(b, hashBytes(b)); err != nil || string(got) != "fixture, never executed" {
		t.Fatalf("%q %v", got, err)
	}
	b = scannerArchive(t, tar.TypeSymlink)
	if _, err := extractScanner(b, hashBytes(b)); err == nil {
		t.Fatal("symlink archive executable accepted")
	}
}

func TestScannerModesAndEvidence(t *testing.T) {
	cases := []struct {
		name    string
		codes   [2]int
		reports [2]string
		want    int
	}{
		{"clean", [2]int{0, 0}, [2]string{"[]", "[]"}, 0},
		{"history finding", [2]int{10, 0}, [2]string{`[{"Secret":"REDACTED"}]`, "[]"}, 1},
		{"file finding", [2]int{0, 10}, [2]string{"[]", `[{}]`}, 1},
		{"tool failure", [2]int{0, 2}, [2]string{"[]", "[]"}, 2},
		{"missing evidence", [2]int{0, 0}, [2]string{"", "[]"}, 2},
		{"malformed evidence", [2]int{0, 0}, [2]string{"{", "[]"}, 2},
		{"null evidence", [2]int{0, 0}, [2]string{"null", "[]"}, 2},
		{"contradictory clean", [2]int{0, 0}, [2]string{"[{}]", "[]"}, 2},
		{"contradictory findings", [2]int{10, 0}, [2]string{"[]", "[]"}, 2},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			output := t.TempDir()
			calls := 0
			run := func(ctx context.Context, cwd string, args []string, out, diagnostic io.Writer) error {
				for _, flag := range []string{"--config", "--gitleaks-ignore-path"} {
					index := slices.Index(args, flag)
					if index < 0 || index+1 >= len(args) {
						t.Fatalf("missing scanner policy flag %s", flag)
					}
					path := args[index+1]
					if contained(cwd, path) {
						t.Fatalf("scanner policy inside target: %s", path)
					}
					content, err := os.ReadFile(path)
					if err != nil {
						t.Fatal(err)
					}
					if flag == "--config" && string(content) != scannerPolicy || flag == "--gitleaks-ignore-path" && len(content) != 0 {
						t.Fatalf("unexpected policy content for %s", flag)
					}
				}
				if !slices.Contains(args, "--ignore-gitleaks-allow") {
					t.Fatal("inline suppressions enabled")
				}
				i := calls
				calls++
				if tc.reports[i] != "" {
					writeFixture(t, output, "gitleaks-"+args[1]+".json", tc.reports[i])
				}
				if tc.codes[i] == 0 {
					return nil
				}
				return exec.CommandContext(ctx, "sh", "-c", "exit "+strconv.Itoa(tc.codes[i])).Run()
			}
			modes, err := scannerResults(t.Context(), t.TempDir(), output, "fixture", run, io.Discard, io.Discard)
			if calls != 2 || len(modes) != 2 || exitCode(err) != tc.want {
				t.Fatalf("calls=%d modes=%+v err=%v", calls, modes, err)
			}
		})
	}
}

func TestScannerOutputBoundary(t *testing.T) {
	root := fixture(t)
	err := scanRepository(t.Context(), root, filepath.Join(root, "reports"), io.Discard, io.Discard)
	if err == nil || errors.Is(err, errFindings) {
		t.Fatal("accepted output inside repository")
	}
	output := t.TempDir()
	if err := scanRepository(t.Context(), root, output, io.Discard, io.Discard); err == nil {
		t.Fatal("accepted stale evidence directory")
	}
	link := filepath.Join(t.TempDir(), "link")
	if err := os.Symlink(root, link); err != nil {
		t.Fatal(err)
	}
	if err := scanRepository(t.Context(), root, filepath.Join(link, "reports"), io.Discard, io.Discard); err == nil {
		t.Fatal("accepted symlink escape")
	}
}

func TestScannerRefusesTargetIgnore(t *testing.T) {
	for _, kind := range []string{"file", "symlink", "directory"} {
		t.Run(kind, func(t *testing.T) {
			root := t.TempDir()
			path := filepath.Join(root, ".gitleaksignore")
			var err error
			switch kind {
			case "file":
				err = os.WriteFile(path, []byte("fixture fingerprint\n"), 0600)
			case "symlink":
				err = os.Symlink(filepath.Join(root, "absent"), path)
			case "directory":
				err = os.Mkdir(path, 0700)
			}
			if err != nil {
				t.Fatal(err)
			}
			run := func(context.Context, string, []string, io.Writer, io.Writer) error {
				t.Fatal("scanner ran with target ignore file")
				return nil
			}
			_, err = scannerResults(t.Context(), root, t.TempDir(), "fixture", run, io.Discard, io.Discard)
			if exitCode(err) != 2 {
				t.Fatalf("expected operational refusal, got %v", err)
			}
			if _, err := os.Lstat(path); err != nil {
				t.Fatalf("target was changed: %v", err)
			}
		})
	}
}

func TestShallowScanFailsWithEvidence(t *testing.T) {
	root := fixture(t)
	if _, err := git(t.Context(), root, "-c", "user.name=Fixture", "-c", "user.email=fixture@example.invalid", "-c", "commit.gpgsign=false", "commit", "--allow-empty", "-m", "fixture"); err != nil {
		t.Fatal(err)
	}
	head, err := git(t.Context(), root, "rev-parse", "HEAD")
	if err != nil {
		t.Fatal(err)
	}
	writeFixture(t, root, ".git/shallow", string(head))
	output := filepath.Join(t.TempDir(), "scan")
	if err := scanRepository(t.Context(), root, output, io.Discard, io.Discard); exitCode(err) != 2 {
		t.Fatalf("shallow scan: %v", err)
	}
	raw, err := os.ReadFile(filepath.Join(output, "summary.json"))
	if err != nil {
		t.Fatal(err)
	}
	var report scanReport
	if err := json.Unmarshal(raw, &report); err != nil {
		t.Fatal(err)
	}
	if report.Status != "error" || len(report.Modes) != 0 {
		t.Fatalf("misleading evidence: %+v", report)
	}
}
