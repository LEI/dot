package main

import (
	"fmt"

	"github.com/spf13/cobra"
)

// TemplateOptions collects all options for the template command.
type TemplateOptions struct {
}

var templateOptions TemplateOptions

var cmdTemplate = &cobra.Command{
	Use:     "template [flags]",
	Aliases: []string{"tpl"},
	Short:   "Template managed files",
	Long: `
The "template" task installs or removes templates.
`,
	DisableAutoGenTag: true,
	Args:              cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		return runTemplate(templateOptions, globalOptions, args)
	},
}

func init() {
	cmdInstall.AddCommand(cmdTemplate)
	cmdRemove.AddCommand(cmdTemplate)

	// flags := cmdTemplate.Flags()
}

func runTemplate(opts TemplateOptions, gopts GlobalOptions, args []string) error {
	fmt.Println("runTemplate", args)
	return nil
}
