package main

import (
	"fmt"

	"github.com/LEI/dot/internal/dot"
	"github.com/spf13/cobra"
)

func preRunTask(cmd *cobra.Command, args []string) error {
	action := cmd.Parent().Name()
	switch action {
	// case "list":
	// 	return preRunList(cmd, args)
	case "install":
		return preRunInstall(cmd, args)
	case "remove":
		return preRunRemove(cmd, args)
	default:
		return fmt.Errorf("%s: invalid action", action)
	}
}

func runTask(action string, i interface{}) error {
	t := i.(dot.Tasker)
	// act := action + " " + t.Type()
	switch action {
	case "install":
		if err := doTask(t); err != nil && !dot.IsSkip(err) {
			// act := t.String()
			return err // fmt.Errorf("%s: %s", act, err)
		}
	case "remove":
		if err := undoTask(t); err != nil && !dot.IsSkip(err) {
			// act := t.UndoString()
			return err // fmt.Errorf("%s: %s", act, err)
		}
	default:
		return fmt.Errorf("%s: unknown action", action)
	}
	return nil
}

func doTask(t dot.Tasker) error {
	t.SetAction("install")
	if err := t.Check(); err != nil {
		return err
	}
	err := t.Status()
	ok := dot.IsExist(err)
	if !ok && err != nil {
		return err
	}
	str := t.String()
	// if str == "" {
	// 	fmt.Fprintln(os.Stderr, "warning: empty task string")
	// }
	if ok {
		if str != "" && dotOpts.verbosity >= 2 {
			fmt.Fprintf(dotOpts.stdout, "# %s\n", str)
		}
		return nil
	}
	if str != "" && dotOpts.verbosity >= 1 {
		fmt.Fprintf(dotOpts.stdout, "$ %s\n", str)
	}
	if dotOpts.DryRun {
		return nil
	}
	return t.Do()
}

func undoTask(t dot.Tasker) error {
	t.SetAction("remove")
	if err := t.Check(); err != nil {
		return err
	}
	err := t.Status()
	ok := dot.IsExist(err)
	if !ok && err != nil {
		return err
	}
	str := t.String()
	// if str == "" {
	// 	fmt.Fprintln(os.Stderr, "warning: empty task string")
	// }
	if !ok {
		if str != "" && dotOpts.verbosity >= 2 {
			fmt.Fprintf(dotOpts.stdout, "# %s\n", str)
		}
		// if !prompt.AskConfirmation("Proceed to execute?") {
		return nil
		// }
	}
	if str != "" && dotOpts.verbosity >= 1 {
		fmt.Fprintf(dotOpts.stdout, "$ %s\n", str)
	}
	if dotOpts.DryRun {
		return nil
	}
	return t.Undo()
}
