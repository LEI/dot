package cmd

import (
	"fmt"
	"github.com/LEI/dot/role"
	"github.com/spf13/cobra"
	"os"
)

func init() {
	RootCmd.AddCommand(syncCmd)
}

var syncCmd = &cobra.Command{
	// Hidden: true,
	Use:   "sync [clone...]",
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
		ok := checkPackage(pkg)
		if !ok {
			continue
		}
		fmt.Printf("[%s]\n", pkg.Name)
		err := syncPackage(pkg)
		if err != nil {
			return err
		}
		fmt.Printf("[%s] Done\n", pkg.Name)
	}
	return nil
}

func checkPackage(pkg *role.Package) bool {
	for _, osType := range pkg.Os {
		switch osType {
		case OS:
			return true
		default:
			if Debug {
				fmt.Fprintf(os.Stderr, "[%s] %s: unsupported platform, only for %+v\n", pkg.Name, OS, pkg.Os)
			}
			return false
		}
	}
	return true
}

func syncPackage(pkg *role.Package) error {
	err := initPackage(pkg)
	if err != nil {
		return err
	}
	if pkg.Config != nil {
		fmt.Println("-> DIRS", pkg.Config.GetStringSlice("dirs"))
		fmt.Println("-> LINKS", pkg.Config.Get("links"))
		fmt.Println("-> LINES", pkg.Config.GetStringMapString("lines"))
	} else {
		fmt.Println("-> NIL")
	}
	return nil
}
