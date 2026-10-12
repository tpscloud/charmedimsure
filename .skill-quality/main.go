// Copyright 2026 John Deatherage (@coversheets). SPDX-License-Identifier: MIT
package main

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/signal"
	"runtime/debug"
	"syscall"
	"time"

	"github.com/spf13/cobra"
)

// Release builds set version; ordinary builds remain visibly development builds.
var version = "dev"

var errFindings = errors.New("checks found problems")

func buildRevision() string {
	if info, ok := debug.ReadBuildInfo(); ok {
		for _, setting := range info.Settings {
			if setting.Key == "vcs.revision" {
				return setting.Value
			}
		}
	}
	return "unrecorded"
}

func buildDirty() bool {
	if info, ok := debug.ReadBuildInfo(); ok {
		for _, setting := range info.Settings {
			if setting.Key == "vcs.modified" {
				return setting.Value != "false"
			}
		}
	}
	return true
}

func newRoot(out, diagnostic io.Writer) *cobra.Command {
	root := &cobra.Command{Use: "skillforge", Short: "Validate skill repositories and run local CI checks", Version: fmt.Sprintf("%s (revision %s, dirty %t)", version, buildRevision(), buildDirty()), SilenceUsage: true, SilenceErrors: true}
	root.SetOut(out)
	root.SetErr(diagnostic)
	repository := func(args []string) string {
		if len(args) == 0 {
			return "."
		}
		return args[0]
	}
	var reportPath string
	validate := &cobra.Command{Use: "validate [repository]", Short: "Inspect declared skills without running repository code", Args: cobra.MaximumNArgs(1), RunE: func(cmd *cobra.Command, args []string) error {
		report := validateRepository(cmd.Context(), repository(args))
		if reportPath != "" {
			if err := writeJSON(reportPath, report); err != nil {
				return err
			}
		}
		for _, f := range report.Findings {
			fmt.Fprintf(cmd.OutOrStdout(), "%s: %s: %s\n", f.Severity, f.Path, f.Message)
		}
		fmt.Fprintf(cmd.OutOrStdout(), "Skill validation: %s\n", report.Status)
		if report.Status != "pass" {
			return errFindings
		}
		return nil
	}}
	validate.Flags().StringVar(&reportPath, "output", "", "Write JSON evidence to this file")
	checks := &cobra.Command{Use: "check [repository]", Short: "Execute explicitly configured argument arrays (trusted repositories only)", Args: cobra.MaximumNArgs(1), RunE: func(cmd *cobra.Command, args []string) error {
		return runChecks(cmd.Context(), repository(args), cmd.OutOrStdout(), cmd.ErrOrStderr())
	}}
	var outputDirectory string
	scan := &cobra.Command{Use: "scan [repository]", Short: "Scan Git history and current files using checksum-pinned Gitleaks", Args: cobra.MaximumNArgs(1), RunE: func(cmd *cobra.Command, args []string) error {
		return scanRepository(cmd.Context(), repository(args), outputDirectory, cmd.OutOrStdout(), cmd.ErrOrStderr())
	}}
	scan.Flags().StringVar(&outputDirectory, "output-directory", "", "Directory outside the repository for redacted evidence (required)")
	if err := scan.MarkFlagRequired("output-directory"); err != nil {
		panic(err)
	}
	var update bool
	install := &cobra.Command{Use: "install repository", Short: "Copy the embedded Go CI bundle into a selected repository", Args: cobra.ExactArgs(1), RunE: func(cmd *cobra.Command, args []string) error {
		if err := installBundle(cmd.Context(), args[0], update); err != nil {
			return err
		}
		fmt.Fprintln(cmd.OutOrStdout(), "Installed Skillforge CI bundle; inspect the diff before committing.")
		return nil
	}}
	install.Flags().BoolVar(&update, "update", false, "Replace recognized managed files; preserve repository configuration")
	var forks forkSyncOptions
	forkSync := &cobra.Command{Use: "fork-sync [directory]", Short: "Fast-forward declared forks to their upstream branch and report the result", Long: "Reads .skill-forks.json from the directory, fast-forwards each declared fork to its upstream branch, and reports every fork that is behind, diverged, missing, or not writable.\n\nThe fork token is read from SKILLFORGE_FORK_TOKEN and an optional webhook URL from SKILLFORGE_NOTIFY_WEBHOOK. A Markdown summary is appended to the file named by GITHUB_STEP_SUMMARY when that variable is set.", Args: cobra.MaximumNArgs(1), RunE: func(cmd *cobra.Command, args []string) error {
		if forks.Timeout <= 0 {
			return fmt.Errorf("timeout must be positive")
		}
		forks.Directory = repository(args)
		forks.Token, forks.Webhook, forks.SummaryPath = os.Getenv(forkTokenEnv), os.Getenv(notifyWebhookEnv), os.Getenv("GITHUB_STEP_SUMMARY")
		forks.Client = &http.Client{}
		return runForkSync(cmd.Context(), forks, cmd.OutOrStdout())
	}}
	forkSync.Flags().StringVar(&forks.Output, "output", "", "Write the JSON report to this file")
	forkSync.Flags().DurationVar(&forks.Timeout, "timeout", 10*time.Minute, "Maximum duration for each fork, including all Git operations")
	forkSync.Flags().BoolVar(&forks.CheckOnly, "check", false, "Compare forks with upstream without pushing")
	forkSync.Flags().StringVar(&forks.GitBase, "git-base", "https://github.com", "Base URL or directory that holds owner/repository.git")
	root.AddCommand(validate, checks, scan, install, forkSync)
	return root
}

func exitCode(err error) int {
	if err == nil {
		return 0
	}
	if errors.Is(err, errFindings) {
		return 1
	}
	return 2
}

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()
	err := newRoot(os.Stdout, os.Stderr).ExecuteContext(ctx)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
	}
	os.Exit(exitCode(err))
}
