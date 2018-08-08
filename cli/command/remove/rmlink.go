package remove

import (
	"fmt"
	// "golang.org/x/net/context"

	// "github.com/docker/docker/api/types"
	"github.com/LEI/dot/cli"
	"github.com/LEI/dot/cli/command"
	// "github.com/LEI/dot/cli/config"
	// "github.com/docker/docker/cli/command/formatter"
	// "github.com/docker/docker/opts"
	"github.com/spf13/cobra"
)

type linkOptions struct {
}

// NewRmLinkCommand creates a new `dot link` command
func NewRmLinkCommand(dotCli *command.DotCli) *cobra.Command {
	opts := linkOptions{} // filter: opts.NewFilterOpt()
	cmd := &cobra.Command{
		Use:   "link [ACTION] [OPTIONS]",
		Aliases: []string{"ln"},
		Short: "Unlink",
		Args: cli.NoArgs, // RequiresMaxArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			// if len(args) > 0 {
			// 	opts.matchName = args[0]
			// }
			// fmt.Printf("CMD: %+v\n", cmd)
			return runRmLink(dotCli, opts)
		},
	}
	// flags := cmd.Flags()
	// flags.BoolVarP(&opts.quiet, "quiet", "q", false, "Only show numeric IDs")
	return cmd
}

// func newRmLinkCommand(dotCli *command.DotCli) *cobra.Command {
// 	cmd := *NewRmLinkCommand(dotCli)
// 	cmd.Aliases = []string{"ln"}
// 	cmd.Use = "link [OPTIONS]"
// 	return &cmd
// }

func runRmLink(dotCli *command.DotCli, opts linkOptions) error {
	fmt.Fprintf(dotCli.Out(), "REMOVE LINK %+v\n", opts)
	// Check links
	for _, r := range dotCli.Roles() {
		if err := r.Links.Check(); err != nil {
			return err
		}
	}
	// Remove links
	for _, r := range dotCli.Roles() {
		if err := r.Links.Execute(); err != nil {
			return err
		}
	}
	return nil
}
