package cmd

import (
	"fmt"
	// "github.com/LEI/dot/git"
	"github.com/LEI/dot/role"
	"github.com/spf13/cobra"
)

// var Packages role.PackageSlice

func init() {
	RootCmd.AddCommand(SyncCmd)
}

var SyncCmd = &cobra.Command{
	// Hidden: true,
	Use: "sync [clone...]",
	Short: "",
	Long: ``,
	Run: func(cmd *cobra.Command, args []string) {
		// packages := viper.GetStringMap("packages")
		// fmt.Printf("cmd: %+v\n", cmd)
		fmt.Printf("Synchronize packages: %+v\n", Packages)
		err := sync(Packages)
		if err != nil {
			er(err)
		}
	},
}

func sync(packages []role.Package) error {
	fmt.Println(Source, "->", Target)
	for name, pkg := range packages {
		// pkg, ok := p.(role.Package)
		// if !ok {
		// 	er(fmt.Errorf("%s: could not assert package type", name))
		// }
		fmt.Printf("Package %s: %+v\n", name, pkg)
		// fmt.Printf("Repo %+v\n", pkg.Repo)
		// err := pkg.Repo.CloneOrPull()
		// if err != nil {
		// 	er(err)
		// }
	}
	return nil
}
