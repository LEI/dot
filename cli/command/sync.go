package command
// https://github.com/moby/moby/tree/17.05.x/cli/command/image

import (
	"fmt"
	"github.com/spf13/cobra"

	// "github.com/LEI/dot/cli/config"
	"github.com/LEI/dot/pkg/git"
	"github.com/LEI/dot/system"
)

type syncOptions struct {
	// check bool
	// filter []string
	force bool
}

// NewSyncCommand returns a cobra command for `sync` subcommands
func NewSyncCommand(dotCli *DotCli) *cobra.Command {
	opts := syncOptions{}
	cmd := &cobra.Command{
		Use:   "sync", // [OPTIONS] [FILTER...]
		Aliases: []string{"s"},
		Short: "Synchronize",
		// Args: cobra.NoArgs,
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
	for _, r := range dotCli.Roles() {
		// fmt.Fprintf(dotCli.Out(), "Syncing %s...\n", r.Name)
		repo, err := git.NewRepo(r.Path, r.URL)
		if err != nil {
			return err
		}
		exists, err := system.IsDir(r.Path)
		if err != nil {
			return err
		}
		if exists {
			fmt.Fprintf(dotCli.Out(), "Checking %s...\n", r.Name)
			if err := repo.Status(); err != nil {
				return err
			}
			if err := repo.Pull(); err != nil {
				return err
			}
		} else {
			fmt.Fprintf(dotCli.Out(), "Cloning %s...\n", r.Name)
			if err := repo.Clone(); err != nil {
				return err
			}
		}
	}
	return nil
}
