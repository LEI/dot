package command

import (
	"fmt"
	// "golang.org/x/net/context"

	// "github.com/docker/docker/api/types"
	// "github.com/LEI/dot/cli/config"
	// "github.com/docker/docker/cli/command/formatter"
	// "github.com/docker/docker/opts"
	"github.com/spf13/cobra"
)

type rmLinkOptions struct {
}

// NewUnlinkCommand creates a new `dot link` command
func NewUnlinkCommand(dotCli *DotCli) *cobra.Command {
	opts := rmLinkOptions{} // filter: opts.NewFilterOpt()
	cmd := &cobra.Command{
		Use:   "unlink [ACTION] [OPTIONS]",
		// Aliases: []string{"rmln"},
		Short: "Remove symbolic link",
		Args: cobra.NoArgs, // RequiresMaxArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			// if len(args) > 0 {
			// 	opts.matchName = args[0]
			// }
			// fmt.Printf("CMD: %+v\n", cmd)
			return runUnlink(dotCli, opts)
		},
	}
	// flags := cmd.Flags()
	// flags.BoolVarP(&opts.quiet, "quiet", "q", false, "Only show numeric IDs")
	return cmd
}

// func newUnlinkCommand(dotCli *DotCli) *cobra.Command {
// 	cmd := *NewUnlinkCommand(dotCli)
// 	cmd.Aliases = []string{"ln"}
// 	cmd.Use = "link [OPTIONS]"
// 	return &cmd
// }

func runUnlink(dotCli *DotCli, opts rmLinkOptions) error {
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
