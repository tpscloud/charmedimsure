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
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

func hashBytes(b []byte) string { sum := sha256.Sum256(b); return hex.EncodeToString(sum[:]) }

func writeJSON(path string, value any) error {
	b, err := json.MarshalIndent(value, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(path, append(b, '\n'), 0600)
}

// Resolve existing ancestors too, so a not-yet-created path cannot hide an escape.
func canonical(path string) (string, error) {
	path, err := filepath.Abs(path)
	if err != nil {
		return "", err
	}
	if _, err := os.Lstat(path); err == nil {
		return filepath.EvalSymlinks(path)
	} else if !os.IsNotExist(err) {
		return "", err
	}
	parent := filepath.Dir(path)
	if parent == path {
		return "", fmt.Errorf("cannot resolve %s", path)
	}
	parent, err = canonical(parent)
	if err != nil {
		return "", err
	}
	return filepath.Join(parent, filepath.Base(path)), nil
}

func contained(root, path string) bool {
	resolved, err := canonical(path)
	if err != nil {
		return false
	}
	rel, err := filepath.Rel(root, resolved)
	return err == nil && filepath.IsLocal(rel)
}

func git(ctx context.Context, root string, args ...string) ([]byte, error) {
	cmd := exec.CommandContext(ctx, "git", append([]string{"-C", root}, args...)...)
	var stderr bytes.Buffer
	cmd.Stderr = &stderr
	b, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("git %s: %w: %s", args[0], err, strings.TrimSpace(stderr.String()))
	}
	return b, nil
}

func repositoryRoot(ctx context.Context, path string) (string, error) {
	root, err := canonical(path)
	if err != nil {
		return "", err
	}
	b, err := git(ctx, root, "rev-parse", "--show-toplevel")
	if err != nil {
		return "", err
	}
	top, err := canonical(strings.TrimSpace(string(b)))
	if err != nil {
		return "", err
	}
	if top != root {
		return "", fmt.Errorf("select the repository root: %s", top)
	}
	return root, nil
}

// Reject duplicate keys rather than letting a later value change what gets checked.
func uniqueJSON(raw []byte, target any) error {
	d := json.NewDecoder(bytes.NewReader(raw))
	var value func() error
	value = func() error {
		t, err := d.Token()
		if err != nil {
			return err
		}
		if delim, ok := t.(json.Delim); ok {
			switch delim {
			case '{':
				seen := map[string]bool{}
				for d.More() {
					key, err := d.Token()
					if err != nil {
						return err
					}
					s, ok := key.(string)
					if !ok {
						return fmt.Errorf("invalid JSON key")
					}
					if seen[s] {
						return fmt.Errorf("duplicate JSON key: %s", s)
					}
					seen[s] = true
					if err := value(); err != nil {
						return err
					}
				}
			case '[':
				for d.More() {
					if err := value(); err != nil {
						return err
					}
				}
			default:
				return fmt.Errorf("unexpected JSON delimiter")
			}
			_, err = d.Token()
			return err
		}
		return nil
	}
	if err := value(); err != nil {
		return err
	}
	if _, err := d.Token(); err != io.EOF {
		return fmt.Errorf("trailing JSON data")
	}
	return json.Unmarshal(raw, target)
}

type check struct {
	Name    string   `json:"name"`
	Command []string `json:"command"`
	CWD     string   `json:"cwd,omitempty"`
}
type configuration struct {
	Version     int               `json:"version"`
	Skills      []string          `json:"skills"`
	Checks      []check           `json:"checks"`
	Frontmatter frontmatterPolicy `json:"frontmatter,omitzero"`
}

type frontmatterPolicy struct {
	AllowedKeys          []string `json:"allowedKeys,omitempty"`
	RequireDirectoryName bool     `json:"requireDirectoryName,omitzero"`
}

func readConfig(root string) (configuration, error) {
	var c configuration
	path := filepath.Join(root, ".skill-quality.json")
	if !contained(root, path) {
		return c, fmt.Errorf("configuration escapes repository")
	}
	b, err := os.ReadFile(path)
	if err != nil {
		return c, err
	}
	if err := uniqueJSON(b, &c); err != nil {
		return c, err
	}
	var shape map[string]any
	if err := json.Unmarshal(b, &shape); err != nil {
		return c, err
	}
	for key := range shape {
		switch key {
		case "version", "skills", "checks", "frontmatter":
		default:
			return c, fmt.Errorf("unknown configuration property: %s", key)
		}
	}
	if raw, exists := shape["frontmatter"]; exists {
		policy, ok := raw.(map[string]any)
		if !ok {
			return c, fmt.Errorf("frontmatter must be an object")
		}
		for key, value := range policy {
			switch key {
			case "allowedKeys":
				keys, ok := value.([]any)
				if !ok || len(keys) == 0 {
					return c, fmt.Errorf("frontmatter.allowedKeys must be a nonempty array of strings")
				}
				seen := map[string]bool{}
				for _, item := range keys {
					key, ok := item.(string)
					if !ok || strings.TrimSpace(key) == "" || seen[key] {
						return c, fmt.Errorf("frontmatter.allowedKeys entries must be nonempty, unique strings")
					}
					seen[key] = true
				}
			case "requireDirectoryName":
				if _, ok := value.(bool); !ok {
					return c, fmt.Errorf("frontmatter.requireDirectoryName must be a boolean")
				}
			default:
				return c, fmt.Errorf("unknown frontmatter policy: %s", key)
			}
		}
	}
	if rawChecks, exists := shape["checks"]; exists {
		checks, ok := rawChecks.([]any)
		if !ok {
			return c, fmt.Errorf("checks must be an array")
		}
		for _, item := range checks {
			object, ok := item.(map[string]any)
			if !ok {
				return c, fmt.Errorf("each check must be an object")
			}
			for key := range object {
				switch key {
				case "name", "command", "cwd":
				default:
					return c, fmt.Errorf("unknown check property: %s", key)
				}
			}
			if cwd, exists := object["cwd"]; exists {
				if _, ok := cwd.(string); !ok {
					return c, fmt.Errorf("check cwd must be a string")
				}
			}
			args, ok := object["command"].([]any)
			if !ok {
				return c, fmt.Errorf("check commands must be argument arrays")
			}
			for _, arg := range args {
				if _, ok := arg.(string); !ok {
					return c, fmt.Errorf("check arguments must be strings")
				}
			}
		}
	}
	if c.Version != 0 && c.Version != 1 {
		return c, fmt.Errorf("unsupported configuration version %d", c.Version)
	}
	if len(c.Skills) == 0 {
		return c, fmt.Errorf("skills must be a nonempty array")
	}
	seen := map[string]bool{}
	for _, path := range c.Skills {
		if path == "" || seen[path] {
			return c, fmt.Errorf("empty or duplicate skill path")
		}
		seen[path] = true
	}
	for _, check := range c.Checks {
		if strings.TrimSpace(check.Name) == "" || len(check.Command) == 0 || check.Command[0] == "" {
			return c, fmt.Errorf("checks need a name and nonempty command argument array")
		}
		cwd := check.CWD
		if cwd == "" {
			cwd = "."
		}
		if !filepath.IsLocal(cwd) || !contained(root, filepath.Join(root, cwd)) {
			return c, fmt.Errorf("check cwd escapes repository")
		}
	}
	return c, nil
}
