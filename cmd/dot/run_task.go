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

func checkTask(action, name string, t dot.Tasker, c chan<- actionResult, wg *sync.WaitGroup) {
	t.SetAction(action)
	// err := t.Check()
	if err := t.Check(); err != nil {
		c <- actionResult{name, t, err}
		wg.Done()
		return
	}
	err := t.Status()
	c <- actionResult{name, t, err}
	wg.Done()
}

// Run after preRunInstall and preRunRemove
func preRunAction(cmd *cobra.Command, args []string) error {
	action := cmd.Name()
	c := make(chan actionResult)
	roles := dotConfig.Roles
	go func() {
		var wg sync.WaitGroup
		for _, r := range roles {
			// if dotOpts.verbosity >= 1 {
			// 	fmt.Fprintf(dotOpts.stdout, "## Checking %s...\n", r.Name)
			// }
			if dotOpts.pkg {
				for _, p := range r.Pkgs {
					wg.Add(1)
					go checkTask(action, r.Name, p, c, &wg)
				}
			}
			for _, d := range r.Dirs {
				wg.Add(1)
				go checkTask(action, r.Name, d, c, &wg)
			}
			for _, f := range r.Files {
				wg.Add(1)
				go checkTask(action, r.Name, f, c, &wg)
			}
			for _, l := range r.Links {
				wg.Add(1)
				go checkTask(action, r.Name, l, c, &wg)
			}
			for _, t := range r.Tpls {
				wg.Add(1)
				go checkTask(action, r.Name, t, c, &wg)
			}
			for _, l := range r.Lines {
				wg.Add(1)
				go checkTask(action, r.Name, l, c, &wg)
			}
			for _, h := range r.Install {
				wg.Add(1)
				go checkTask(action, r.Name, h, c, &wg)
			}
			for _, h := range r.PostInstall {
				wg.Add(1)
				go checkTask(action, r.Name, h, c, &wg)
			}
			for _, h := range r.Remove {
				wg.Add(1)
				go checkTask(action, r.Name, h, c, &wg)
			}
			for _, h := range r.PostRemove {
				wg.Add(1)
				go checkTask(action, r.Name, h, c, &wg)
			}
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
		fmt.Fprintf(dotOpts.stderr, "# error in %s role: %s\n", r.name, r.err)
		failed++
	}
	// if total == exists+skipped && !dotOpts.Force {
	// 	return &appError{
	// 		Err:  nil, // dot.ErrExist
	// 		Msg:  "nothing to do",
	// 		Code: 0,
	// 	}
	// }
	if failed > 0 && action != "list" {
		return fmt.Errorf("%d error(s) while checking %d roles", failed, len(roles))
	}
	return nil
}

/* func checkTask(action, name string, i interface{}) error {
	t, ok := i.(dot.Tasker)
	if !ok {
		return fmt.Errorf("%s: not a tasker", i)
	}
	t.SetAction(action)
	if err := t.Check(); err != nil && !dot.IsSkip(err) {
		return err
	}
	if err := t.Status(); err != nil && !dot.IsExist(err) {
		return err
	}
	return nil
}

// Run after preRunInstall and preRunRemove
func preRunAction(cmd *cobra.Command, args []string) error {
	action := cmd.Name()
	roles := dotConfig.Roles
	for _, r := range roles {
		// if dotOpts.verbosity >= 1 {
		// 	fmt.Fprintf(dotOpts.stdout, "## Checking %s...\n", r.Name)
		// }
		if dotOpts.pkg { // action != "list"
			for _, p := range r.Pkgs {
				if err := checkTask(action, r.Name, p); err != nil {
					return err
				}
			}
		}
		for _, d := range r.Dirs {
			if err := checkTask(action, r.Name, d); err != nil {
				return err
			}
		}
		for _, f := range r.Files {
			if err := checkTask(action, r.Name, f); err != nil {
				return err
			}
		}
		for _, l := range r.Links {
			if err := checkTask(action, r.Name, l); err != nil {
				return err
			}
		}
		for _, t := range r.Tpls {
			if err := checkTask(action, r.Name, t); err != nil {
				return err
			}
		}
		for _, l := range r.Lines {
			if err := checkTask(action, r.Name, l); err != nil {
				return err
			}
		}
		for _, h := range r.Install {
			if err := checkTask(action, r.Name, h); err != nil {
				return err
			}
		}
		for _, h := range r.PostInstall {
			if err := checkTask(action, r.Name, h); err != nil {
				return err
			}
		}
		for _, h := range r.Remove {
			if err := checkTask(action, r.Name, h); err != nil {
				return err
			}
		}
		for _, h := range r.PostRemove {
			if err := checkTask(action, r.Name, h); err != nil {
				return err
			}
		}
	}
	return nil
} */

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
