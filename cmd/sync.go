package cmd

import (
	"fmt"
	"github.com/LEI/dot/role"
	"github.com/spf13/cobra"
	// "os"
)

func init() {
	RootCmd.AddCommand(syncCmd)
}

var syncCmd = &cobra.Command{
	// Hidden: true,
	Use:   "sync [flags]",
	Aliases: []string{"s"},
	Short: "",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		// if Debug {
		// 	fmt.Printf("Synchronize packages: %+v\n", Packages)
		// }
		err := syncPackages(Source, Target, Packages)
		if err != nil {
			fatal(err)
		}
	},
}

func syncPackages(source string, target string, packages []*role.Package) error {
	for _, pkg := range packages {
		ok := pkg.Check(OS)
		if !ok {
			continue
		}
		fmt.Printf("[%s]\n", pkg.Name)
		// err := initPackage(pkg)
		err := pkg.InitRepo()
		if err != nil {
			return err
		}
		err = pkg.Repo.CloneOrPull()
		if err != nil {
			return err
		}
		err = pkg.InitConfig(ConfigName)
		if err != nil {
			return err
		}
		err = pkg.Sync(source, target)
		if err != nil {
			return err
		}
		fmt.Printf("[%s] Done\n", pkg.Name)
	}
	return nil
}
