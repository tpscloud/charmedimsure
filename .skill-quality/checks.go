// Copyright 2026 John Deatherage (@coversheets). SPDX-License-Identifier: MIT
package main

import (
	"context"
	"fmt"
	"io"
	"os/exec"
	"path/filepath"
	"syscall"
	"time"
)

// Linux and macOS are supported. Kill the process group so a timeout does not
// leave grandchildren running after their command's parent exits.
func execute(ctx context.Context, cwd string, argv []string, out, diagnostic io.Writer) error {
	cmd := exec.CommandContext(ctx, argv[0], argv[1:]...)
	cmd.Dir, cmd.Stdout, cmd.Stderr = cwd, out, diagnostic
	configureProcessGroup(cmd)
	return cmd.Run()
}

func configureProcessGroup(cmd *exec.Cmd) {
	cmd.SysProcAttr = &syscall.SysProcAttr{Setpgid: true}
	cmd.Cancel = func() error { return syscall.Kill(-cmd.Process.Pid, syscall.SIGKILL) }
	cmd.WaitDelay = 2 * time.Second
}

func runChecks(ctx context.Context, path string, out, diagnostic io.Writer) error {
	root, err := repositoryRoot(ctx, path)
	if err != nil {
		return err
	}
	c, err := readConfig(root)
	if err != nil {
		return err
	}
	if len(c.Checks) == 0 {
		fmt.Fprintln(out, "No repository checks declared; behavioral evaluation remains untested.")
	}
	for _, check := range c.Checks {
		fmt.Fprintf(out, "Running %s\n", check.Name)
		cwd := filepath.Join(root, check.CWD)
		if !contained(root, cwd) {
			return fmt.Errorf("check cwd escapes repository")
		}
		checkCtx, cancel := context.WithTimeout(ctx, 10*time.Minute)
		err := execute(checkCtx, cwd, check.Command, out, diagnostic)
		cancel()
		if err != nil {
			return fmt.Errorf("check %q failed: %w", check.Name, err)
		}
	}
	return nil
}
