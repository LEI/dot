package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func init() {
	RootCmd.AddCommand(RemoveCmd)
}

var RemoveCmd = &cobra.Command{
	Hidden: true,
	Use: "remove",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("cmd: %+v\n", cmd)
		fmt.Printf("args: %+v\n", args)

		// fmt.Printf("cmd: %+v\n", cmd)
		// fmt.Printf("args: %+v\n", args)
		fmt.Printf("%s: %+v\n", "Packages", viper.Get("packages"))
		fmt.Println(viper.Get("dir"), viper.Get("dirs"))
		fmt.Println(viper.Get("link"), viper.Get("links"))
		fmt.Println(viper.Get("line"), viper.Get("lines"))

		// Packages = viper.Get("packages").(PackageSlice)
	},
}
