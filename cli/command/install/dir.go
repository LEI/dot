package install

import (
	"fmt"

	"github.com/LEI/dot/cli"
	"github.com/LEI/dot/cli/command"
	"github.com/spf13/cobra"
)

type dirOptions struct {
	quiet       bool
}

// NewDirCommand creates a new `dot dir` command
func NewDirCommand(dotCli *command.DotCli) *cobra.Command {
	opts := dirOptions{} // filter: opts.NewFilterOpt()
	cmd := &cobra.Command{
		Use:   "dir [ACTION] [OPTIONS]",
		Aliases: []string{"ln"},
		Short: "Dir",
		Args: cli.NoArgs, // RequiresMaxArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return runDir(dotCli, opts)
		},
	}
	flags := cmd.Flags()
	flags.BoolVarP(&opts.quiet, "quiet", "q", false, "Only show numeric IDs")
	return cmd
}

// func newDirCommand(dotCli *command.DotCli) *cobra.Command {
// 	cmd := *NewDirCommand(dotCli)
// 	cmd.Aliases = []string{"ln"}
// 	cmd.Use = "dir [OPTIONS]"
// 	return &cmd
// }

func runDir(dotCli *command.DotCli, opts dirOptions) error {
	fmt.Fprintf(dotCli.Out(), "RUN DIR %+v\n", opts)
	// Check dirs
	for _, r := range dotCli.Roles() {
		if err := r.Dirs.Check(); err != nil {
			return err
		}
	}
	// Install dirs
	for _, r := range dotCli.Roles() {
		if err := r.Dirs.Execute(); err != nil {
			return err
		}
	}
	return nil
}
