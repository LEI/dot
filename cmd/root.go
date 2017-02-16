package cmd

import (
	"fmt"
	// "github.com/LEI/dot/cmd"
	"github.com/LEI/dot/role"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"log"
	"os"
	"runtime"
)

const (
	OS      = runtime.GOOS
	version = "master"
)

var (
	Config = viper.New()
	// PkgConfig = make(map[string]*viper.Viper, 0)
	// flag = pflag.NewFlagSet(os.Args[0], pflag.ExitOnError)
	// Skip = fmt.Errorf("Skip this path")
	HomeDir, CurrentDir string
	Source, Target      string
	Debug, ForceYes     bool
	cfgFile             = ""
	cfgName             = ".dotrc"
	IgnoreFiles         = []string{".git", ".*\\.md"}
	Packages            []*role.Package
	// Packages            role.PackageSlice //= make(role.PackageSlice, 0)
)

var (
	InfoSymbol    = "›"
	SuccessSymbol = "✓" // ✓ ✔
	ErrorSymbol   = "✘" // × ✕ ✖ ✗ ✘
	WarnSymbol    = "!" // ⚠ !
	logInfo       = log.New(os.Stdout, InfoSymbol+" ", 0)
	logSuccess    = log.New(os.Stdout, SuccessSymbol+" ", 0)
	logWarn       = log.New(os.Stderr, WarnSymbol+" ", log.Lshortfile)
	logError      = log.New(os.Stderr, ErrorSymbol+" ", log.Llongfile)
)

// type Configuration struct {
// 	Packages *role.PackageSlice
// }

func init() {
	cobra.OnInitialize(initConfig)

	HomeDir = os.Getenv("HOME")   // user.Current().HomeDir
	CurrentDir, err := os.Getwd() // os.Getenv("PWD")
	if err != nil {
		fatal(err)
	}

	// flags := RootCmd.Flags()
	// pflags := RootCmd.PersistentFlags()

	// RootCmd.Flags().BoolVarP(&Sync, "sync", "S", Sync, "Synchronize packages")
	// RootCmd.Flags().BoolVarP(&Remove, "remove", "R", Remove, "Remove packages")

	RootCmd.PersistentFlags().StringVarP(&Source, "source", "s", CurrentDir, "Source `directory`")
	RootCmd.PersistentFlags().StringVarP(&Target, "target", "t", HomeDir, "Destination `directory`")

	RootCmd.PersistentFlags().BoolVarP(&Debug, "debug", "d", Debug, "Check mode")
	RootCmd.PersistentFlags().BoolVarP(&ForceYes, "force", "f", ForceYes, "Force yes")
	RootCmd.PersistentFlags().StringVarP(&cfgFile, "config", "c", cfgFile, "Configuration `file`")
	// RootCmd.PersistentFlags().VarP(&Packages, "package", "p", "List of packages `[name=]user/repo`")

	// 	viper.SetDefault("Source", CurrentDir)
	// 	viper.SetDefault("Target", HomeDir)

	// viper.RegisterAlias("src", "source")
	// viper.RegisterAlias("dst", "target")

	Config.BindPFlags(RootCmd.PersistentFlags())
	// RootCmd.PersistentFlags().Parse(os.Args[1:])
}

var RootCmd = &cobra.Command{
	Use:   "dot",
	Short: "Manage dotfiles",
	Long:  ``,
	PersistentPreRun: func(cmd *cobra.Command, args[]string) {
		Source = Config.GetString("source")
		Target = Config.GetString("target")
		err := Config.UnmarshalKey("packages", &Packages)
		if err != nil {
			fatal(err)
		}
	},
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			os.Args = append(os.Args, "sync")
			if err := syncCmd.Execute(); err != nil {
				fatal(err)
			}
		} else {
			fmt.Println("Extra arguments: %s", args)
			err := cmd.Help()
			if err != nil {
				fatal(err)
			}
		}
	},
}

func Execute() {
	if err := RootCmd.Execute(); err != nil {
		fatal(err)
	}
}

// func GetCommand(cmd *cobra.Command, name string) *cobra.Command {
// 	subCommands := cmd.Commands()
// 	index := GetCommandIndex(subCommands, func(i int) bool {
// 		fmt.Println(subCommands[i])
// 		return subCommands[i].name == name
// 	})
// 	if index < 0 {
//		fatal(fmt.Errorf("%s: sub-command not found", name))
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

func initConfig() {
	if cfgFile != "" {
		Config.SetConfigFile(cfgFile)
	}

	configPaths := []string{HomeDir, Source}
	// if CurrentDir != Source {
	// 	configPaths = append(configPaths, ".")
	// }
	// viper.GetViper()
	err := readConfig(Config, cfgName, configPaths...)
	if err != nil && os.IsExist(err) {
		fatal(err)
	}
}

func readConfig(v *viper.Viper, name string, paths ...string) error {
	v.SetConfigName(name)
	for _, path := range paths {
		v.AddConfigPath(path)
	}
	// Read in environment variables that match
	v.AutomaticEnv()
	// if err := v.ReadInConfig(); err != nil {
	// 	return err
	// }
	err := v.ReadInConfig()
	cfgFileUsed := v.ConfigFileUsed()
	if cfgFileUsed != "" {
		fmt.Printf("Found config file: %s\n", cfgFileUsed)
	} else {
		fmt.Println("No config file found")
	}
	return err
}

func InitPackage(pkg *role.Package) error {
	err := pkg.InitRepo()
	if err != nil {
		return err
	}
	err = pkg.Repo.CloneOrPull()
	if err != nil {
		return err
	}

	pkg.Config = viper.New()
	err = readConfig(pkg.Config, cfgName, pkg.Path)
	if err != nil {
		return err
	}
	dir := pkg.Config.GetString("dir")
	if dir != "" {
		pkg.Config.Set("dirs", append(pkg.Config.GetStringSlice("dirs"), dir))
	}
	link := pkg.Config.Get("link") // .(role.Link) // GetString
	links := pkg.Config.Get("links")
	if links == nil {
		links = make([]interface{}, 0)
	}
	if link != nil {
		// pkg.Config.Set("links", append(pkg.Config.GetStringSlice("links"), link))
		// pkg.Config.Set("links", append(pkg.Config.Get("links").(role.Link), link))
		pkg.Config.Set("links", append(links.([]interface{}), link))
	}

	// PkgConfig[pkg.Name] = pkg.Config
	return nil
}
