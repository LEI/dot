package main

import (
	"fmt"

	"github.com/LEI/dot/internal/dot"
)

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
		if dotOptions.Verbose > 0 {
			fmt.Println("#", t.DoString())
		}
		return nil
	}
	fmt.Println("$", t.DoString())
	if dotOptions.DryRun {
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
		// if dotOptions.Verbose > 0 {
		// 	fmt.Println("#", t.UndoString())
		// }
		return nil
	}
	fmt.Println("$", t.UndoString())
	if dotOptions.DryRun {
		return nil
	}
	return t.Undo()
}
