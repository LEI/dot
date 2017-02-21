package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
)

var syncCmd = &cobra.Command{
	Hidden: true,
	Use:    "Sync [flags]",
	Short:  "Sync dotfiles",
	// Long:   ``,
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Println("Not implemented")
		if len(args) > 0 {
			fmt.Println("Extra arguments:", args)
			return cmd.Help()
		}
		return syncCommand(Dot, filter)
	},
}

func init() {
	DotCmd.AddCommand(syncCmd)
	// syncCmd.Flags().BoolVarP(&, "", "", , "")
	// Config.BindPFlags(syncCmd.Flags())
}
