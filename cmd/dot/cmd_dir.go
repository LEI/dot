package main

import (
	"fmt"
	"path/filepath"

	"github.com/LEI/dot/internal/dot"
	"github.com/spf13/cobra"
)

var cmdDir = &cobra.Command{
	Use:     "dir [flags]",
	Aliases: []string{"d"},
	Short:   "Dir managed files",
	Long: `
The "dir" task installs or removes directories.
`,
	DisableAutoGenTag: true,
	Args:              cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		return runDir(cmd.Parent().Name(), dirOptions, globalOptions, args)
	},
}

// DirOptions collects all options for the dir command.
type DirOptions struct {
	// Empty bool
}

var dirOptions DirOptions

func init() {
	cmdInstall.AddCommand(cmdDir)
	cmdRemove.AddCommand(cmdDir)

	// flags := cmdDir.Flags()
}

func runDir(action string, opts DirOptions, gopts GlobalOptions, args []string) error {
	for _, r := range globalConfig.Roles {
		for _, d := range r.Dirs {
			switch action {
			case "install":
				if err := doDir(d, opts, gopts); err != nil {
					return err
				}
			case "remove":
				if err := undoDir(d, opts, gopts); err != nil {
					return err
				}
			}
		}
	}
	return nil
}

func doDir(d *dot.Dir, opts DirOptions, gopts GlobalOptions) error {
	if !filepath.IsAbs(d.Path) {
		d.Path = filepath.Join(gopts.Target, d.Path)
	}
	err := d.Status()
	ok := dot.IsOk(err)
	if !ok && err != nil {
		return err
	}
	if ok {
		if gopts.Verbose > 0 {
			fmt.Println("#", d.DoString())
		}
		return nil
	}
	fmt.Println("$", d.DoString())
	return doTask(d, gopts)
}

func undoDir(d *dot.Dir, opts DirOptions, gopts GlobalOptions) error {
	if !filepath.IsAbs(d.Path) {
		d.Path = filepath.Join(gopts.Target, d.Path)
	}
	err := d.Status()
	ok := dot.IsOk(err)
	if !ok && err != nil {
		return err
	}
	if !ok {
		fmt.Println("# TODO AskConfirmation", d.UndoString())
		// if gopts.Verbose > 0 {
		// 	fmt.Println("#", d.UndoString())
		// }
		return nil
	}
	fmt.Println("$", d.UndoString())
	return undoTask(d, gopts)
}

func doTask(i interface{}, gopts GlobalOptions) error {
	t := i.(dot.Tasker)
	if gopts.DryRun {
		return nil
	}
	return t.Do()
}

func undoTask(i interface{}, gopts GlobalOptions) error {
	t := i.(dot.Tasker)
	if gopts.DryRun {
		return nil
	}
	return t.Undo()
}
