package main

import (
	"fmt"
	"os"
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

// https://blog.golang.org/pipelines#TOC_8.
type syncResult struct {
	role *dot.Role
	out  string
	err  error
}

func preRunSync(cmd *cobra.Command, args []string) error {
	return nil // preRunAction(cmd, args)
}

func runSync(cmd *cobra.Command, args []string) error {
	c := syncRoles(dotConfig.Roles)
	failed := false
	for r := range c {
		// fmt.Printf("# Synced %s:\n", r.role.Name)
		if r.out != "" {
			fmt.Fprintln(dotOpts.stdout, r.out)
		}
		if r.err != nil {
			fmt.Fprintln(os.Stderr, "ERR", r.err)
			failed = true
		}
	}
	if failed {
		return fmt.Errorf("failed to sync roles")
	}
	return nil
}

// done <-chan struct{}
func syncRoles(roles []*dot.Role) <-chan syncResult {
	c := make(chan syncResult)
	go func() {
		var wg sync.WaitGroup
		for _, r := range roles {
			wg.Add(1)
			go func(r *dot.Role) {
				out, err := r.Sync()
				if err != nil {
					c <- syncResult{r, out, err}
					wg.Done()
					return
				}
				if err := r.Load(); err != nil {
					c <- syncResult{r, out, err}
					wg.Done()
					return
				}
				// fmt.Println("Loaded", r.Name)
				// select {
				// case c <- syncResult{r, "ok", err}:
				// }
				c <- syncResult{r, out, nil}
				wg.Done()
			}(r)
		}
		// All calls to wg.Add are done. Start a goroutine
		// to close c once all the sends are done.
		go func() {
			wg.Wait()
			close(c)
		}()
	}()
	return c
}
