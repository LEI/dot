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
	if ok {
		if dotOpts.Verbose > 0 {
			fmt.Println("#", t.DoString())
		}
		return nil
	}
	fmt.Println("$", t.DoString())
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
	if !ok {
		fmt.Println("# TODO AskConfirmation", t.UndoString())
		// if dotOpts.Verbose > 0 {
		// 	fmt.Println("#", t.UndoString())
		// }
		return nil
	}
	fmt.Println("$", t.UndoString())
	if dotOpts.DryRun {
		return nil
	}
	return t.Undo()
}
