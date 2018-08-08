package remove

import (
	"fmt"

	"github.com/LEI/dot/cli"
	"github.com/LEI/dot/cli/command"
	"github.com/spf13/cobra"
)

type dirOptions struct {
	quiet       bool
}

// NewRmDirCommand creates a new `dot dir` command
func NewRmDirCommand(dotCli *command.DotCli) *cobra.Command {
	opts := dirOptions{} // filter: opts.NewFilterOpt()
	cmd := &cobra.Command{
		Use:   "dir [ACTION] [OPTIONS]",
		Aliases: []string{"ln"},
		Short: "Dir",
		Args: cli.NoArgs, // RequiresMaxArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return runRmDir(dotCli, opts)
		},
	}
	flags := cmd.Flags()
	flags.BoolVarP(&opts.quiet, "quiet", "q", false, "Only show numeric IDs")
	return cmd
}

// func newRmDirCommand(dotCli *command.DotCli) *cobra.Command {
// 	cmd := *NewRmDirCommand(dotCli)
// 	cmd.Aliases = []string{"ln"}
// 	cmd.Use = "dir [OPTIONS]"
// 	return &cmd
// }

func runRmDir(dotCli *command.DotCli, opts dirOptions) error {
	fmt.Fprintf(dotCli.Out(), "RUN DIR %+v\n", opts)
	// Check dirs
	for _, r := range dotCli.Roles() {
		if err := r.Dirs.Check(); err != nil {
			return err
		}
	}
	// Remove dirs
	for _, r := range dotCli.Roles() {
		if err := r.Dirs.Execute(); err != nil {
			return err
		}
	}
	return nil
}
