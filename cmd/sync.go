package cmd

import (
	"fmt"
	// "github.com/LEI/dot/git"
	"github.com/LEI/dot/role"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// var Packages role.PackageSlice

func init() {
	RootCmd.AddCommand(syncCmd)
}

var syncCmd = &cobra.Command{
	// Hidden: true,
	Use: "sync",
	Short: "",
	Long: ``,
	Run: func(cmd *cobra.Command, args []string) {
		packages := viper.GetStringMap("packages")
		// fmt.Printf("cmd: %+v\n", cmd)
		fmt.Printf("Synchronize packages: %+v\n", packages)
		for name, p := range packages {
			pkg, ok := p.(role.Package)
			if !ok {
				er(fmt.Errorf("%s: could not assert package type", name))
			}
			fmt.Printf("Package %s: %+v\n", name, pkg)
			fmt.Printf("Repo %+v\n", pkg.Repo)
			err := pkg.Repo.CloneOrPull()
			if err != nil {
				er(err)
			}
		}

		// fmt.Println(viper.Get("dir"), viper.Get("dirs"))
		// fmt.Println(viper.Get("link"), viper.Get("links"))
		// fmt.Println(viper.Get("line"), viper.Get("lines"))

		// Packages = viper.Get("packages").(PackageSlice)
	},
}
