// Copyright 2026 John Deatherage (@coversheets). SPDX-License-Identifier: MIT
package main

import (
	"context"
	"crypto/sha256"
	"embed"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"slices"
	"strings"
)

//go:embed *.go go.mod go.sum LICENSE NOTICE.md licenses templates source.json
var bundle embed.FS

type manifestFile struct {
	Path   string `json:"path"`
	SHA256 string `json:"sha256"`
}
type manifest struct {
	Version     int            `json:"version"`
	Source      string         `json:"source"`
	Revision    string         `json:"revision"`
	Dirty       bool           `json:"sourceDirty"`
	License     string         `json:"license"`
	ContentHash string         `json:"contentHash"`
	Files       []manifestFile `json:"files"`
}

func bundleFiles() []string {
	var files []string
	_ = fs.WalkDir(bundle, ".", func(path string, entry fs.DirEntry, err error) error {
		if err == nil && !entry.IsDir() {
			files = append(files, path)
		}
		return err
	})
	slices.Sort(files)
	return files
}

func bundleHash() string {
	h := sha256.New()
	for _, path := range bundleFiles() {
		if path == "source.json" {
			continue
		}
		b, _ := bundle.ReadFile(path)
		io.WriteString(h, path+"\x00")
		h.Write(b)
		h.Write([]byte{0})
	}
	return hex.EncodeToString(h.Sum(nil))
}

type sourceProvenance struct {
	Revision    string `json:"revision"`
	Dirty       bool   `json:"dirty"`
	ContentHash string `json:"contentHash"`
}

func producerSource() sourceProvenance {
	b, _ := bundle.ReadFile("source.json")
	var source sourceProvenance
	_ = json.Unmarshal(b, &source)
	currentHash := bundleHash()
	if source.ContentHash != "" {
		source.Dirty = source.Dirty || source.ContentHash != currentHash
		source.ContentHash = currentHash
		return source
	}
	source.Revision = buildRevision()
	source.Dirty = buildDirty()
	source.ContentHash = currentHash
	return source
}

func noSymlinks(root, relative string) error {
	if !filepath.IsLocal(relative) {
		return fmt.Errorf("unsafe destination: %s", relative)
	}
	path := root
	for component := range strings.SplitSeq(filepath.ToSlash(relative), "/") {
		path = filepath.Join(path, component)
		info, err := os.Lstat(path)
		if os.IsNotExist(err) {
			continue
		}
		if err != nil {
			return err
		}
		if info.Mode()&os.ModeSymlink != 0 {
			return fmt.Errorf("refusing symlink destination: %s", relative)
		}
	}
	return nil
}

func checkExisting(root string, update bool) error {
	for _, relative := range []string{".skill-quality", ".skill-quality/bundle.json", ".github/workflows/skill-quality.yml", ".skill-quality.json"} {
		if err := noSymlinks(root, relative); err != nil {
			return err
		}
	}
	_, bundleErr := os.Stat(filepath.Join(root, ".skill-quality"))
	_, workflowErr := os.Stat(filepath.Join(root, ".github/workflows/skill-quality.yml"))
	if os.IsNotExist(bundleErr) && os.IsNotExist(workflowErr) {
		return nil
	}
	if !update {
		return fmt.Errorf("bundle or workflow exists; inspect before using --update")
	}
	var old manifest
	b, err := os.ReadFile(filepath.Join(root, ".skill-quality/bundle.json"))
	if err != nil {
		return err
	}
	if err := uniqueJSON(b, &old); err != nil {
		return err
	}
	if !((old.Version == 2 && old.Source == "tpscloud/skillforge-ci") || (old.Version == 1 && old.Source == "coversheets/one-skill-to-rule-them-all-skill")) || len(old.Files) == 0 {
		return fmt.Errorf("unrecognized bundle; refusing replacement")
	}
	expected := map[string]bool{"bundle.json": true}
	expectedDirs := map[string]bool{".": true}
	workflowFound := false
	for _, f := range old.Files {
		if !filepath.IsLocal(f.Path) {
			return fmt.Errorf("unsafe manifest path: %s", f.Path)
		}
		relative := filepath.Join(".skill-quality", f.Path)
		if f.Path == ".github/workflows/skill-quality.yml" {
			relative = f.Path
			workflowFound = true
		} else {
			expected[f.Path] = true
			for dir := filepath.Dir(f.Path); dir != "."; dir = filepath.Dir(dir) {
				expectedDirs[filepath.ToSlash(dir)] = true
			}
		}
		if err := noSymlinks(root, relative); err != nil {
			return err
		}
		b, err := os.ReadFile(filepath.Join(root, relative))
		if err != nil {
			return err
		}
		if hashBytes(b) != f.SHA256 {
			return fmt.Errorf("managed file was modified: %s; preserve or reconcile it before updating", relative)
		}
	}
	if !workflowFound {
		return fmt.Errorf("bundle does not own the existing workflow")
	}
	return filepath.WalkDir(filepath.Join(root, ".skill-quality"), func(path string, entry fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		relative, err := filepath.Rel(filepath.Join(root, ".skill-quality"), path)
		if err != nil {
			return err
		}
		if entry.IsDir() {
			if !expectedDirs[filepath.ToSlash(relative)] {
				return fmt.Errorf("unmanaged directory in bundle: %s", relative)
			}
			return nil
		}
		if !expected[filepath.ToSlash(relative)] {
			return fmt.Errorf("unmanaged file in bundle: %s", relative)
		}
		return nil
	})
}

type replacement struct {
	Target      string
	Staged      string
	Backup      string
	HadOriginal bool
	Installed   bool
}

// The recovery directory survives any rollback failure. A caller can recover
// its original files from the named location instead of losing them to cleanup.
func replaceFiles(changes []replacement, recovery string, rename func(string, string) error) error {
	var failed error
	for i := range changes {
		change := &changes[i]
		if err := os.MkdirAll(filepath.Dir(change.Target), 0755); err != nil {
			failed = err
			break
		}
		if _, err := os.Lstat(change.Target); err == nil {
			if err := rename(change.Target, change.Backup); err != nil {
				failed = err
				break
			}
			change.HadOriginal = true
		} else if !os.IsNotExist(err) {
			failed = err
			break
		}
		if err := rename(change.Staged, change.Target); err != nil {
			failed = err
			break
		}
		change.Installed = true
	}
	if failed == nil {
		return os.RemoveAll(recovery)
	}
	var rollbackErrors []error
	for i := len(changes) - 1; i >= 0; i-- {
		change := changes[i]
		if change.Installed {
			if err := os.RemoveAll(change.Target); err != nil {
				rollbackErrors = append(rollbackErrors, err)
				continue
			}
		}
		if change.HadOriginal {
			if err := rename(change.Backup, change.Target); err != nil {
				rollbackErrors = append(rollbackErrors, err)
			}
		}
	}
	if len(rollbackErrors) != 0 {
		return fmt.Errorf("installation failed; recover original files from %s: %w", recovery, errors.Join(append(rollbackErrors, failed)...))
	}
	return errors.Join(failed, os.RemoveAll(recovery))
}

func installBundle(ctx context.Context, path string, update bool) error {
	root, err := repositoryRoot(ctx, path)
	if err != nil {
		return err
	}
	if err := checkExisting(root, update); err != nil {
		return err
	}
	var configBytes []byte
	if _, err := os.Stat(filepath.Join(root, ".skill-quality.json")); os.IsNotExist(err) {
		b, err := git(ctx, root, "ls-files", "--cached", "--others", "--exclude-standard", "-z")
		if err != nil {
			return err
		}
		var skills []string
		for path := range strings.SplitSeq(string(b), "\x00") {
			if path == "SKILL.md" || strings.HasSuffix(path, "/SKILL.md") {
				skills = append(skills, path)
			}
		}
		slices.Sort(skills)
		skills = slices.Compact(skills)
		if len(skills) == 0 {
			return fmt.Errorf("no SKILL.md found; declare entrypoints before installing")
		}
		configBytes, err = json.MarshalIndent(configuration{Version: 1, Skills: skills, Checks: []check{}}, "", "  ")
		if err != nil {
			return err
		}
		configBytes = append(configBytes, '\n')
	} else if err != nil {
		return err
	}
	recovery, err := os.MkdirTemp(root, ".skillforge-recovery-")
	if err != nil {
		return err
	}
	staged := filepath.Join(recovery, "new-bundle")
	if err := os.Mkdir(staged, 0755); err != nil {
		return errors.Join(err, os.RemoveAll(recovery))
	}
	stagingFailed := func(err error) error { return errors.Join(err, os.RemoveAll(recovery)) }
	source := producerSource()
	m := manifest{Version: 2, Source: "tpscloud/skillforge-ci", Revision: source.Revision, Dirty: source.Dirty, License: "MIT", ContentHash: bundleHash(), Files: []manifestFile{}}
	for _, name := range bundleFiles() {
		b, err := bundle.ReadFile(name)
		if err != nil {
			return stagingFailed(err)
		}
		if name == "source.json" {
			b, err = json.Marshal(source)
			if err != nil {
				return stagingFailed(err)
			}
			b = append(b, '\n')
		}
		destination := filepath.Join(staged, filepath.FromSlash(name))
		if err := os.MkdirAll(filepath.Dir(destination), 0755); err != nil {
			return stagingFailed(err)
		}
		if err := os.WriteFile(destination, b, 0644); err != nil {
			return stagingFailed(err)
		}
		m.Files = append(m.Files, manifestFile{name, hashBytes(b)})
	}
	workflow, err := bundle.ReadFile("templates/skill-quality.yml")
	if err != nil {
		return stagingFailed(err)
	}
	workflowPath := filepath.Join(recovery, "new-workflow")
	if err := os.WriteFile(workflowPath, workflow, 0644); err != nil {
		return stagingFailed(err)
	}
	m.Files = append(m.Files, manifestFile{".github/workflows/skill-quality.yml", hashBytes(workflow)})
	if err := writeJSON(filepath.Join(staged, "bundle.json"), m); err != nil {
		return stagingFailed(err)
	}
	changes := []replacement{{Target: filepath.Join(root, ".skill-quality"), Staged: staged, Backup: filepath.Join(recovery, "original-bundle")}, {Target: filepath.Join(root, ".github/workflows/skill-quality.yml"), Staged: workflowPath, Backup: filepath.Join(recovery, "original-workflow")}}
	if configBytes != nil {
		configPath := filepath.Join(recovery, "new-config")
		if err := os.WriteFile(configPath, configBytes, 0644); err != nil {
			return stagingFailed(err)
		}
		changes = append(changes, replacement{Target: filepath.Join(root, ".skill-quality.json"), Staged: configPath, Backup: filepath.Join(recovery, "original-config")})
	}
	return replaceFiles(changes, recovery, os.Rename)
}
