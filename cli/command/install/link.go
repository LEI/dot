package install

import (
	"fmt"
	// "golang.org/x/net/context"

	// "github.com/docker/docker/api/types"
	"github.com/LEI/dot/cli"
	"github.com/LEI/dot/cli/command"
	// "github.com/docker/docker/cli/command/formatter"
	// "github.com/docker/docker/opts"
	"github.com/spf13/cobra"
)

type linkOptions struct {
	// matchName string

	quiet       bool
	// all         bool
	// noTrunc     bool
	// showDigests bool
	// format      string
	// filter      opts.FilterOpt
}

// NewLinkCommand creates a new `dot link` command
func NewLinkCommand(dotCli *command.DotCli) *cobra.Command {
	opts := linkOptions{} // filter: opts.NewFilterOpt()
	cmd := &cobra.Command{
		Use:   "link [ACTION] [OPTIONS]",
		Aliases: []string{"ln"},
		Short: "Link",
		Args: cli.NoArgs, // RequiresMaxArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			// if len(args) > 0 {
			// 	opts.matchName = args[0]
			// }
			return runLink(dotCli, opts)
		},
	}
	flags := cmd.Flags()
	flags.BoolVarP(&opts.quiet, "quiet", "q", false, "Only show numeric IDs")
	// flags.BoolVarP(&opts.all, "all", "a", false, "Show all images (default hides intermediate images)")
	// flags.BoolVar(&opts.noTrunc, "no-trunc", false, "Don't truncate output")
	// flags.BoolVar(&opts.showDigests, "digests", false, "Show digests")
	// flags.StringVar(&opts.format, "format", "", "Pretty-print images using a Go template")
	// flags.VarP(&opts.filter, "filter", "f", "Filter output based on conditions provided")
	return cmd
}

// func newLinkCommand(dotCli *command.DotCli) *cobra.Command {
// 	cmd := *NewLinkCommand(dotCli)
// 	cmd.Aliases = []string{"ln"}
// 	cmd.Use = "link [OPTIONS]"
// 	return &cmd
// }

func runLink(dotCli *command.DotCli, opts linkOptions) error {
	fmt.Fprintf(dotCli.Out(), "RUN LINK %+v\n", opts)
	// client := dotCli.Client()
	// options := types.NetworkLinkOptions{Filters: opts.filter.Value()}
	// networkResources, err := client.NetworkLink(context.Background(), options)
	// if err != nil {
	// 	return err
	// }

	// format := opts.format
	// if len(format) == 0 {
	// 	if len(dotCli.ConfigFile().NetworksFormat) > 0 && !opts.quiet {
	// 		format = dotCli.ConfigFile().NetworksFormat
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
