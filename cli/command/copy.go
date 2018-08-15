package command

import (
	"fmt"
	// "golang.org/x/net/context"

	// "github.com/docker/docker/api/types"
	// "github.com/docker/docker/cli/command/formatter"
	// "github.com/docker/docker/opts"
	// "github.com/LEI/dot/cli/config"
	"github.com/LEI/dot/cli/config/tasks"
	"github.com/spf13/cobra"
)

type copyOptions struct {
	action string
	// matchName string
	// quiet       bool
	// all         bool
	// noTrunc     bool
	// showDigests bool
	// format      string
	// filter      opts.FilterOpt
}

// NewCopyCommand creates a new `dot copy` command
func NewCopyCommand(dotCli *DotCli) *cobra.Command {
	opts := copyOptions{} // filter: opts.NewFilterOpt()
	cmd := &cobra.Command{
		Use:     "copy [OPTIONS]",
		Aliases: []string{"cp"},
		Short:   "Copy file",
		Args:    cobra.NoArgs, // RequiresMaxArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			// if len(args) > 0 {
			// 	opts.matchName = args[0]
			// }
			// fmt.Fprintf(dotCli.Out(), "CMD: %+v\n", cmd)
			opts.action = cmd.Parent().Name()
			return runCopy(dotCli, opts)
		},
	}
	Options.InstallTaskFlags(cmd.Flags())
	Options.InstallTaskPersistentFlags(cmd.PersistentFlags())
	// flags := cmd.Flags()
	// flags.BoolVarP(&opts.quiet, "quiet", "q", false, "Only show numeric IDs")
	// flags.BoolVarP(&opts.all, "all", "a", false, "Show all images (default hides intermediate images)")
	// flags.BoolVar(&opts.noTrunc, "no-trunc", false, "Don't truncate output")
	// flags.BoolVar(&opts.showDigests, "digests", false, "Show digests")
	// flags.StringVar(&opts.format, "format", "", "Pretty-print images using a Go template")
	// flags.VarP(&opts.filter, "filter", "f", "Filter output based on conditions provided")
	return cmd
}

// func newCopyCommand(dotCli *DotCli) *cobra.Command {
// 	cmd := *NewCopyCommand(dotCli)
// 	cmd.Aliases = []string{"ln"}
// 	cmd.Use = "copy [OPTIONS]"
// 	return &cmd
// }

func runCopy(dotCli *DotCli, opts copyOptions) error {
	roles := dotCli.Roles()
	if len(roles) == 0 {
		return fmt.Errorf("no roles to copy")
	}
	for _, r := range roles {
		if err := tasks.Check(r.Files); err != nil {
			return err
		}
	}
	for _, r := range roles {
		switch opts.action {
		case "install":
			if err := tasks.Install(r.Files); err != nil {
				return err
			}
		case "remove":
			if err := tasks.Remove(r.Files); err != nil {
				return err
			}
		default:
			return fmt.Errorf("%s: not implemented", opts.action)
		}
	}
	return nil
}
