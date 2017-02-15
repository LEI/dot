package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func init() {
	RootCmd.AddCommand(syncCmd)
}

var syncCmd = &cobra.Command{
	// Hidden: true,
	Use: "sync",
	Short: "",
	Long: ``,
	Run: func(cmd *cobra.Command, args []string) {
		// fmt.Printf("cmd: %+v\n", cmd)
		// fmt.Printf("args: %+v\n", args)

		for name, pkg := range viper.GetStringMap("packages") {
			fmt.Printf("Package %s: %+v\n", name, pkg)
			pkg.Repo.CloneOrPull()
		}

		// fmt.Println(viper.Get("dir"), viper.Get("dirs"))
		// fmt.Println(viper.Get("link"), viper.Get("links"))
		// fmt.Println(viper.Get("line"), viper.Get("lines"))

		// Packages = viper.Get("packages").(PackageSlice)
	},
}
