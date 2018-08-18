package main

import (
	"fmt"
	"os"

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
	switch action {
	case "install":
		if err := doTask(t); err != nil {
			return err
		}
	case "remove":
		if err := undoTask(t); err != nil {
			return err
		}
	default:
		return fmt.Errorf("%s: unknown action", action)
	}
	return nil
}

func doTask(t dot.Tasker) error {
	err := t.Status()
	ok := dot.IsOk(err)
	if !ok && err != nil {
		return err
	}
	str := t.DoString()
	if str == "" {
		fmt.Fprintln(os.Stderr, "warning: empty task string")
	}
	if ok {
		if str != "" && dotOpts.verbosity >= 2 {
			fmt.Printf("# %s\n", str)
		}
		return nil
	}
	if str != "" && dotOpts.verbosity >= 1 {
		fmt.Printf("$ %s\n", str)
	}
	if dotOpts.DryRun {
		return nil
	}
	return t.Do()
}

func undoTask(t dot.Tasker) error {
	err := t.Status()
	ok := dot.IsOk(err)
	if !ok && err != nil {
		return err
	}
	str := t.UndoString()
	if str == "" {
		fmt.Fprintln(os.Stderr, "warning: empty task string")
	}
	if !ok {
		fmt.Println("# TODO AskConfirmation")
		if str != "" && dotOpts.verbosity >= 2 {
			fmt.Printf("# %s\n", str)
		}
		return nil
	}
	if str != "" && dotOpts.verbosity >= 1 {
		fmt.Printf("$ %s\n", str)
	}
	if dotOpts.DryRun {
		return nil
	}
	return t.Undo()
}