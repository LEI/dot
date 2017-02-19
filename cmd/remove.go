package cmd

import (
	"fmt"
	// "github.com/LEI/dot/config"
	// "github.com/LEI/dot/fileutil"
	// "github.com/LEI/dot/role"
	"github.com/spf13/cobra"
	// "os"
	// "path"
	// "strings"
)

var ()

var removeCmd = &cobra.Command{
	Use:   "remove [flags]",
	Short: "Remove dotfiles",
	// Long:  ``,
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Println("Not implemented")
		if len(args) > 0 {
			fmt.Printf("Args: %s\n", args)
			// return cmd.Help()
		}
		return nil
	},
}

func init() {
	DotCmd.AddCommand(removeCmd)

	// removeCmd.Flags().BoolVarP(&, "", "", , "")
	// Config.BindPFlags(removeCmd.Flags())
}
