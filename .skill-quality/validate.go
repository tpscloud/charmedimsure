// Copyright 2026 John Deatherage (@coversheets). SPDX-License-Identifier: MIT
package main

import (
	"bytes"
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"maps"
	"net/url"
	"os"
	"path/filepath"
	"regexp"
	"slices"
	"strings"
	"time"
	"unicode/utf8"

	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/text"
	"gopkg.in/yaml.v3"
)

type finding struct {
	Code     string `json:"code"`
	Path     string `json:"path"`
	Message  string `json:"message"`
	Severity string `json:"severity"`
}
type skill struct {
	Name   string `json:"name"`
	Path   string `json:"path"`
	SHA256 string `json:"sha256"`
}
type report struct {
	Version            int       `json:"version"`
	CheckedAt          string    `json:"checkedAt"`
	Revision           *string   `json:"revision"`
	ContentHash        string    `json:"contentHash"`
	ContentHashVersion int       `json:"contentHashVersion"`
	ValidatorHash      string    `json:"validatorHash"`
	Dirty              bool      `json:"dirty"`
	Status             string    `json:"status"`
	Scope              string    `json:"scope"`
	Skills             []skill   `json:"skills"`
	Findings           []finding `json:"findings"`
	NotAssessed        []string  `json:"notAssessed"`
}

func (r *report) issue(code, path, message string, warning bool) {
	severity := "error"
	if warning {
		severity = "warning"
	} else {
		r.Status = "fail"
	}
	r.Findings = append(r.Findings, finding{code, path, message, severity})
}

var namePattern = regexp.MustCompile(`^[a-z0-9]+(-[a-z0-9]+)*$`)

func frontmatter(raw []byte, policy frontmatterPolicy) (string, []byte, error) {
	if !utf8.Valid(raw) {
		return "", nil, fmt.Errorf("entrypoint is not UTF-8")
	}
	normalized := bytes.ReplaceAll(raw, []byte("\r\n"), []byte("\n"))
	if !bytes.HasPrefix(normalized, []byte("---\n")) {
		return "", nil, fmt.Errorf("missing YAML frontmatter")
	}
	lines := bytes.Split(normalized[4:], []byte("\n"))
	end := slices.IndexFunc(lines, func(line []byte) bool { return string(line) == "---" })
	if end < 0 {
		return "", nil, fmt.Errorf("unterminated YAML frontmatter")
	}
	var meta map[string]any
	if err := yaml.Unmarshal(bytes.Join(lines[:end], []byte("\n")), &meta); err != nil {
		return "", nil, err
	}
	name, ok := meta["name"].(string)
	if !ok || len(name) > 64 || !namePattern.MatchString(name) {
		return "", nil, fmt.Errorf("name must be 1-64 lowercase letters, digits, and single hyphens")
	}
	description, ok := meta["description"].(string)
	if !ok || strings.TrimSpace(description) == "" || utf8.RuneCountInString(description) > 1024 {
		return "", nil, fmt.Errorf("description must be a nonempty string of at most 1024 characters")
	}
	if len(policy.AllowedKeys) > 0 {
		var extra []string
		for _, key := range slices.Sorted(maps.Keys(meta)) {
			if !slices.Contains(policy.AllowedKeys, key) {
				extra = append(extra, key)
			}
		}
		if len(extra) > 0 {
			return "", nil, fmt.Errorf("unexpected frontmatter keys: %s", strings.Join(extra, ", "))
		}
	}
	return name, bytes.Join(lines[end+1:], []byte("\n")), nil
}

func validateRepository(ctx context.Context, path string) report {
	r := report{Version: 1, ContentHashVersion: 2, CheckedAt: time.Now().UTC().Format(time.RFC3339), Status: "pass", Dirty: true,
		Scope: "declared entrypoints, page config, tracked file hygiene", Skills: []skill{}, Findings: []finding{},
		ValidatorHash: bundleHash(), NotAssessed: []string{"semantic security", "behavioral quality", "provider audits", "redistribution rights", "installer compatibility"}}
	root, err := repositoryRoot(ctx, path)
	if err != nil {
		r.issue("repository", ".", err.Error(), false)
		return r
	}
	c, err := readConfig(root)
	if err != nil {
		r.issue("configuration", ".skill-quality.json", err.Error(), false)
	} else {
		for _, relative := range c.Skills {
			path := filepath.Join(root, relative)
			if !filepath.IsLocal(relative) || !contained(root, path) {
				r.issue("path-escape", relative, "skill path escapes repository", false)
				continue
			}
			raw, err := os.ReadFile(path)
			if err != nil {
				r.issue("frontmatter", relative, err.Error(), false)
				continue
			}
			name, body, err := frontmatter(raw, c.Frontmatter)
			if err != nil {
				r.issue("frontmatter", relative, err.Error(), false)
				continue
			}
			r.Skills = append(r.Skills, skill{name, relative, hashBytes(raw)})
			if dir := filepath.Base(filepath.Dir(path)); c.Frontmatter.RequireDirectoryName && name != dir {
				r.issue("skill-directory-name", relative, fmt.Sprintf("name %q does not match directory %q", name, dir), false)
			}
			doc := goldmark.DefaultParser().Parse(text.NewReader(body))
			_ = ast.Walk(doc, func(n ast.Node, entering bool) (ast.WalkStatus, error) {
				link, ok := n.(*ast.Link)
				if !ok || !entering {
					return ast.WalkContinue, nil
				}
				target := string(link.Destination)
				u, err := url.Parse(target)
				if err != nil {
					r.issue("invalid-link", relative, "invalid link: "+target, false)
					return ast.WalkContinue, nil
				}
				if u.Scheme != "" || u.Host != "" || u.Path == "" || strings.HasPrefix(u.Path, "/") {
					return ast.WalkContinue, nil
				}
				linked := filepath.Join(filepath.Dir(path), filepath.FromSlash(u.Path))
				if !contained(root, linked) {
					r.issue("link-escape", relative, "relative link leaves repository: "+target, false)
				} else if _, err := os.Stat(linked); err != nil {
					r.issue("missing-link", relative, "missing relative file: "+target, false)
				}
				return ast.WalkContinue, nil
			})
			if bytes.Count(body, []byte("\n")) > 500 {
				r.issue("entrypoint-length", relative, "entrypoint exceeds 500 lines; review progressive disclosure", true)
			}
		}
	}
	pagePath := filepath.Join(root, "skills.sh.json")
	var grouped []string
	if !contained(root, pagePath) {
		err = fmt.Errorf("page configuration escapes repository")
	} else {
		var raw []byte
		raw, err = os.ReadFile(pagePath)
		if err == nil {
			grouped, err = pageGroups(raw)
		}
	}
	if err != nil {
		r.issue("page-schema", "skills.sh.json", err.Error(), false)
	} else {
		names, seen := map[string]bool{}, map[string]bool{}
		for _, s := range r.Skills {
			names[s.Name] = true
		}
		for _, name := range grouped {
			if seen[name] {
				r.issue("duplicate-grouping", "skills.sh.json", "skill occurs in multiple groups: "+name, false)
			}
			seen[name] = true
			if !names[name] {
				r.issue("unknown-grouping", "skills.sh.json", "grouped skill is absent from inventory: "+name, false)
			}
		}
		for _, s := range r.Skills {
			if !seen[s.Name] {
				r.issue("ungrouped", "skills.sh.json", "skill is not explicitly grouped: "+s.Name, true)
			}
		}
	}
	tracked, err := git(ctx, root, "ls-files", "-z")
	if err != nil {
		r.issue("git", ".", err.Error(), false)
	} else {
		for relative := range strings.SplitSeq(string(tracked), "\x00") {
			if relative == "" {
				continue
			}
			path := filepath.Join(root, relative)
			info, err := os.Lstat(path)
			if err == nil && info.Mode()&os.ModeSymlink != 0 {
				if !contained(root, path) {
					r.issue("unsafe-symlink", relative, "tracked symlink is broken or leaves repository", false)
				}
			}
			if slices.Contains(strings.Split(relative, "/"), "__pycache__") || strings.HasSuffix(relative, ".pyc") || filepath.Base(relative) == ".DS_Store" {
				r.issue("generated-file", relative, "generated local file is tracked", false)
			}
		}
	}
	if b, err := git(ctx, root, "rev-parse", "HEAD"); err == nil {
		r.Revision = new(strings.TrimSpace(string(b)))
	}
	if b, err := git(ctx, root, "status", "--porcelain"); err == nil {
		r.Dirty = len(b) != 0
	} else {
		r.issue("git", ".", err.Error(), false)
	}
	r.ContentHash, err = contentHash(ctx, root)
	if err != nil {
		r.issue("content-hash", ".", err.Error(), false)
	}
	return r
}

// This is our implementation of the supported page fields, not a bundled copy
// of a third-party schema. JSON shapes and bounds are checked before grouping.
func pageGroups(raw []byte) ([]string, error) {
	var page map[string]json.RawMessage
	if err := uniqueJSON(raw, &page); err != nil {
		return nil, err
	}
	if page == nil {
		return nil, fmt.Errorf("page must be an object")
	}
	for k, v := range page {
		switch k {
		case "$schema", "schema", "notGrouped":
			var s string
			if bytes.Equal(v, []byte("null")) || json.Unmarshal(v, &s) != nil {
				return nil, fmt.Errorf("%s must be a string", k)
			}
			if k == "notGrouped" && s != "top" && s != "bottom" {
				return nil, fmt.Errorf("notGrouped must be top or bottom")
			}
		case "groupings":
		default:
			return nil, fmt.Errorf("unknown page property: %s", k)
		}
	}
	var groups []map[string]json.RawMessage
	if err := json.Unmarshal(page["groupings"], &groups); err != nil || len(groups) < 1 || len(groups) > 50 {
		return nil, fmt.Errorf("groupings must contain 1-50 objects")
	}
	var names []string
	for _, group := range groups {
		for key := range group {
			if key != "title" && key != "description" && key != "skills" {
				return nil, fmt.Errorf("unknown group property: %s", key)
			}
		}
		var title string
		if err := json.Unmarshal(group["title"], &title); err != nil || utf8.RuneCountInString(title) < 1 || utf8.RuneCountInString(title) > 120 {
			return nil, fmt.Errorf("group title must contain 1-120 characters")
		}
		if value, ok := group["description"]; ok {
			var description string
			if bytes.Equal(value, []byte("null")) || json.Unmarshal(value, &description) != nil || utf8.RuneCountInString(description) > 500 {
				return nil, fmt.Errorf("group description must be a string of at most 500 characters")
			}
		}
		var entries []string
		if err := json.Unmarshal(group["skills"], &entries); err != nil || len(entries) < 1 || len(entries) > 500 {
			return nil, fmt.Errorf("group skills must contain 1-500 strings")
		}
		for _, name := range entries {
			if utf8.RuneCountInString(name) < 1 || utf8.RuneCountInString(name) > 120 {
				return nil, fmt.Errorf("group skill names must contain 1-120 characters")
			}
		}
		names = append(names, entries...)
	}
	return names, nil
}

func contentHash(ctx context.Context, root string) (string, error) {
	b, err := git(ctx, root, "ls-files", "--cached", "--others", "--exclude-standard", "-z")
	if err != nil {
		return "", err
	}
	paths := strings.Split(strings.TrimSuffix(string(b), "\x00"), "\x00")
	slices.Sort(paths)
	paths = slices.Compact(paths)
	h := sha256.New()
	encoder := json.NewEncoder(h)
	// JSON records frame arbitrary path/content bytes without delimiter collisions.
	// A fixed digest represents regular file content; type and execute bits are separate.
	for _, relative := range paths {
		if relative == "" {
			continue
		}
		if err := ctx.Err(); err != nil {
			return "", err
		}
		record := struct {
			Path       []byte
			Kind       string
			Executable bool
			Content    []byte
		}{Path: []byte(relative)}
		path := filepath.Join(root, relative)
		info, err := os.Lstat(path)
		switch {
		case os.IsNotExist(err):
			record.Kind = "missing"
		case err != nil:
			return "", err
		case info.Mode()&os.ModeSymlink != 0:
			target, err := os.Readlink(path)
			if err != nil {
				return "", err
			}
			record.Kind, record.Content = "symlink", []byte(target)
		case info.Mode().IsRegular():
			if !contained(root, path) {
				return "", fmt.Errorf("cannot hash file outside repository: %s", relative)
			}
			f, err := os.Open(path)
			if err != nil {
				return "", err
			}
			fileHash := sha256.New()
			_, err = io.Copy(fileHash, f)
			closeErr := f.Close()
			if err != nil {
				return "", err
			}
			if closeErr != nil {
				return "", closeErr
			}
			record.Kind, record.Content = "file", fileHash.Sum(nil)
			record.Executable = info.Mode().Perm()&0111 != 0
		default:
			return "", fmt.Errorf("unsupported tracked or unignored file type: %s", relative)
		}
		if err := encoder.Encode(record); err != nil {
			return "", err
		}
	}
	return hex.EncodeToString(h.Sum(nil)), nil
}
