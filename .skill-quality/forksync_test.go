// Copyright 2026 John Deatherage (@coversheets). SPDX-License-Identifier: MIT
package main

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"os/exec"
	"path/filepath"
	"slices"
	"strings"
	"sync/atomic"
	"testing"
	"time"
)

// forkFixture holds local bare repositories laid out as base/owner/name.git,
// so the command under test runs real fetches and pushes without a network.
type forkFixture struct {
	t                *testing.T
	base, work, conf string
}

func newForkFixture(t *testing.T) *forkFixture {
	t.Helper()
	// Keep the developer's signing, hooks, and credential helpers out of the fixtures.
	t.Setenv("GIT_CONFIG_GLOBAL", os.DevNull)
	t.Setenv("GIT_CONFIG_SYSTEM", os.DevNull)
	for _, name := range []string{"GIT_AUTHOR_NAME", "GIT_COMMITTER_NAME"} {
		t.Setenv(name, "Fixture")
	}
	for _, name := range []string{"GIT_AUTHOR_EMAIL", "GIT_COMMITTER_EMAIL"} {
		t.Setenv(name, "fixture@example.invalid")
	}
	root, err := canonical(t.TempDir())
	if err != nil {
		t.Fatal(err)
	}
	f := &forkFixture{t: t, base: filepath.Join(root, "remotes"), work: filepath.Join(root, "work"), conf: filepath.Join(root, "conf")}
	for _, directory := range []string{filepath.Join(f.base, "up"), filepath.Join(f.base, "me"), f.work, f.conf} {
		if err := os.MkdirAll(directory, 0755); err != nil {
			t.Fatal(err)
		}
	}
	f.run(f.base, "init", "-q", "--bare", "-b", "main", "up/proj.git")
	f.run(f.work, "init", "-q", "-b", "main")
	f.run(f.work, "commit", "-q", "--allow-empty", "-m", "root")
	f.run(f.work, "push", "-q", f.remote("up/proj"), "main")
	f.run(f.base, "clone", "-q", "--bare", f.remote("up/proj"), "me/proj.git")
	f.configure(`{"version":1,"forks":[{"fork":"me/proj","upstream":"up/proj","branch":"main"}]}`)
	return f
}

func (f *forkFixture) remote(name string) string { return remoteURL(f.base, name) }

func (f *forkFixture) run(directory string, args ...string) string {
	f.t.Helper()
	out, err := exec.CommandContext(f.t.Context(), "git", append([]string{"-C", directory}, args...)...).CombinedOutput()
	if err != nil {
		f.t.Fatalf("git %v: %v: %s", args, err, out)
	}
	return strings.TrimSpace(string(out))
}

// commit adds one commit on top of the named remote's main branch and pushes it there.
func (f *forkFixture) commit(name, message string) string {
	f.t.Helper()
	f.run(f.work, "fetch", "-q", f.remote(name), "main")
	f.run(f.work, "checkout", "-q", "-B", "main", "FETCH_HEAD")
	f.run(f.work, "commit", "-q", "--allow-empty", "-m", message)
	f.run(f.work, "push", "-q", f.remote(name), "main")
	return f.head(name)
}

func (f *forkFixture) head(name string) string {
	f.t.Helper()
	return f.run(f.remote(name), "rev-parse", "refs/heads/main")
}

func (f *forkFixture) configure(text string) {
	f.t.Helper()
	if err := os.WriteFile(filepath.Join(f.conf, forkConfigName), []byte(text), 0644); err != nil {
		f.t.Fatal(err)
	}
}

func (f *forkFixture) sync(change func(*forkSyncOptions)) (forkReport, error) {
	f.t.Helper()
	o := forkSyncOptions{Directory: f.conf, GitBase: f.base, Output: filepath.Join(f.conf, "report.json"), Client: http.DefaultClient}
	if change != nil {
		change(&o)
	}
	runErr := runForkSync(f.t.Context(), o, io.Discard)
	var report forkReport
	b, err := os.ReadFile(o.Output)
	if err != nil {
		f.t.Fatal(err)
	}
	if err := json.Unmarshal(b, &report); err != nil {
		f.t.Fatal(err)
	}
	return report, runErr
}

func TestForkSyncFastForwardsThenReportsCurrent(t *testing.T) {
	f := newForkFixture(t)
	f.commit("up/proj", "one")
	want := f.commit("up/proj", "two")
	var posts atomic.Int32
	var body atomic.Value
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		b, _ := io.ReadAll(r.Body)
		body.Store(string(b))
		posts.Add(1)
	}))
	defer server.Close()
	summary := filepath.Join(f.conf, "summary.md")
	options := func(o *forkSyncOptions) { o.Webhook, o.SummaryPath = server.URL, summary }
	report, err := f.sync(options)
	if err != nil || report.Status != "pass" || report.Forks[0].Status != "synced" || report.Forks[0].Behind != 2 {
		t.Fatalf("sync: %v %+v", err, report)
	}
	if f.head("me/proj") != want || report.Forks[0].UpstreamRevision != want {
		t.Fatal("fork was not fast-forwarded to upstream")
	}
	if report.Notification.Webhook != "delivered" || posts.Load() != 1 || !strings.Contains(body.Load().(string), "synced") {
		t.Fatalf("webhook: %+v", report.Notification)
	}
	if b, err := os.ReadFile(summary); err != nil || !strings.Contains(string(b), "| me/proj | up/proj | main | synced | 2 | 0 |") {
		t.Fatalf("summary: %v %s", err, b)
	}
	report, err = f.sync(options)
	if err != nil || report.Forks[0].Status != "current" || report.Notification.Webhook != "not-needed" || posts.Load() != 1 {
		t.Fatalf("second run: %v %+v", err, report)
	}
}

func TestForkSyncCheckModeDoesNotPush(t *testing.T) {
	f := newForkFixture(t)
	before := f.head("me/proj")
	f.commit("up/proj", "one")
	report, err := f.sync(func(o *forkSyncOptions) { o.CheckOnly = true })
	if !errors.Is(err, errFindings) || report.Mode != "check" || report.Forks[0].Status != "behind" || f.head("me/proj") != before {
		t.Fatalf("check: %v %+v", err, report)
	}
}

func TestForkSyncLeavesDivergedForkUnchanged(t *testing.T) {
	f := newForkFixture(t)
	mine := f.commit("me/proj", "local work")
	f.commit("up/proj", "one")
	report, err := f.sync(nil)
	r := report.Forks[0]
	if !errors.Is(err, errFindings) || report.Status != "fail" || r.Status != "diverged" || r.Ahead != 1 || r.Behind != 1 {
		t.Fatalf("diverged: %v %+v", err, report)
	}
	if f.head("me/proj") != mine {
		t.Fatal("diverged fork was rewritten")
	}
}

func TestForkSyncMissingForkDoesNotStopOtherForks(t *testing.T) {
	f := newForkFixture(t)
	want := f.commit("up/proj", "one")
	f.configure(`{"version":1,"forks":[{"fork":"me/absent","upstream":"up/proj","branch":"main"},{"fork":"me/proj","upstream":"up/proj","branch":"gone"},{"fork":"me/proj","upstream":"up/proj","branch":"main"}]}`)
	report, err := f.sync(nil)
	if !errors.Is(err, errForkSync) || exitCode(err) != 2 || report.Status != "error" {
		t.Fatalf("missing: %v %+v", err, report)
	}
	if report.Forks[0].Status != "missing" || report.Forks[1].Status != "missing" || report.Forks[2].Status != "synced" || f.head("me/proj") != want {
		t.Fatalf("statuses: %+v", report.Forks)
	}
}

func TestForkSyncReportsInsufficientToken(t *testing.T) {
	f := newForkFixture(t)
	before := f.head("me/proj")
	f.commit("up/proj", "workflow change")
	hook := "#!/bin/sh\necho 'refusing to allow a Personal Access Token to create or update workflow `.github/workflows/ci.yml` without `workflow` scope' >&2\nexit 1\n"
	if err := os.WriteFile(filepath.Join(f.remote("me/proj"), "hooks", "pre-receive"), []byte(hook), 0755); err != nil {
		t.Fatal(err)
	}
	var posted atomic.Value
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		b, _ := io.ReadAll(r.Body)
		posted.Store(string(b))
	}))
	defer server.Close()
	report, err := f.sync(func(o *forkSyncOptions) { o.Webhook = server.URL })
	r := report.Forks[0]
	if !errors.Is(err, errForkSync) || r.Status != "token-insufficient" || !strings.Contains(r.Message, forkTokenEnv) || !strings.Contains(r.Message, "workflow") {
		t.Fatalf("token: %v %+v", err, report)
	}
	if f.head("me/proj") != before {
		t.Fatal("fork changed after a refused push")
	}
	if text, _ := posted.Load().(string); !strings.Contains(text, "token-insufficient") {
		t.Fatalf("owner was not notified: %q", text)
	}
}

func TestForkSyncNotificationFailureIsAnError(t *testing.T) {
	f := newForkFixture(t)
	want := f.commit("up/proj", "one")
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "no", http.StatusInternalServerError)
	}))
	defer server.Close()
	report, err := f.sync(func(o *forkSyncOptions) { o.Webhook = server.URL + "/secret-path" })
	if !errors.Is(err, errForkSync) || report.Status != "error" || report.Notification.Webhook != "failed" || report.Forks[0].Status != "synced" || f.head("me/proj") != want {
		t.Fatalf("webhook failure: %v %+v", err, report)
	}
	f.commit("up/proj", "two")
	server.Close()
	report, err = f.sync(func(o *forkSyncOptions) {
		o.Webhook, o.SummaryPath = server.URL+"/secret-path", filepath.Join(f.conf, "absent", "summary.md")
	})
	if !errors.Is(err, errForkSync) || report.Notification.Webhook != "failed" || report.Notification.Summary != "failed" {
		t.Fatalf("unreachable webhook: %v %+v", err, report)
	}
	if strings.Contains(report.Notification.Message, "secret-path") {
		t.Fatalf("webhook URL leaked: %s", report.Notification.Message)
	}
}

func TestForkSyncPin(t *testing.T) {
	f := newForkFixture(t)
	pinned := f.commit("up/proj", "pinned")
	f.commit("up/proj", "later")
	f.configure(`{"version":1,"forks":[{"fork":"me/proj","upstream":"up/proj","branch":"main","pin":"` + pinned[:12] + `"}]}`)
	report, err := f.sync(nil)
	pin := report.Forks[0].Pin
	if err != nil || report.Forks[0].Status != "synced" || pin == nil || !pin.InFork || pin.UpstreamCommitsPast != 1 {
		t.Fatalf("pin: %v %+v %+v", err, report, pin)
	}
	f.configure(`{"version":1,"forks":[{"fork":"me/proj","upstream":"up/proj","branch":"main","pin":"0123456789abcdef"}]}`)
	report, err = f.sync(nil)
	if !errors.Is(err, errFindings) || report.Forks[0].Status != "pin-missing" || report.Forks[0].Pin.InFork {
		t.Fatalf("absent pin: %v %+v", err, report)
	}
}

func TestForkConfigRejectsUnsafeOrAmbiguousInput(t *testing.T) {
	entry := func(fork, upstream, branch, extra string) string {
		return `{"version":1,"forks":[{"fork":"` + fork + `","upstream":"` + upstream + `","branch":"` + branch + `"` + extra + `}]}`
	}
	for name, text := range map[string]string{
		"empty list":      `{"version":1,"forks":[]}`,
		"version":         `{"version":2,"forks":[{"fork":"a/b","upstream":"c/d","branch":"main"}]}`,
		"unknown key":     entry("a/b", "c/d", "main", `,"force":true`),
		"duplicate key":   `{"version":1,"version":1,"forks":[{"fork":"a/b","upstream":"c/d","branch":"main"}]}`,
		"option as name":  entry("-x/b", "c/d", "main", ""),
		"path escape":     entry("../b", "c/d", "main", ""),
		"url":             entry("https://example.invalid/a/b", "c/d", "main", ""),
		"option branch":   entry("a/b", "c/d", "--force", ""),
		"refspec branch":  entry("a/b", "c/d", "+main:other", ""),
		"self upstream":   entry("a/b", "A/B", "main", ""),
		"pin not a hash":  entry("a/b", "c/d", "main", `,"pin":"v1.0"`),
		"duplicate entry": `{"version":1,"forks":[{"fork":"a/b","upstream":"c/d","branch":"main"},{"fork":"a/b","upstream":"e/f","branch":"main"}]}`,
	} {
		directory := t.TempDir()
		if err := os.WriteFile(filepath.Join(directory, forkConfigName), []byte(text), 0644); err != nil {
			t.Fatal(err)
		}
		if _, err := readForkConfig(directory); err == nil {
			t.Errorf("%s: accepted %s", name, text)
		}
	}
	err := runForkSync(t.Context(), forkSyncOptions{Directory: t.TempDir()}, io.Discard)
	if err == nil || errors.Is(err, errFindings) || exitCode(err) != 2 {
		t.Fatalf("missing configuration: %v", err)
	}
}

func TestForkTokenHandling(t *testing.T) {
	env := forkGitEnv("https://github.com", "s3cret-token")
	if !slices.Contains(env, "GIT_CONFIG_KEY_0=http.https://github.com/.extraheader") || !slices.Contains(env, "GIT_TERMINAL_PROMPT=0") {
		t.Fatal("token header was not scoped to the fork host")
	}
	for _, base := range []string{"http://example.invalid", "/tmp/remotes"} {
		if strings.Contains(strings.Join(forkGitEnv(base, "s3cret-token"), "\n"), "extraheader") {
			t.Fatalf("token offered to %s", base)
		}
	}
	_, detail, err := forkGit(t.Context(), env, "s3cret-token", t.TempDir(), "ls-remote", "/nonexistent/s3cret-token.git")
	if err == nil || strings.Contains(detail, "s3cret-token") || classifyForkFailure(detail) != "missing" {
		t.Fatalf("redaction: %v %q", err, detail)
	}
	for text, want := range map[string]string{
		"remote: Permission to me/proj.git denied to someone. fatal: The requested URL returned error: 403":                    "token-insufficient",
		"refusing to allow a GitHub App to create or update workflow `.github/workflows/a.yml` without `workflows` permission": "token-insufficient",
		"remote: Repository not found.":                   "missing",
		"fatal: unable to access: Could not resolve host": "failed",
	} {
		if got := classifyForkFailure(text); got != want {
			t.Errorf("%q classified as %s, want %s", text, got, want)
		}
	}
}

func TestForkSyncCommand(t *testing.T) {
	f := newForkFixture(t)
	f.commit("up/proj", "one")
	summary := filepath.Join(f.conf, "summary.md")
	t.Setenv("GITHUB_STEP_SUMMARY", summary)
	t.Setenv(forkTokenEnv, "")
	t.Setenv(notifyWebhookEnv, "")
	var out bytes.Buffer
	root := newRoot(&out, io.Discard)
	root.SetArgs([]string{"fork-sync", f.conf, "--git-base", f.base, "--check"})
	if err := root.ExecuteContext(t.Context()); exitCode(err) != 1 || !strings.Contains(out.String(), "behind: me/proj from up/proj (main)") {
		t.Fatalf("check command: %v %s", err, out.String())
	}
	root = newRoot(&out, io.Discard)
	root.SetArgs([]string{"fork-sync", f.conf, "--git-base", f.base})
	if err := root.ExecuteContext(t.Context()); err != nil || f.head("me/proj") != f.head("up/proj") {
		t.Fatalf("sync command: %v %s", err, out.String())
	}
	if b, err := os.ReadFile(summary); err != nil || !strings.Contains(string(b), "## Fork sync: pass") {
		t.Fatalf("summary: %v %s", err, b)
	}
}

func TestWebhookErrorsNeverDiscloseURLs(t *testing.T) {
	for _, location := range []string{"http://invalid:port", "http://unreachable.invalid"} {
		t.Run(location, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("Location", location+r.URL.RequestURI())
				w.WriteHeader(http.StatusTemporaryRedirect)
			}))
			defer server.Close()
			err := postWebhook(t.Context(), server.Client(), server.URL+"/synthetic-credential", forkReport{})
			if err == nil || strings.Contains(err.Error(), "synthetic-credential") || strings.Contains(err.Error(), server.URL) || strings.Contains(err.Error(), location) {
				t.Fatalf("unsafe webhook error: %v", err)
			}
		})
	}
}

func TestForkTimeoutContinuesAndWritesReport(t *testing.T) {
	f := newForkFixture(t)
	gitPath, err := exec.LookPath("git")
	if err != nil {
		t.Fatal(err)
	}
	bin := t.TempDir()
	// Only the first fork stalls. The child retains pipes until group cancellation.
	script := "#!/bin/sh\ncase \"$*\" in *me/stalled.git*) sleep 30 & wait; exit 1;; esac\nexec \"$SKILLFORGE_TEST_GIT\" \"$@\"\n"
	if err := os.WriteFile(filepath.Join(bin, "git"), []byte(script), 0755); err != nil {
		t.Fatal(err)
	}
	t.Setenv("SKILLFORGE_TEST_GIT", gitPath)
	t.Setenv("PATH", bin+string(os.PathListSeparator)+os.Getenv("PATH"))
	f.configure(`{"version":1,"forks":[{"fork":"me/stalled","upstream":"up/proj","branch":"main"},{"fork":"me/proj","upstream":"up/proj","branch":"main"}]}`)
	start := time.Now()
	report, err := f.sync(func(o *forkSyncOptions) { o.Timeout = time.Second })
	if exitCode(err) != 2 || len(report.Forks) != 2 || report.Forks[0].Status != "failed" || report.Forks[1].Status != "current" || !strings.Contains(report.Forks[0].Message, "deadline exceeded") {
		t.Fatalf("timeout: %v %+v", err, report)
	}
	if time.Since(start) > 5*time.Second {
		t.Fatal("child process delayed cancellation")
	}
}

func TestForkGitCancellation(t *testing.T) {
	ctx, cancel := context.WithCancel(t.Context())
	cancel()
	_, _, err := forkGit(ctx, os.Environ(), "", t.TempDir(), "status")
	if !errors.Is(err, context.Canceled) {
		t.Fatalf("canceled git: %v", err)
	}
}
