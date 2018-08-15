package command

// https://github.com/ncw/rclone

import (
	"fmt"

	"github.com/LEI/dot/cli"
	cliconfig "github.com/LEI/dot/cli/config"
	"github.com/spf13/cobra"
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
	// dotCli.AddCommands(cmd)
	// Options.InstallActionFlags(cmd.Flags())
	Options.InstallActionPersistentFlags(cmd.PersistentFlags())

	// flags := cmd.Flags() // var flags *pflag.FlagSet
	// flags.StringVarP(&opts.format, "format", "f", "", "List format")

	return cmd
}

func runSync(dotCli *DotCli, opts syncOptions) error {
	roles := dotCli.Roles()
	length := len(roles)
	errs := make(chan error, length)
	for _, r := range roles {
		// fmt.Fprintf(dotCli.Out(), "Syncing %s...\n", r.Name)
		go func(r *cliconfig.Role) {
			if err := r.Sync(); err != nil {
				errs <- err
				return
			}
			errs <- nil
		}(r)
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
