package main

import (
	"github.com/spf13/cobra"
)

var cmdSync = &cobra.Command{
	Use:     "sync",
	Aliases: []string{"s"},
	Short:   "Synchronize a role",
	Long: `
The "sync" command clone or pull a role repository.
`,
	DisableAutoGenTag: true,
	Args:              cobra.NoArgs,
	RunE:              runSync,
}

func init() {
	cmdRoot.AddCommand(cmdSync)
}

func runSync(cmd *cobra.Command, args []string) error {
	// for _, r := range globalConfig.Roles {
	// 	if err := r.Sync(); err != nil {
	// 		return err
	// 	}
	// }
	return nil
}
