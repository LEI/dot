package command

import (
	"fmt"

	"github.com/spf13/cobra"
)

type rmDirOptions struct {
	quiet       bool
}

// NewRmDirCommand creates a new `dot dir` command
func NewRmDirCommand(dotCli *DotCli) *cobra.Command {
	opts := rmDirOptions{} // filter: opts.NewFilterOpt()
	cmd := &cobra.Command{
		Use:   "rmdir [ACTION] [OPTIONS]",
		// Aliases: []string{"rmd"},
		Short: "Remove directory",
		Args: cobra.NoArgs, // RequiresMaxArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return runRmDir(dotCli, opts)
		},
	}
	flags := cmd.Flags()
	flags.BoolVarP(&opts.quiet, "quiet", "q", false, "Only show numeric IDs")
	return cmd
}

// func newRmDirCommand(dotCli *DotCli) *cobra.Command {
// 	cmd := *NewRmDirCommand(dotCli)
// 	cmd.Aliases = []string{"ln"}
// 	cmd.Use = "dir [OPTIONS]"
// 	return &cmd
// }

func runRmDir(dotCli *DotCli, opts rmDirOptions) error {
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
