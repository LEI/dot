package cmd

import (
	"fmt"
	// "github.com/LEI/dot/cmd"
	"github.com/LEI/dot/role"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"log"
	"os"
	"regexp"
	"runtime"
	"strings"
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
	Sync, Remove        bool
	HomeDir, CurrentDir string
	Source, Target      string
	Debug, ForceYes     bool
	ConfigFile          = ""
	ConfigName          = ".dotrc"
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
	// cobra.OnInitialize(initConfig)

	HomeDir = os.Getenv("HOME")   // user.Current().HomeDir
	CurrentDir, err := os.Getwd() // os.Getenv("PWD")
	if err != nil {
		fatal(err)
	}

	// RootCmd.Flags().BoolVarP(&Clone, "clone", "C", Clone, "Clone remote packages")
	RootCmd.Flags().BoolVarP(&Sync, "sync", "S", Sync, "Synchronize packages")
	RootCmd.Flags().BoolVarP(&Remove, "remove", "R", Remove, "Remove packages")

	RootCmd.PersistentFlags().StringVarP(&Source, "source", "s", CurrentDir, "Source `directory`")
	RootCmd.PersistentFlags().StringVarP(&Target, "target", "t", HomeDir, "Destination `directory`")

	RootCmd.PersistentFlags().BoolVarP(&Debug, "debug", "d", Debug, "Check mode")
	RootCmd.PersistentFlags().BoolVarP(&ForceYes, "force-yes", "f", ForceYes, "Force yes")

	RootCmd.PersistentFlags().StringVarP(&ConfigFile, "config", "c", ConfigFile, "Configuration `file`")
	RootCmd.PersistentFlags().StringVarP(&ConfigName, "config-name", "", ConfigName, "Configuration `file`")
	// RootCmd.PersistentFlags().VarP(&Packages, "package", "p", "List of packages `[name=]user/repo`")

	// 	viper.SetDefault("Source", CurrentDir)
	// 	viper.SetDefault("Target", HomeDir)

	// viper.RegisterAlias("src", "source")
	// viper.RegisterAlias("dst", "target")

	Config.BindPFlags(RootCmd.PersistentFlags())
	// RootCmd.PersistentFlags().Parse(os.Args[1:])

	initConfig()
}

var RootCmd = &cobra.Command{
	Use:   "dot",
	Short: "Manage dotfiles",
	Long:  ``,
	PersistentPreRun: func(cmd *cobra.Command, args[]string) {
		Source = Config.GetString("source")
		Target = Config.GetString("target")
		// Debug = Config.GetString("debug")
		// ForceYes = Config.GetString("force-yes")
		// ConfigFile = Config.GetString("config")
		// ConfigName = Config.GetString("config-name")
		err := Config.UnmarshalKey("packages", &Packages)
		if err != nil {
			fatal(err)
		}
	},
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) > 0 {
			fatal(fmt.Errorf("Extra arguments: %s", args))
		}
		switch {
		case Remove:
			if err := flagToArg("R", "remove"); err != nil {
				fatal(err)
			}
			if err := removeCmd.Execute(); err != nil {
				fatal(err)
			}
		case Sync:
			if err := flagToArg("S", "sync"); err != nil {
				fatal(err)
			}
			if err := syncCmd.Execute(); err != nil {
				fatal(err)
			}
		default:
			err := cmd.Help()
			if err != nil {
				fatal(err)
			}
		}
	},
}

// Only handle long and uppercase short flags
func flagToArg(short string, long string) error {
	var new = os.Args // []string{os.Args[0], os.Args[1]}
	re := regexp.MustCompile("^-[a-z]*" + short + "[a-z]*$")
	for i, arg := range os.Args {
		if i == 0 {
			continue
		}
		if arg == "-" + short || arg == "--" + long {
			os.Args[i] = long
			return nil
		}
		matched := re.MatchString(arg)
		if matched {
			new[i] = strings.Replace(arg, short, "", 1)
			new = append(new, long)
			break
		}
	}
	if len(new) > len(os.Args) {
		os.Args = new
		return nil
	}
	return fmt.Errorf("Did not replace -%s or --%s in args: '%+v'", short, long, os.Args)

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
	// fmt.Println("INIT CONFIG")

	if ConfigFile != "" {
		Config.SetConfigFile(ConfigFile)
	}

	configPaths := []string{HomeDir, Source}
	// if CurrentDir != Source {
	// 	configPaths = append(configPaths, ".")
	// }
	// viper.GetViper()
	err := readConfig(Config, ConfigName, configPaths...)
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
	cfgPath := v.ConfigFileUsed()
	if cfgPath != "" {
		fmt.Printf("Using: %s\n", cfgPath)
		if Debug {
			fmt.Printf("%s >>> %+v\n", name, v)
		}
	} else {
		fmt.Println(name, "not found in", paths)
	}
	return err
}

func fatal(msg interface{}) {
	// log.Fatal*
	fmt.Fprintln(os.Stderr, "Error:", msg)
	os.Exit(1)
}
