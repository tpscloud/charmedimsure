// Copyright 2026 John Deatherage (@coversheets). SPDX-License-Identifier: MIT
package main

import (
	"archive/tar"
	"bytes"
	"compress/gzip"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"time"
)

const gitleaksVersion = "8.30.1"

// Extend the pinned scanner's built-in rules, never the target's configuration.
const scannerPolicy = "[extend]\nuseDefault = true\n"

var archiveHashes = map[string]string{
	"darwin/arm64": "b40ab0ae55c505963e365f271a8d3846efbc170aa17f2607f13df610a9aeb6a5",
	"darwin/amd64": "dfe101a4db2255fc85120ac7f3d25e4342c3c20cf749f2c20a18081af1952709",
	"linux/arm64":  "e4a487ee7ccd7d3a7f7ec08657610aa3606637dab924210b3aee62570fb4b080",
	"linux/amd64":  "551f6fc83ea457d62a0d98237cbad105af8d557003051f41f3e7ca7b3f2470eb",
}

func extractScanner(raw []byte, expected string) ([]byte, error) {
	if hashBytes(raw) != expected {
		return nil, fmt.Errorf("Gitleaks archive checksum mismatch")
	}
	gz, err := gzip.NewReader(bytes.NewReader(raw))
	if err != nil {
		return nil, err
	}
	defer gz.Close()
	t := tar.NewReader(gz)
	for {
		header, err := t.Next()
		if err == io.EOF {
			return nil, fmt.Errorf("Gitleaks archive has no executable")
		}
		if err != nil {
			return nil, err
		}
		if header.Name != "gitleaks" {
			continue
		}
		if header.Typeflag != tar.TypeReg || header.Size <= 0 || header.Size > 100<<20 {
			return nil, fmt.Errorf("invalid Gitleaks executable entry")
		}
		return io.ReadAll(t)
	}
}

func downloadScanner(ctx context.Context) ([]byte, error) {
	key := runtime.GOOS + "/" + runtime.GOARCH
	expected, ok := archiveHashes[key]
	if !ok {
		return nil, fmt.Errorf("unsupported scanner platform: %s", key)
	}
	arch := runtime.GOARCH
	if arch == "amd64" {
		arch = "x64"
	}
	url := fmt.Sprintf("https://github.com/gitleaks/gitleaks/releases/download/v%s/gitleaks_%s_%s_%s.tar.gz", gitleaksVersion, gitleaksVersion, runtime.GOOS, arch)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}
	client := &http.Client{Timeout: time.Minute}
	response, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()
	if response.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("Gitleaks download: HTTP %d", response.StatusCode)
	}
	raw, err := io.ReadAll(io.LimitReader(response.Body, (50<<20)+1))
	if err != nil {
		return nil, err
	}
	if len(raw) > 50<<20 {
		return nil, fmt.Errorf("Gitleaks archive exceeds size limit")
	}
	return extractScanner(raw, expected)
}

type scanMode struct {
	Mode     string `json:"mode"`
	Status   string `json:"status"`
	ExitCode int    `json:"exitCode"`
	Findings int    `json:"findings"`
}
type scanReport struct {
	Version   int        `json:"version"`
	Scanner   string     `json:"scanner"`
	Policy    string     `json:"policy"`
	CheckedAt string     `json:"checkedAt"`
	Revision  string     `json:"revision"`
	Status    string     `json:"status"`
	Modes     []scanMode `json:"modes"`
}
type executor func(context.Context, string, []string, io.Writer, io.Writer) error

func scannerResults(ctx context.Context, root, output, binary string, run executor, out, diagnostic io.Writer) ([]scanMode, error) {
	// Gitleaks also loads this file from the source even with an explicit
	// --gitleaks-ignore-path. Refuse it without changing the target checkout.
	if _, err := os.Lstat(filepath.Join(root, ".gitleaksignore")); err == nil {
		return nil, fmt.Errorf("baseline scan refuses target .gitleaksignore; review exclusions and scan a separate copy without that file")
	} else if !os.IsNotExist(err) {
		return nil, err
	}
	policyDir, err := os.MkdirTemp("", "skillforge-scan-policy-")
	if err != nil {
		return nil, err
	}
	defer os.RemoveAll(policyDir)
	configPath := filepath.Join(policyDir, "gitleaks.toml")
	if err := os.WriteFile(configPath, []byte(scannerPolicy), 0600); err != nil {
		return nil, err
	}
	ignorePath := filepath.Join(policyDir, ".gitleaksignore")
	if err := os.WriteFile(ignorePath, nil, 0600); err != nil {
		return nil, err
	}
	var modes []scanMode
	var failures []error
	found := false
	for _, mode := range []string{"git", "dir"} {
		path := filepath.Join(output, "gitleaks-"+mode+".json")
		args := []string{binary, mode, root, "--redact=100", "--no-banner", "--exit-code", "10", "--report-format", "json", "--report-path", path,
			"--config", configPath, "--gitleaks-ignore-path", ignorePath, "--ignore-gitleaks-allow"}
		if mode == "git" {
			args = append(args, "--log-opts=--all")
		}
		scanCtx, cancel := context.WithTimeout(ctx, 10*time.Minute)
		err := run(scanCtx, root, args, out, diagnostic)
		cancel()
		code := 0
		if err != nil {
			code = -1
			if exit, ok := errors.AsType[*exec.ExitError](err); ok {
				code = exit.ExitCode()
			}
		}
		result := scanMode{Mode: mode, Status: "error", ExitCode: code}
		raw, readErr := os.ReadFile(path)
		var entries []json.RawMessage
		validEvidence := readErr == nil && len(bytes.TrimSpace(raw)) > 0 && bytes.TrimSpace(raw)[0] == '[' && json.Unmarshal(raw, &entries) == nil
		result.Findings = len(entries)
		switch {
		case code != 0 && code != 10:
			failures = append(failures, fmt.Errorf("Gitleaks %s did not complete: %w", mode, err))
		case !validEvidence:
			failures = append(failures, fmt.Errorf("Gitleaks %s has missing or invalid report evidence", mode))
		case code == 0 && len(entries) != 0 || code == 10 && len(entries) == 0:
			failures = append(failures, fmt.Errorf("Gitleaks %s exit code disagrees with evidence", mode))
		case code == 10:
			found = true
			result.Status = "findings"
		default:
			result.Status = "pass"
		}
		modes = append(modes, result)
	}
	if len(failures) != 0 {
		return modes, errors.Join(failures...)
	}
	if found {
		return modes, errFindings
	}
	return modes, nil
}

func scanRepository(ctx context.Context, path, output string, out, diagnostic io.Writer) (resultErr error) {
	root, err := repositoryRoot(ctx, path)
	if err != nil {
		return err
	}
	if output == "" {
		return fmt.Errorf("output directory is required")
	}
	output, err = canonical(output)
	if err != nil {
		return err
	}
	if contained(root, output) {
		return fmt.Errorf("write scanner reports outside the repository")
	}
	// A fresh directory prevents a crashed scanner from reusing old clean evidence.
	if err := os.MkdirAll(filepath.Dir(output), 0700); err != nil {
		return err
	}
	if err := os.Mkdir(output, 0700); err != nil {
		return fmt.Errorf("select a new evidence directory: %w", err)
	}
	report := scanReport{Version: 1, Scanner: "gitleaks " + gitleaksVersion, Policy: "pinned-defaults-no-target-suppressions", CheckedAt: time.Now().UTC().Format(time.RFC3339), Status: "error", Modes: []scanMode{}}
	if revision, err := git(ctx, root, "rev-parse", "HEAD"); err == nil {
		report.Revision = string(bytes.TrimSpace(revision))
	}
	defer func() { resultErr = errors.Join(resultErr, writeJSON(filepath.Join(output, "summary.json"), report)) }()
	shallow, err := git(ctx, root, "rev-parse", "--is-shallow-repository")
	if err != nil {
		return err
	}
	if string(bytes.TrimSpace(shallow)) != "false" {
		return fmt.Errorf("history scan requires a complete clone; fetch full history before scanning")
	}
	b, err := downloadScanner(ctx)
	if err != nil {
		return err
	}
	temp, err := os.MkdirTemp("", "skillforge-gitleaks-")
	if err != nil {
		return err
	}
	defer os.RemoveAll(temp)
	binary := filepath.Join(temp, "gitleaks")
	if err := os.WriteFile(binary, b, 0700); err != nil {
		return err
	}
	report.Modes, err = scannerResults(ctx, root, output, binary, execute, out, diagnostic)
	switch {
	case err == nil:
		report.Status = "pass"
	case errors.Is(err, errFindings):
		report.Status = "findings"
	}
	return err
}
