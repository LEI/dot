package cmd

import (
	"fmt"
	"github.com/LEI/dot/fileutil"
	"github.com/LEI/dot/role"
	"github.com/spf13/cobra"
	"os"
	"path/filepath"
	"strings"
)

var Ignore = []string{".git", ".*\\.md"}

func init() {
	RootCmd.AddCommand(syncCmd)
}

var syncCmd = &cobra.Command{
	// Hidden: true,
	Use:     "sync [flags]",
	Aliases: []string{"s"},
	Short:   "",
	Long:    ``,
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

func syncPackages(source, target string, packages []*role.Package) error {
	for _, pkg := range packages {
		ok := pkg.CheckOsType(OsTypes)
		if !ok {
			continue
		}
		fmt.Printf("[%s]\n", pkg.Name)
		err := pkg.InitRepo(Https)
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
		if Debug {
			fmt.Printf("Using: %s for %s package\n", pkg.Config.ConfigFileUsed(), pkg.Name)
		}
		if pkg.Config == nil {
			return fmt.Errorf("%s: no config", pkg.Name)
		}
		for _, dir := range pkg.GetDirs() {
			dir = os.ExpandEnv(dir)
			dir = filepath.Join(target, dir)
			err := fileutil.MakeDir(dir)
			if err != nil {
				return err
			}
			fmt.Printf("Directory: %s\n", dir)
		}
		for _, value := range pkg.GetLinks() {
			var file *role.File
			switch val := value.(type) {
			case string:
				file = &role.File{Type: "", Path: val}
			case map[string]interface{}:
				file = &role.File{Type: val["type"].(string), Path: val["path"].(string)}
			case *role.File:
				file = val
			// case interface{}:
			default:
				fmt.Printf("Unknown type %T for %+v\n", val, val)
				file = val.(*role.File)
				// return file, fmt.Errorf("Unknown type %T for %+v\n", val, val)
			}
			if file.Path == "" {
				return fmt.Errorf("Empty link path\n")
			}
			// if filepath.IsAbs(file.Path) {
			// 	fmt.Printf("%s: file path is not absolute\n", file.Path)
			// }
			file.Path = os.ExpandEnv(file.Path)
			// if err != nil {
			// 	return err
			// }
			link := &role.Link{File: file}
			// link, err := role.NewLink(glob)
			// if err != nil {
			// 	return err
			// }
			link.Path = filepath.Join(pkg.Path, link.Path)
			links, err := link.GlobAsLink()
			if err != nil {
				return err
			}
			if len(links) == 0 {
				fmt.Fprintf(os.Stderr, "%s: No match\n", link.Path)
			}
			LOOP:
			for _, link := range links {
				matched, err := link.NameMatches(Ignore)
				if err != nil {
					return err
				}
				if matched {
					fmt.Printf("Ignoring link: %s\n", link)
					continue LOOP
				}
				dst := strings.Replace(link.Path, pkg.Path, target, 1)
				err = link.Sync(dst)
				if err != nil {
					return err
				}
				fmt.Printf("Link: %s into %s\n", link.Path, dst)
			}
		}
		for file, line := range pkg.GetLines() {
			file = os.ExpandEnv(file)
			file = filepath.Join(target, file)
			err := fileutil.LineInFile(file, line)
			if err != nil {
				return err
			}
			fmt.Printf("Line: '%s' in %s\n", line, file)
		}
		// fmt.Printf("[%s] Done\n", pkg.Name)
	}
	return nil
}
