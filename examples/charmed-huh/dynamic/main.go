// examples/charmed-huh/dynamic/main.go
// Fields that recompute from earlier answers. The *Func variants take a
// binding (a pointer to the value they depend on); huh re-evaluates the
// function whenever the bound value changes.
package main

import (
	"fmt"
	"os"

	"charm.land/huh/v2"
)

var tasksByProject = map[string][]string{
	"site":   {"preview", "build", "test"},
	"docs":   {"preview", "check links"},
	"worker": {"test"},
}

func main() {
	var project, task string
	var details bool

	err := huh.NewForm(
		huh.NewGroup(
			huh.NewSelect[string]().
				Title("Project").
				Options(
					huh.NewOption("Website", "site"),
					huh.NewOption("Documentation", "docs"),
					huh.NewOption("Worker", "worker"),
				).
				Value(&project),

			huh.NewSelect[string]().
				TitleFunc(func() string { return "Task for " + project }, &project).
				OptionsFunc(func() []huh.Option[string] {
					return huh.NewOptions(tasksByProject[project]...)
				}, &project).
				Value(&task),

			huh.NewConfirm().
				Title("Show task details?").
				Value(&details),
		),
		huh.NewGroup(
			huh.NewNote().
				TitleFunc(func() string { return "Details for " + project }, &project).
				Description("This example only prints your choices; it does not run tasks."),
		).WithHideFunc(func() bool { return !details }), // skipped entirely when false
	).Run()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	fmt.Printf("%s %s (details: %v)\n", project, task, details)
}
