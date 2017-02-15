package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	// "os"
)

var RootCmd = &cobra.Command{
	Use: "dot",
	Short: "Manage dotfiles",
	Long: ``,
	// PersistentPreRun: func(cmd *cobra.Command, args[]string) {
	// },
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Error: no command specified", args)
		err := cmd.RootCmd.Help()
		if err != nil {
			er(err)
		}
	},
}

// func GetCommand(cmd *cobra.Command, name string) *cobra.Command {
// 	subCommands := cmd.Commands()
// 	index := GetCommandIndex(subCommands, func(i int) bool {
// 		fmt.Println(subCommands[i])
// 		return subCommands[i].name == name
// 	})
// 	if index < 0 {
// 		er(fmt.Errorf("%s: sub-command not found", name))
// 	}
// 	return subCommands[index]
// }

// func GetCommandIndex(slice []*cobra.Command, predicate func(i int) bool) int {
// 	for index, _ := range slice {
// 		if predicate(index) {
// 			return index
// 		}
// 	}
// 	return -1
// }
