package cli
// https://github.com/moby/moby/blob/master/cli/cobra.go

import (
	"github.com/spf13/cobra"
)

// SetupRootCommand sets default usage
func SetupRootCommand(rootCmd *cobra.Command) {
	// cobra.AddTemplateFunc("hasSubCommands", hasSubCommands)
	// cobra.AddTemplateFunc("hasManagementSubCommands", hasManagementSubCommands)
	// cobra.AddTemplateFunc("operationSubCommands", operationSubCommands)
	// cobra.AddTemplateFunc("managementSubCommands", managementSubCommands)
	// cobra.AddTemplateFunc("wrappedFlagUsages", wrappedFlagUsages)

	// rootCmd.SetUsageTemplate(usageTemplate)
	// rootCmd.SetHelpTemplate(helpTemplate)
	// rootCmd.SetFlagErrorFunc(FlagErrorFunc)
	// rootCmd.SetHelpCommand(helpCommand)

	rootCmd.PersistentFlags().BoolP("help", "h", false, "Print usage")
}
