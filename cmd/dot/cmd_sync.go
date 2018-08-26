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
	DisableAutoGenTag: true,
	Args:              cobra.NoArgs,
	RunE:              runSync,
}

func init() {
	cmdRoot.AddCommand(cmdSync)
}

// func runSync(cmd *cobra.Command, args []string) error {
// 	// dotOpts.stdout
// 	// fmt.Print(term.ClearEntireScreen())
// 	fmt.Println("Syncing roles...")
// 	l := len(dotConfig.Roles)
// 	for i, r := range dotConfig.Roles {
// 		// var buf bytes.Buffer
// 		// reader := io.MultiReader(os.Stdout, os.Stderr)
// 		// w := io.MultiWriter(os.Stdout, &buf)
// 		// git.Stdout = w
// 		// git.Stderr = ?
// 		fmt.Printf("... Syncing %s (%d/%d) ...\n", r.Name, i+1, l)
// 		// https://golang.org/pkg/os/exec/#Cmd.StdoutPipe
// 		// Clone or pull git repository
// 		if err := r.Sync(); err != nil {
// 			return err
// 		}
// 		// Parse config file (again)
// 		if err := r.LoadConfig(); err != nil {
// 			return err
// 		}
// 		// lc := len(strings.Split(buf.String(), "\n")) + 2
// 		// nl := bytes.Count(buf.Bytes(), []byte{'\n'})
// 		// fmt.Print(term.MoveUp(nl))
// 		// fmt.Print(term.ClearScreenDown())
// 		// fmt.Printf("--- Synced %s ---\n%v\n---\n", r.Name, nl)
// 	}
// 	fmt.Println("All done!")
// 	return nil
// }

/*func runSync(cmd *cobra.Command, args []string) error {
	// dotOpts.stdout
	length := len(dotConfig.Roles)
	errs := make(chan error, length)
	for i, r := range dotConfig.Roles {
		// var buf bytes.Buffer
		// reader := io.MultiReader(os.Stdout, os.Stderr)
		// w := io.MultiWriter(os.Stdout, &buf)
		// git.Stdout = w
		// git.Stderr = ?
		// https://golang.org/pkg/os/exec/#Cmd.StdoutPipe
		go func(i int, r *dot.Role) {
			// n := fmt.Sprintf("%d/%d", i+1, length)
			// fmt.Printf("Syncing %s (%s) ...\n", r.Name, n)
			//buf := bytes.Buffer{}
			//git.Stdout = &buf
			//git.Stderr = &buf
			// Clone or pull git repository
			if err := r.Sync(); err != nil {
				errs <- err
				return
			}
			//fmt.Printf("Sync %s:\n[%s]\n", r.Path, strings.TrimSuffix(buf.String(), "\n"))
			// Parse config file (again)
			if err := r.Load(); err != nil {
				errs <- err
				return
			}
			// lc := len(strings.Split(buf.String(), "\n")) + 2
			// nl := bytes.Count(buf.Bytes(), []byte{'\n'})
			// fmt.Print(term.MoveUp(nl))
			// fmt.Print(term.ClearScreenDown())
			// fmt.Printf("--- Synced %s ---\n%v\n---\n", r.Name, nl)
			errs <- nil
		}(i, r)
	}
	//git.Stdout = os.Stdout
	//git.Stderr = os.Stderr
	for i := 0; i < length; i++ {
		if err := <-errs; err != nil {
			return err
		}
	}
	return nil
}*/

// https://blog.golang.org/pipelines#TOC_8.
type syncResult struct {
	role *dot.Role
	out  string
	err  error
}

func runSync(cmd *cobra.Command, args []string) error {
	c := syncRoles(dotConfig.Roles)
	for r := range c {
		fmt.Printf("-> Syncing %s...\n%s\n", r.role.Name, r.out)
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
