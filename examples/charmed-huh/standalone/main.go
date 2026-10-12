// examples/charmed-huh/standalone/main.go
// A standalone form that collects a local task plan and prints the choices.
package main

import (
	"errors"
	"fmt"
	"os"
	"strings"

	"charm.land/huh/v2"
)

func main() {
	var (
		projects []string
		tasks    []string
		location string
		custom   string
		notes    string
		confirm  bool
	)

	form := huh.NewForm(
		huh.NewGroup(
			huh.NewNote().
				Title("Local task planner").
				Description("Choose projects and tasks. Use _space_ to toggle, `/` to filter, *ctrl+a* for all."),

			huh.NewMultiSelect[string]().
				Title("Projects").
				Options(
					huh.NewOption("Website", "site").Selected(true),
					huh.NewOption("API", "api"),
					huh.NewOption("Documentation", "docs"),
					huh.NewOption("Worker", "worker"),
					huh.NewOption("CLI", "cli"),
					huh.NewOption("Tests", "tests"),
				).
				Filterable(true).
				Height(8).
				Validate(func(v []string) error {
					if len(v) == 0 {
						return errors.New("pick at least one project")
					}
					return nil
				}).
				Value(&projects),

			huh.NewMultiSelect[string]().
				Title("Tasks").
				Options(huh.NewOptions("build", "test", "lint", "format")...).
				Limit(4).
				Value(&tasks),
		).Title("What"),

		huh.NewGroup(
			huh.NewSelect[string]().
				Title("Working directory").
				Description("Choose a local directory for the task plan.").
				Options(
					huh.NewOption("Current directory", "."),
					huh.NewOption("Parent directory", ".."),
					huh.NewOption("Custom path", "custom"),
				).
				Value(&location),

			huh.NewInput().
				Title("Custom path").
				Placeholder("./workspace").
				Suggestions([]string{"./workspace", "../workspace"}).
				Validate(func(s string) error {
					if location == "custom" && strings.TrimSpace(s) == "" {
						return errors.New("enter a directory")
					}
					return nil
				}).
				Value(&custom),

			huh.NewText().
				Title("Notes").
				Placeholder("optional").
				Lines(3).
				CharLimit(200).
				Value(&notes),
		).Title("Where"),

		huh.NewGroup(
			huh.NewConfirm().
				Title("Print task plan?").
				Affirmative("Print").
				Negative("Cancel").
				Value(&confirm),
		),
	).WithShowHelp(true)

	err := form.Run()
	switch {
	case errors.Is(err, huh.ErrUserAborted):
		fmt.Fprintln(os.Stderr, "cancelled")
		os.Exit(130)
	case err != nil:
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	if !confirm {
		fmt.Println("no task plan printed")
		return
	}
	fmt.Printf("planned projects=%v tasks=%v directory=%s custom=%s\n", projects, tasks, location, custom)
}
