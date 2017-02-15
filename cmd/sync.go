package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
)

func init() {
	RootCmd.AddCommand(SyncCmd)
}

var SyncCmd = &cobra.Command{
	Hidden: true,
	Use: "sync",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("cmd: %+v\n", cmd)
		fmt.Printf("args: %+v\n", args)
	},
}
