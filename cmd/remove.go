package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
)

func init() {
	RootCmd.AddCommand(RemoveCmd)
}

var RemoveCmd = &cobra.Command{
	// Hidden: true,
	Use: "remove",
	Short: "",
	Long: ``,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("cmd: %+v\n", cmd)
		fmt.Printf("args: %+v\n", args)
	},
}
