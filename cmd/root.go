package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"os"
)

var RootCmd = &cobra.Command{
	Use: "dot",
	Short: "Manage dotfiles",
	Long: ``,
	// PersistentPreRun: func(cmd *cobra.Command, args[]string) {
	// },
	Run: func(cmd *cobra.Command, args []string) {
		switch {
		// case viper.Get("remove"):
		// 	err := RemoveCmd.Execute()
		// 	if err != nil {
		// 		fmt.Printf("%v", err)
		// 		os.Exit(1)
		// 	}
		case viper.Get("sync"):
			err := SyncCmd.Execute()
			if err != nil {
				fmt.Printf("%v", err)
				os.Exit(1)
			}
		default:
			cmd.HelpFunc()(cmd, args)
			os.Exit(127) // Command not found
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
// 		panic(fmt.Errorf("%s: sub-command not found", name))
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
