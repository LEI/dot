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
	Short:   "Synchronize a role",
	Long: `
The "sync" command clone or pull a role repository.
`,
	Args: cobra.NoArgs,
	RunE: runSync,
	// DisableAutoGenTag: true,
}

func init() {
	cmdRoot.AddCommand(cmdSync)
}

// https://golang.org/pkg/os/exec/#Cmd.StdoutPipe

// https://blog.golang.org/pipelines#TOC_8.
type syncResult struct {
	role *dot.Role
	out  string
	err  error
}

func runSync(cmd *cobra.Command, args []string) error {
	c := syncRoles(dotConfig.Roles)
	for r := range c {
		// fmt.Printf("# Synced %s:\n", r.role.Name)
		fmt.Println(r.out)
		if r.err != nil {
			fmt.Fprintln(os.Stderr, "ERR", r.err)
		}
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
					return
				}
				if err := r.Load(); err != nil {
					c <- syncResult{r, out, err}
					return
				}
				c <- syncResult{r, out, nil}
				// select {
				// case c <- syncResult{r, "ok", err}:
				// }
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
