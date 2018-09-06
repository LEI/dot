package main

import (
	"fmt"
	"sync"

	"github.com/LEI/dot/internal/dot"
	"github.com/spf13/cobra"
)

var cmdSync = &cobra.Command{
	Use:     "sync",
	Aliases: []string{"s"},
	Short:   "Synchronize roles",
	Long:    `The "sync" command clone or pull a role repository.`,
	Args:    cobra.NoArgs,
	PreRunE: preRunSync,
	RunE:    runSync,
	// DisableAutoGenTag: true,
}

func init() {
	cmdRoot.AddCommand(cmdSync)

	// addActionFlags(cmdSync)
}

// https://golang.org/pkg/os/exec/#Cmd.StdoutPipe

func preRunSync(cmd *cobra.Command, args []string) error {
	return nil // preRunAction(cmd, args)
}

func runSync(cmd *cobra.Command, args []string) error {
	roles := dotConfig.Roles
	ch := syncRoles(roles)
	errs := []error{} // failed := false
	for r := range ch {
		// fmt.Printf("# Synced %s:\n", r.role.Name)
		if r.out != "" {
			fmt.Fprintln(dotOpts.stdout, r.out)
		}
		if r.err != nil {
			fmt.Fprintln(dotOpts.stderr, r.err)
			errs = append(errs, r.err) // failed = true
		}
	}
	if len(errs) > 0 { // failed
		return fmt.Errorf(
			"%d error(s) occurred while syncing %d roles",
			len(errs),
			len(roles),
		)
	}
	return nil
}

// https://blog.golang.org/pipelines#TOC_8.
type syncResult struct {
	role *dot.Role
	out  string
	err  error
}

// done <-chan struct{}
func syncRoles(roles []*dot.Role) <-chan syncResult {
	ch := make(chan syncResult)
	go func() {
		var wg sync.WaitGroup
		for _, r := range roles {
			if !r.ShouldRun() {
				continue
			}
			wg.Add(1)
			go func(r *dot.Role) {
				out, err := r.Sync()
				if err != nil {
					// Git error
					ch <- syncResult{r, out, err}
					wg.Done()
					return
				}
				if err := r.Load(); err != nil {
					// Cfg error
					ch <- syncResult{r, out, err}
					wg.Done()
					return
				}
				// fmt.Println("Loaded", r.Name)
				// select {
				// case ch <- syncResult{r, "ok", err}:
				// }
				ch <- syncResult{r, out, nil}
				wg.Done()
			}(r)
		}
		// All calls to wg.Add are done. Start a goroutine
		// to close ch once all the sends are done.
		go func() {
			wg.Wait()
			close(ch)
		}()
	}()
	return ch
}
