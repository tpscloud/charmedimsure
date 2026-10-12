// Copyright 2026 John Deatherage (@coversheets). SPDX-License-Identifier: MIT
package main

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"time"
)

const (
	forkConfigName   = ".skill-forks.json"
	forkTokenEnv     = "SKILLFORGE_FORK_TOKEN"
	notifyWebhookEnv = "SKILLFORGE_NOTIFY_WEBHOOK"
	tokenHelp        = "The token in " + forkTokenEnv + " needs write access to the fork's contents and workflows (fine-grained: Contents and Workflows read and write; classic: repo and workflow scopes)."
)

var (
	errForkSync    = errors.New("fork sync could not complete or notify; see the report")
	repositoryName = regexp.MustCompile(`^[A-Za-z0-9_][A-Za-z0-9_.-]*/[A-Za-z0-9_][A-Za-z0-9_.-]*$`)
	branchName     = regexp.MustCompile(`^[A-Za-z0-9_][A-Za-z0-9_./-]*$`)
	pinRevision    = regexp.MustCompile(`^[0-9a-f]{7,64}$`)
)

type forkEntry struct {
	Fork     string `json:"fork"`
	Upstream string `json:"upstream"`
	Branch   string `json:"branch"`
	Pin      string `json:"pin,omitempty"`
}
type forkConfig struct {
	Version int         `json:"version"`
	Forks   []forkEntry `json:"forks"`
}
type pinResult struct {
	Revision            string `json:"revision"`
	InFork              bool   `json:"inFork"`
	UpstreamCommitsPast int    `json:"upstreamCommitsPast"`
}
type forkResult struct {
	Fork             string     `json:"fork"`
	Upstream         string     `json:"upstream"`
	Branch           string     `json:"branch"`
	Status           string     `json:"status"`
	Message          string     `json:"message"`
	ForkRevision     string     `json:"forkRevision,omitempty"`
	UpstreamRevision string     `json:"upstreamRevision,omitempty"`
	Behind           int        `json:"behind"`
	Ahead            int        `json:"ahead"`
	Pin              *pinResult `json:"pin,omitempty"`
}
type forkNotification struct {
	Summary string `json:"summary"`
	Webhook string `json:"webhook"`
	Message string `json:"message,omitempty"`
}
type forkReport struct {
	Version      int              `json:"version"`
	CheckedAt    string           `json:"checkedAt"`
	Mode         string           `json:"mode"`
	Status       string           `json:"status"`
	Forks        []forkResult     `json:"forks"`
	Notification forkNotification `json:"notification"`
}
type forkSyncOptions struct {
	Directory, GitBase, Output, Token, Webhook, SummaryPath string
	CheckOnly                                               bool
	Timeout                                                 time.Duration
	Client                                                  *http.Client
}

func readForkConfig(directory string) (forkConfig, error) {
	var c forkConfig
	raw, err := os.ReadFile(filepath.Join(directory, forkConfigName))
	if err != nil {
		return c, err
	}
	if err := uniqueJSON(raw, &c); err != nil {
		return c, err
	}
	// A misspelled key must not silently drop a fork or a pin.
	strict := json.NewDecoder(bytes.NewReader(raw))
	strict.DisallowUnknownFields()
	c = forkConfig{}
	if err := strict.Decode(&c); err != nil {
		return c, err
	}
	if c.Version != 1 {
		return c, fmt.Errorf("unsupported fork configuration version %d", c.Version)
	}
	if len(c.Forks) == 0 {
		return c, fmt.Errorf("forks must be a nonempty array")
	}
	seen := map[string]bool{}
	for _, e := range c.Forks {
		for _, name := range []string{e.Fork, e.Upstream} {
			if !repositoryName.MatchString(name) || strings.Contains(name, "..") {
				return c, fmt.Errorf("fork and upstream must be owner/repository names: %q", name)
			}
		}
		if strings.EqualFold(e.Fork, e.Upstream) {
			return c, fmt.Errorf("fork %s names itself as upstream", e.Fork)
		}
		if !branchName.MatchString(e.Branch) || strings.Contains(e.Branch, "..") || strings.HasSuffix(e.Branch, "/") || strings.HasSuffix(e.Branch, ".lock") {
			return c, fmt.Errorf("fork %s needs a plain branch name", e.Fork)
		}
		if e.Pin != "" && !pinRevision.MatchString(e.Pin) {
			return c, fmt.Errorf("fork %s pin must be a lowercase hexadecimal commit ID of at least 7 characters", e.Fork)
		}
		key := strings.ToLower(e.Fork) + "\x00" + e.Branch
		if seen[key] {
			return c, fmt.Errorf("duplicate fork and branch: %s %s", e.Fork, e.Branch)
		}
		seen[key] = true
	}
	return c, nil
}

// The token travels in the environment so it never appears in argv, a URL, or
// a stored remote. It is only attached to HTTPS requests for the fork's host.
func forkGitEnv(base, token string) []string {
	env := append(os.Environ(), "GIT_TERMINAL_PROMPT=0")
	if token == "" || !strings.HasPrefix(base, "https://") {
		return env
	}
	header := "Authorization: Basic " + base64.StdEncoding.EncodeToString([]byte("x-access-token:"+token))
	return append(env, "GIT_CONFIG_COUNT=2",
		"GIT_CONFIG_KEY_0=http."+base+"/.extraheader", "GIT_CONFIG_VALUE_0="+header,
		"GIT_CONFIG_KEY_1=credential.helper", "GIT_CONFIG_VALUE_1=")
}

func forkGit(ctx context.Context, env []string, token, directory string, args ...string) (string, string, error) {
	cmd := exec.CommandContext(ctx, "git", append([]string{"-C", directory}, args...)...)
	cmd.Env = env
	configureProcessGroup(cmd)
	var stdout, stderr bytes.Buffer
	cmd.Stdout, cmd.Stderr = &stdout, &stderr
	err := cmd.Run()
	detail := strings.Join(strings.Fields(stderr.String()), " ")
	if err != nil && ctx.Err() != nil {
		detail = ctx.Err().Error()
	}
	if token != "" {
		detail = strings.ReplaceAll(detail, token, "[redacted]")
		detail = strings.ReplaceAll(detail, base64.StdEncoding.EncodeToString([]byte("x-access-token:"+token)), "[redacted]")
	}
	if len(detail) > 600 {
		detail = detail[:600]
	}
	return strings.TrimSpace(stdout.String()), detail, err
}

// classifyForkFailure maps Git's remote error text to a report status. Text
// that matches nothing stays "failed" so an unknown error is never softened.
func classifyForkFailure(detail string) string {
	text := strings.ToLower(detail)
	for _, marker := range []string{"without `workflow` scope", "without `workflows` permission", "permission to ", "write access to repository not granted", "authentication failed", "could not read username", "invalid username or", "returned error: 403", "returned error: 401"} {
		if strings.Contains(text, marker) {
			return "token-insufficient"
		}
	}
	for _, marker := range []string{"repository not found", "does not appear to be a git repository", "couldn't find remote ref", "returned error: 404", "not found"} {
		if strings.Contains(text, marker) {
			return "missing"
		}
	}
	return "failed"
}

func tokenAdvice(token string) string {
	if token == "" {
		return forkTokenEnv + " is not set. " + tokenHelp
	}
	return tokenHelp
}

func remoteURL(base, name string) string {
	return strings.TrimSuffix(base, "/") + "/" + name + ".git"
}

func syncFork(ctx context.Context, e forkEntry, o forkSyncOptions) (r forkResult) {
	timeout := o.Timeout
	if timeout <= 0 {
		timeout = 10 * time.Minute
	}
	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()
	r = forkResult{Fork: e.Fork, Upstream: e.Upstream, Branch: e.Branch, Status: "failed"}
	work, err := os.MkdirTemp("", "skillforge-fork-")
	if err != nil {
		r.Message = err.Error()
		return r
	}
	defer func() {
		if err := os.RemoveAll(work); err != nil && r.Status != "failed" {
			r.Status, r.Message = "failed", "remove temporary repository: "+err.Error()
		}
	}()
	plain, authorized := forkGitEnv(o.GitBase, ""), forkGitEnv(o.GitBase, o.Token)
	forkURL, upstreamURL, ref := remoteURL(o.GitBase, e.Fork), remoteURL(o.GitBase, e.Upstream), "refs/heads/"+e.Branch
	local := func(args ...string) (string, error) {
		out, detail, err := forkGit(ctx, plain, "", work, args...)
		if err != nil {
			return "", fmt.Errorf("git %s: %w: %s", args[0], err, detail)
		}
		return out, nil
	}
	ancestor := func(older, newer string) (bool, error) {
		_, detail, err := forkGit(ctx, plain, "", work, "merge-base", "--is-ancestor", older, newer)
		if exit, ok := errors.AsType[*exec.ExitError](err); ok && exit.ExitCode() == 1 {
			return false, nil
		}
		if err != nil {
			return false, fmt.Errorf("git merge-base: %w: %s", err, detail)
		}
		return true, nil
	}
	count := func(revisions string) (int, error) {
		out, err := local("rev-list", "--count", revisions)
		if err != nil {
			return 0, err
		}
		return strconv.Atoi(out)
	}
	fail := func(err error) forkResult {
		r.Status, r.Message = "failed", err.Error()
		return r
	}
	if _, err := local("init", "--quiet", "--bare"); err != nil {
		return fail(err)
	}
	if _, detail, err := forkGit(ctx, authorized, o.Token, work, "fetch", "--quiet", "--no-tags", forkURL, "+"+ref+":refs/skillforge/fork"); err != nil {
		r.Status = classifyForkFailure(detail)
		switch r.Status {
		case "missing":
			r.Message = "The fork or its branch was not found, or the token cannot read it: " + detail
		case "token-insufficient":
			r.Message = "The fork could not be read: " + detail + " " + tokenAdvice(o.Token)
		default:
			r.Message = "fetch fork: " + detail
		}
		return r
	}
	if _, detail, err := forkGit(ctx, plain, "", work, "fetch", "--quiet", "--no-tags", upstreamURL, "+"+ref+":refs/skillforge/upstream"); err != nil {
		r.Message = "fetch upstream: " + detail
		return r
	}
	if r.ForkRevision, err = local("rev-parse", "refs/skillforge/fork"); err != nil {
		return fail(err)
	}
	if r.UpstreamRevision, err = local("rev-parse", "refs/skillforge/upstream"); err != nil {
		return fail(err)
	}
	if r.Behind, err = count(r.ForkRevision + ".." + r.UpstreamRevision); err != nil {
		return fail(err)
	}
	if r.Ahead, err = count(r.UpstreamRevision + ".." + r.ForkRevision); err != nil {
		return fail(err)
	}
	final := r.ForkRevision
	switch {
	case r.Ahead > 0:
		r.Status = "diverged"
		r.Message = fmt.Sprintf("The fork has %d commits that upstream lacks and is %d behind. It was not changed.", r.Ahead, r.Behind)
		return r
	case r.Behind == 0:
		r.Status, r.Message = "current", "The fork matches upstream."
	case o.CheckOnly:
		r.Status, r.Message = "behind", fmt.Sprintf("The fork is %d commits behind upstream. Check mode made no changes.", r.Behind)
	default:
		// A plain push: the remote itself refuses anything other than a fast-forward.
		if _, detail, err := forkGit(ctx, authorized, o.Token, work, "push", "--quiet", forkURL, r.UpstreamRevision+":"+ref); err != nil {
			r.Status = classifyForkFailure(detail)
			if r.Status == "missing" {
				r.Status = "failed"
			}
			r.Message = fmt.Sprintf("The fork is %d commits behind and the push was refused: %s", r.Behind, detail)
			if r.Status == "token-insufficient" {
				r.Message += " " + tokenAdvice(o.Token)
			}
			return r
		}
		listed, detail, err := forkGit(ctx, authorized, o.Token, work, "ls-remote", forkURL, ref)
		if err != nil {
			r.Message = "The push returned success but the fork could not be read back: " + detail
			return r
		}
		if now, _, _ := strings.Cut(listed, "\t"); now != r.UpstreamRevision {
			r.Message = fmt.Sprintf("The push returned success but the fork is at %q, expected %s.", now, r.UpstreamRevision)
			return r
		}
		final = r.UpstreamRevision
		r.Status, r.Message = "synced", fmt.Sprintf("Fast-forwarded %d commits. The fork was read back at %s.", r.Behind, final)
	}
	if e.Pin == "" {
		return r
	}
	r.Pin = &pinResult{Revision: e.Pin}
	pin, err := local("rev-parse", "--verify", "--quiet", e.Pin+"^{commit}")
	if err == nil {
		if r.Pin.InFork, err = ancestor(pin, final); err != nil {
			return fail(err)
		}
		if r.Pin.UpstreamCommitsPast, err = count(pin + ".." + r.UpstreamRevision); err != nil {
			return fail(err)
		}
	}
	if !r.Pin.InFork && r.Status != "behind" {
		r.Status = "pin-missing"
		r.Message += fmt.Sprintf(" The pinned commit %s is not in the fork branch.", e.Pin)
	} else if r.Pin.UpstreamCommitsPast > 0 {
		r.Message += fmt.Sprintf(" Upstream is %d commits past the pin %s.", r.Pin.UpstreamCommitsPast, e.Pin)
	}
	return r
}

func forkStatus(results []forkResult) string {
	status := "pass"
	for _, r := range results {
		switch r.Status {
		case "current", "synced":
		case "diverged", "behind", "pin-missing":
			if status == "pass" {
				status = "fail"
			}
		default:
			status = "error"
		}
	}
	return status
}

func forkSummary(report forkReport) string {
	var b strings.Builder
	fmt.Fprintf(&b, "## Fork sync: %s\n\n| Fork | Upstream | Branch | Status | Behind | Ahead | Detail |\n| --- | --- | --- | --- | --- | --- | --- |\n", report.Status)
	for _, r := range report.Forks {
		fmt.Fprintf(&b, "| %s | %s | %s | %s | %d | %d | %s |\n", r.Fork, r.Upstream, r.Branch, r.Status, r.Behind, r.Ahead, strings.ReplaceAll(r.Message, "|", "\\|"))
	}
	return b.String()
}

// postWebhook never includes the URL in an error: webhook URLs are credentials.
func postWebhook(ctx context.Context, client *http.Client, target string, report forkReport) error {
	parsed, err := url.Parse(target)
	if err != nil || parsed.Host == "" || (parsed.Scheme != "https" && parsed.Scheme != "http") {
		return errors.New("the webhook URL must be an http or https URL")
	}
	text := "Skillforge fork sync: " + report.Status
	for _, r := range report.Forks {
		if r.Status != "current" {
			text += fmt.Sprintf("\n%s (%s): %s. %s", r.Fork, r.Branch, r.Status, r.Message)
		}
	}
	body, err := json.Marshal(map[string]any{"text": text, "report": report})
	if err != nil {
		return err
	}
	ctx, cancel := context.WithTimeout(ctx, 20*time.Second)
	defer cancel()
	request, err := http.NewRequestWithContext(ctx, http.MethodPost, target, bytes.NewReader(body))
	if err != nil {
		return errors.New("the webhook request could not be built")
	}
	request.Header.Set("Content-Type", "application/json")
	// Redirects can disclose the credential-bearing URL in nested client errors.
	// Treat every redirect as a delivery failure and never return transport text.
	safeClient := *client
	safeClient.CheckRedirect = func(*http.Request, []*http.Request) error { return http.ErrUseLastResponse }
	response, err := safeClient.Do(request)
	if err != nil {
		return errors.New("webhook delivery failed (connection, TLS, timeout, or redirect error)")
	}
	defer response.Body.Close()
	if _, err := io.Copy(io.Discard, io.LimitReader(response.Body, 1<<16)); err != nil {
		return errors.New("webhook response could not be read")
	}
	if response.StatusCode < 200 || response.StatusCode > 299 {
		return fmt.Errorf("webhook returned HTTP %d", response.StatusCode)
	}
	return nil
}

func runForkSync(ctx context.Context, o forkSyncOptions, out io.Writer) error {
	config, err := readForkConfig(o.Directory)
	if err != nil {
		return fmt.Errorf("%s: %w", forkConfigName, err)
	}
	if _, err := exec.LookPath("git"); err != nil {
		return err
	}
	report := forkReport{Version: 1, CheckedAt: time.Now().UTC().Format(time.RFC3339), Mode: "sync", Forks: []forkResult{}, Notification: forkNotification{Summary: "not-configured", Webhook: "not-configured"}}
	if o.CheckOnly {
		report.Mode = "check"
	}
	quiet := true
	for _, e := range config.Forks {
		r := syncFork(ctx, e, o)
		report.Forks = append(report.Forks, r)
		quiet = quiet && r.Status == "current"
		fmt.Fprintf(out, "%s: %s from %s (%s): %s\n", r.Status, r.Fork, r.Upstream, r.Branch, r.Message)
	}
	report.Status = forkStatus(report.Forks)
	var problems []string
	if o.Webhook != "" {
		report.Notification.Webhook = "not-needed"
		if !quiet {
			report.Notification.Webhook = "delivered"
			if err := postWebhook(ctx, o.Client, o.Webhook, report); err != nil {
				report.Notification.Webhook = "failed"
				problems = append(problems, err.Error())
			}
		}
	}
	if o.SummaryPath != "" {
		report.Notification.Summary = "written"
		if err := appendFile(o.SummaryPath, forkSummary(report)); err != nil {
			report.Notification.Summary = "failed"
			problems = append(problems, "job summary: "+err.Error())
		}
	}
	if len(problems) > 0 {
		report.Status = "error"
		report.Notification.Message = strings.Join(problems, "; ")
		fmt.Fprintf(out, "notification failed: %s\n", report.Notification.Message)
	}
	if o.Output != "" {
		if err := writeJSON(o.Output, report); err != nil {
			return err
		}
	}
	fmt.Fprintf(out, "Fork sync: %s\n", report.Status)
	switch report.Status {
	case "pass":
		return nil
	case "fail":
		return errFindings
	}
	return errForkSync
}

func appendFile(path, text string) (err error) {
	f, err := os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0600)
	if err != nil {
		return err
	}
	defer func() { err = errors.Join(err, f.Close()) }()
	_, err = f.WriteString(text)
	return err
}
