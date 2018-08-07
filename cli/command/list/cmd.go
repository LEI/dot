package list
// https://github.com/moby/moby/tree/17.05.x/cli/command/image

import (
	"fmt"
	"github.com/spf13/cobra"

	"github.com/LEI/dot/cli"
	"github.com/LEI/dot/cli/command"
	"github.com/LEI/dot/cli/config"
)

type listOptions struct {}

// NewListCommand returns a cobra command for `list` subcommands
func NewListCommand(dotCli *command.DotCli) *cobra.Command {
	opts := listOptions{}
	cmd := &cobra.Command{
		Use:   "list",
		Short: "List",
		Args:  cli.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			// if len(args) > 0 {
			// 	opts.matchName = args[0]
			// }
			return runList(dotCli, opts)
		},
	}
	return cmd
}

func runList(dotCli *command.DotCli, opts listOptions) error {
	fmt.Fprintf(dotCli.Out(), "RUN LINK %+v\nCONFIG:%+v\n", opts, config.DotConfig)
	// client := dotCli.Client()
	// options := types.NetworkLinkOptions{Filters: opts.filter.Value()}
	// networkResources, err := client.NetworkLink(context.Background(), options)
	// if err != nil {
	// 	return err
	// }

	// format := opts.format
	// if len(format) == 0 {
	// 	if len(dotCli.Config().NetworksFormat) > 0 && !opts.quiet {
	// 		format = dotCli.Config().NetworksFormat
	// 	} else {
	// 		format = formatter.TableFormatKey
	// 	}
	// }

	// sort.Sort(byNetworkName(networkResources))

	// networksCtx := formatter.Context{
	// 	Output: dotCli.Out(),
	// 	Format: formatter.NewNetworkFormat(format, opts.quiet),
	// 	Trunc:  !opts.noTrunc,
	// }
	// return formatter.NetworkWrite(networksCtx, networkResources)
	return nil
}
