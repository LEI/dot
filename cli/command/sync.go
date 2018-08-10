package command

// https://github.com/moby/moby/tree/17.05.x/cli/command/image

import (
	"fmt"
	"github.com/spf13/cobra"

	"github.com/LEI/dot/cli"
	"github.com/LEI/dot/pkg/git"
	"github.com/LEI/dot/system"
)

type syncOptions struct {
	// check bool
	// filter []string
	// force bool
}

// NewSyncCommand returns a cobra command for `sync` subcommands
func NewSyncCommand(dotCli *DotCli) *cobra.Command {
	opts := syncOptions{}
	cmd := &cobra.Command{
		Use:     "sync [OPTIONS]", // [FILTER...]
		Aliases: []string{"s"},
		Short:   "Synchronize",
		Args:    cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			// if len(args) > 0 {
			// 	opts.matchName = args[0]
			// }
			// opts.filter = args
			return runSync(dotCli, opts)
		},
	}
	return cmd
}

func runSync(dotCli *DotCli, opts syncOptions) error {
	roles := dotCli.Roles()
	length := len(roles)
	errs := make(chan error, length)
	for _, r := range roles {
		// fmt.Fprintf(dotCli.Out(), "Syncing %s...\n", r.Name)
		go func(name, path, url string) {
			repo, err := git.NewRepo(path, url)
			if err != nil {
				errs <- err
				return
			}
			exists, err := system.IsDir(path)
			if err != nil {
				errs <- err
				return
			}
			if exists {
				// fmt.Fprintf(dotCli.Out(), "Checking %s...\n", name)
				if err := repo.Status(); err != nil {
					errs <- err
					return
				}
				if err := repo.Pull(); err != nil {
					errs <- err
					return
				}
			} else {
				fmt.Fprintf(dotCli.Out(), "Cloning %s into %s...\n", name, repo.Dir)
				if err := repo.Clone(); err != nil {
					errs <- err
					return
				}
			}
			errs <- nil
		}(r.Name, r.Path, r.URL)
	}
	errors := cli.Errors{}
	for i := 0; i < length; i++ {
		if err := <-errs; err != nil {
			errors = append(errors, err)
		}
		// err := <-errs
		// switch err {
		// case git.ErrDirtyRepo:
		// 	if opts.force {
		// 		continue
		// 	}
		// 	fallthrough
		// default:
		// 	errors = append(errors, err)
		// }
	}
	if len(errors) > 0 {
		fmt.Fprintf(dotCli.Out(), "Synchronisation failed: %d/%d\n", len(errors), length)
		return fmt.Errorf(errors.Error())
	}
	return nil
}
