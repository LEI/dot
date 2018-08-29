package main

import (
	"fmt"
	"sync"

	"github.com/LEI/dot/internal/dot"
	"github.com/spf13/cobra"
)

func preRunTask(cmd *cobra.Command, args []string) error {
	action := cmd.Parent().Name()
	switch action {
	// case "list":
	// 	return preRunList(cmd, args)
	case "sync":
		return nil
	case "install":
		return preRunInstall(cmd, args)
	case "remove":
		return preRunRemove(cmd, args)
	default:
		return fmt.Errorf("%s: invalid action", action)
	}
}

type actionResult struct {
	// role *dot.Role
	name string
	task dot.Tasker
	// out  string
	err error
}

// Run after preRunInstall and preRunRemove
func preRunAction(cmd *cobra.Command, args []string) error {
	action := cmd.Name()
	c := make(chan actionResult)
	ignoreErrors := action == "list"
	roles := dotConfig.Roles
	go func() {
		var wg sync.WaitGroup
		for _, r := range roles {
			wg.Add(1)
			go checkAllTasks(action, r, c, &wg)
		}
		// All calls to wg.Add are done. Start a goroutine
		// to close c once all the sends are done.
		go func() {
			wg.Wait()
			close(c)
		}()
	}()
	// Check all tasks result
	exists := 0
	failed := 0
	skipped := 0
	total := 0
	for r := range c {
		total++
		if r.err == nil {
			continue
		}
		if dot.IsExist(r.err) {
			exists++
			continue
		}
		if dot.IsSkip(r.err) {
			skipped++
			continue
		}
		if !ignoreErrors {
			fmt.Fprintf(dotOpts.stderr, "failed to %s %s role: %s\n", action, r.name, r.err)
		}
		failed++
	}
	// if total == exists+skipped && !dotOpts.Force {
	// 	return &appError{
	// 		Err:  nil, // dot.ErrExist
	// 		Msg:  "nothing to do",
	// 		Code: 0,
	// 	}
	// }
	if failed > 0 && !ignoreErrors {
		return fmt.Errorf("%d error(s) while checking %d roles", failed, len(roles))
	}
	return nil
}

func checkAllTasks(action string, r *dot.Role, c chan<- actionResult, wg *sync.WaitGroup) {
	// var wg sync.WaitGroup
	// if dotOpts.verbosity >= 1 {
	// 	fmt.Fprintf(dotOpts.stdout, "## Checking %s...\n", r.Name)
	// }
	if dotOpts.pkg {
		for _, p := range r.Pkgs {
			wg.Add(1)
			go checkOneTask(action, r, p, c, wg)
		}
	}
	// wg.Add(len(r.Dirs))
	for _, d := range r.Dirs {
		wg.Add(1)
		go checkOneTask(action, r, d, c, wg)
	}
	for _, f := range r.Files {
		wg.Add(1)
		go checkOneTask(action, r, f, c, wg)
	}
	for _, l := range r.Links {
		wg.Add(1)
		go checkOneTask(action, r, l, c, wg)
	}
	for _, t := range r.Tpls {
		wg.Add(1)
		go checkOneTask(action, r, t, c, wg)
	}
	for _, l := range r.Lines {
		wg.Add(1)
		go checkOneTask(action, r, l, c, wg)
	}
	checkTaskHooks(action, r, c, wg)
	wg.Done()
}

func checkTaskHooks(action string, r *dot.Role, c chan<- actionResult, wg *sync.WaitGroup) {
	for _, h := range r.Install {
		wg.Add(1)
		go checkOneTask(action, r, h, c, wg)
	}
	for _, h := range r.PostInstall {
		wg.Add(1)
		go checkOneTask(action, r, h, c, wg)
	}
	for _, h := range r.Remove {
		wg.Add(1)
		go checkOneTask(action, r, h, c, wg)
	}
	for _, h := range r.PostRemove {
		wg.Add(1)
		go checkOneTask(action, r, h, c, wg)
	}
}

func checkOneTask(action string, r *dot.Role, t dot.Tasker, c chan<- actionResult, wg *sync.WaitGroup) {
	t.SetAction(action)
	if err := t.Check(); err != nil {
		c <- actionResult{r.Name, t, err}
		wg.Done()
		return
	}
	err := t.Status()
	c <- actionResult{r.Name, t, err}
	wg.Done()
}

// func runTask(action string, i interface{}) error {
// 	t := i.(dot.Tasker)
func runTask(action string, t dot.Tasker) error {
	//t.SetAction(action)
	switch action {
	case "install":
		if err := doTask(t); err != nil && !dot.IsSkip(err) {
			return err // fmt.Errorf("%s task: %s", action, err)
		}
	case "remove":
		if err := undoTask(t); err != nil && !dot.IsSkip(err) {
			return err // fmt.Errorf("%s task: %s", action, err)
		}
	default:
		return fmt.Errorf("%s: unknown action", action)
	}
	return nil
}

func doTask(t dot.Tasker) error {
	if err := t.Check(); err != nil {
		return err
	}
	err := t.Status()
	ok := dot.IsExist(err)
	if !ok && err != nil {
		return err
	}
	str := t.String() // fmt.Sprintf("%s", t)
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
	if err := t.Check(); err != nil {
		return err
	}
	err := t.Status()
	ok := dot.IsExist(err)
	if !ok && err != nil {
		return err
	}
	str := t.String() // fmt.Sprintf("%s", t)
	// if str == "" {
	// 	fmt.Fprintln(os.Stderr, "warning: empty task string")
	// }
	if !ok {
		if str != "" && dotOpts.verbosity >= 2 {
			fmt.Fprintf(dotOpts.stdout, "# %s\n", str)
		}
		// if !shell.AskConfirmation("Proceed to execute?") {
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
