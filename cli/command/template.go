package command

import (
	"fmt"

	"github.com/LEI/dot/cli/config/tasks"
	"github.com/spf13/cobra"
)

type templateOptions struct {
	action string
}

// NewTemplateCommand creates a new `dot template` command
func NewTemplateCommand(dotCli *DotCli) *cobra.Command {
	opts := templateOptions{} // filter: opts.NewFilterOpt()
	cmd := &cobra.Command{
		Use:     "template [OPTIONS]",
		Aliases: []string{},
		Short:   "Generate template",
		Args:    cobra.NoArgs, // RequiresMaxArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			opts.action = cmd.Parent().Name()
			return runTemplate(dotCli, opts)
		},
	}
	Options.InstallTaskFlags(cmd.Flags())
	Options.InstallTaskPersistentFlags(cmd.PersistentFlags())
	return cmd
}

func runTemplate(dotCli *DotCli, opts templateOptions) error {
	roles := dotCli.Roles()
	if len(roles) == 0 {
		return fmt.Errorf("no roles to template")
	}
	for _, r := range roles {
		if err := tasks.Check(r.Templates); err != nil {
			return err
		}
	}
	for _, r := range roles {
		switch opts.action {
		case "install":
			if err := tasks.Install(r.Templates); err != nil {
				return err
			}
		case "remove":
			if err := tasks.Remove(r.Templates); err != nil {
				return err
			}
		default:
			return fmt.Errorf("%s: not implemented", opts.action)
		}
	}
	return nil
}
