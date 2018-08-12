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

type linkOptions struct {
	action string
	// matchName string
	// quiet       bool
	// all         bool
	// noTrunc     bool
	// showDigests bool
	// format      string
	// filter      opts.FilterOpt
}

// NewLinkCommand creates a new `dot link` command
func NewLinkCommand(dotCli *DotCli) *cobra.Command {
	opts := linkOptions{} // filter: opts.NewFilterOpt()
	cmd := &cobra.Command{
		Use:     "link [OPTIONS]",
		Aliases: []string{"ln"},
		Short:   "Link",
		Args:    cobra.NoArgs, // RequiresMaxArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			// if len(args) > 0 {
			// 	opts.matchName = args[0]
			// }
			// fmt.Fprintf(dotCli.Out(), "CMD: %+v\n", cmd)
			opts.action = cmd.Parent().Name()
			return runLink(dotCli, opts)
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

// func newLinkCommand(dotCli *DotCli) *cobra.Command {
// 	cmd := *NewLinkCommand(dotCli)
// 	cmd.Aliases = []string{"ln"}
// 	cmd.Use = "link [OPTIONS]"
// 	return &cmd
// }

func runLink(dotCli *DotCli, opts linkOptions) error {
	roles := dotCli.Roles()
	if len(roles) == 0 {
		return fmt.Errorf("no roles to link")
	}
	for _, r := range roles {
		if err := tasks.Check(r.Links); err != nil {
			return err
		}
	}
	for _, r := range roles {
		switch opts.action {
		case "install":
			if err := tasks.Install(r.Links); err != nil {
				return err
			}
		case "remove":
			if err := tasks.Remove(r.Links); err != nil {
				return err
			}
		default:
			return fmt.Errorf("%s: not implemented", opts.action)
		}
	}

	// cfg := dotCli.Config().GetAll()
	// for i, r := range cfg["roles"].([]interface{}) {
	// 	// ri := r.(map[string]interface{})
	// 	role := config.NewRole(r)
	// 	fmt.Fprintf(dotCli.Out(), "%d: %+v\n", i, role)
	// }

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
