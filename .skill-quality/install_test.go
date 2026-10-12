// Copyright 2026 John Deatherage (@coversheets). SPDX-License-Identifier: MIT
package main

import (
	"encoding/json"
	"errors"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

func TestInstallAndUpdate(t *testing.T) {
	root := fixture(t)
	before, err := os.ReadFile(filepath.Join(root, ".skill-quality.json"))
	if err != nil {
		t.Fatal(err)
	}
	if err := installBundle(t.Context(), root, false); err != nil {
		t.Fatal(err)
	}
	if err := installBundle(t.Context(), root, false); err == nil {
		t.Fatal("overwrote without update")
	}
	if err := installBundle(t.Context(), root, true); err != nil {
		t.Fatal(err)
	}
	after, err := os.ReadFile(filepath.Join(root, ".skill-quality.json"))
	if err != nil || string(before) != string(after) {
		t.Fatal("configuration changed")
	}
	var m manifest
	b, err := os.ReadFile(filepath.Join(root, ".skill-quality/bundle.json"))
	if err != nil {
		t.Fatal(err)
	}
	if err := json.Unmarshal(b, &m); err != nil {
		t.Fatal(err)
	}
	if m.License != "MIT" || m.Source != "tpscloud/skillforge-ci" || len(m.Files) == 0 {
		t.Fatalf("%+v", m)
	}
	if r := validateRepository(t.Context(), root); r.Status != "pass" {
		t.Fatalf("%+v", r)
	}
	// The copied sources must compile without access to the producer checkout.
	cmd := exec.CommandContext(t.Context(), "go", "build", "-mod=readonly", "-o", filepath.Join(t.TempDir(), "skillforge"), ".")
	cmd.Dir = filepath.Join(root, ".skill-quality")
	if out, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("isolated build: %v\n%s", err, out)
	}
}

func TestInstallPreservesUnmanagedAndModifiedFiles(t *testing.T) {
	for _, name := range []string{".skill-quality/extra.txt", ".skill-quality/main.go", ".github/workflows/skill-quality.yml"} {
		t.Run(name, func(t *testing.T) {
			root := fixture(t)
			if err := installBundle(t.Context(), root, false); err != nil {
				t.Fatal(err)
			}
			writeFixture(t, root, name, "local work")
			if err := installBundle(t.Context(), root, true); err == nil {
				t.Fatal("overwrote local work")
			}
			b, err := os.ReadFile(filepath.Join(root, name))
			if err != nil || string(b) != "local work" {
				t.Fatal("lost local work")
			}
		})
	}
}

func TestInstallRejectsSymlinkDestination(t *testing.T) {
	root := fixture(t)
	outside := t.TempDir()
	if err := os.Symlink(outside, filepath.Join(root, ".github")); err != nil {
		t.Fatal(err)
	}
	if err := installBundle(t.Context(), root, true); err == nil {
		t.Fatal("symlink destination accepted")
	}
	entries, err := os.ReadDir(outside)
	if err != nil || len(entries) != 0 {
		t.Fatal("wrote outside repository")
	}
}

func TestLegacyMigration(t *testing.T) {
	root := fixture(t)
	var files []manifestFile
	for _, name := range []string{"validate.py", "run_checks.py", "scan_secrets.py", "skills.sh.schema.json", "requirements.txt"} {
		writeFixture(t, root, ".skill-quality/"+name, "legacy fixture")
		files = append(files, manifestFile{name, hashBytes([]byte("legacy fixture"))})
	}
	writeFixture(t, root, ".github/workflows/skill-quality.yml", "legacy workflow")
	files = append(files, manifestFile{".github/workflows/skill-quality.yml", hashBytes([]byte("legacy workflow"))})
	if err := writeJSON(filepath.Join(root, ".skill-quality/bundle.json"), manifest{Version: 1, Source: "coversheets/one-skill-to-rule-them-all-skill", Files: files}); err != nil {
		t.Fatal(err)
	}
	if err := installBundle(t.Context(), root, true); err != nil {
		t.Fatal(err)
	}
	for _, name := range []string{"validate.py", "requirements.txt", "skills.sh.schema.json"} {
		if _, err := os.Stat(filepath.Join(root, ".skill-quality", name)); !os.IsNotExist(err) {
			t.Fatalf("legacy file retained: %s", name)
		}
	}
}

func TestReplacementRollback(t *testing.T) {
	for _, rollbackFails := range []bool{false, true} {
		t.Run(map[bool]string{false: "restore", true: "preserve backup"}[rollbackFails], func(t *testing.T) {
			root := t.TempDir()
			recovery := filepath.Join(root, "recovery")
			if err := os.Mkdir(recovery, 0700); err != nil {
				t.Fatal(err)
			}
			target := filepath.Join(root, "target")
			staged := filepath.Join(recovery, "staged")
			backup := filepath.Join(recovery, "original")
			writeFixture(t, root, "target", "original bytes")
			writeFixture(t, recovery, "staged", "replacement")
			rename := func(from, to string) error {
				if from == staged || rollbackFails && from == backup {
					return errors.New("injected rename failure")
				}
				return os.Rename(from, to)
			}
			err := replaceFiles([]replacement{{Target: target, Staged: staged, Backup: backup}}, recovery, rename)
			if err == nil {
				t.Fatal("injected failure passed")
			}
			path := target
			if rollbackFails {
				path = backup
				if !strings.Contains(err.Error(), recovery) {
					t.Fatal("recovery location missing")
				}
			}
			b, readErr := os.ReadFile(path)
			if readErr != nil || string(b) != "original bytes" {
				t.Fatalf("lost original: %s %v", b, readErr)
			}
		})
	}
}

func TestInstallDiscoversSkills(t *testing.T) {
	root := fixture(t)
	if err := os.Remove(filepath.Join(root, ".skill-quality.json")); err != nil {
		t.Fatal(err)
	}
	if err := installBundle(t.Context(), root, false); err != nil {
		t.Fatal(err)
	}
	c, err := readConfig(root)
	if err != nil || len(c.Skills) != 1 || c.Skills[0] != "SKILL.md" {
		t.Fatalf("%+v %v", c, err)
	}
}

func TestRebuiltConsumerRetainsOriginAndDetectsEdits(t *testing.T) {
	root := fixture(t)
	if err := installBundle(t.Context(), root, false); err != nil {
		t.Fatal(err)
	}
	toolkit := filepath.Join(root, ".skill-quality")
	// Represent a copy from a clean producer, then edit its source before building.
	if err := writeJSON(filepath.Join(toolkit, "source.json"), sourceProvenance{Revision: "producer-revision", Dirty: false, ContentHash: bundleHash()}); err != nil {
		t.Fatal(err)
	}
	writeFixture(t, toolkit, "NOTICE.md", "Modified local notice.\n")
	binary := filepath.Join(t.TempDir(), "skillforge")
	build := exec.CommandContext(t.Context(), "go", "build", "-mod=readonly", "-o", binary, ".")
	build.Dir = toolkit
	if out, err := build.CombinedOutput(); err != nil {
		t.Fatalf("%s: %v", out, err)
	}
	target := fixture(t)
	if out, err := exec.CommandContext(t.Context(), binary, "install", target).CombinedOutput(); err != nil {
		t.Fatalf("%s: %v", out, err)
	}
	var m manifest
	b, err := os.ReadFile(filepath.Join(target, ".skill-quality/bundle.json"))
	if err != nil {
		t.Fatal(err)
	}
	if err := json.Unmarshal(b, &m); err != nil {
		t.Fatal(err)
	}
	if m.Revision != "producer-revision" || !m.Dirty || m.ContentHash == bundleHash() {
		t.Fatalf("modified consumer provenance: %+v", m)
	}
}

func TestInstallPreservesUnmanagedEmptyDirectory(t *testing.T) {
	root := fixture(t)
	if err := installBundle(t.Context(), root, false); err != nil {
		t.Fatal(err)
	}
	path := filepath.Join(root, ".skill-quality", "licenses", "local-empty")
	if err := os.Mkdir(path, 0755); err != nil {
		t.Fatal(err)
	}
	if err := installBundle(t.Context(), root, true); err == nil {
		t.Fatal("accepted unmanaged directory")
	}
	if info, err := os.Stat(path); err != nil || !info.IsDir() {
		t.Fatal("lost local directory")
	}
}
