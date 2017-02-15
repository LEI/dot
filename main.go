package main

import (
	"fmt"
	"log"
	"github.com/LEI/dot/cmd"
	"github.com/spf13/viper"
	"os"
	"runtime"
)

const (
	OS      = runtime.GOOS
	version = "master"
)

var (
	HomeDir         = os.Getenv("HOME")
	CurrentDir      = os.Getenv("PWD")
	Sync            = true
	Remove          = false
	Source, Target  string
	Debug, ForceYes bool
	ConfigFile      = ""
	ConfigName      = ".dotrc"
	IgnoreFiles     = []string{".git", ".*\\.md"}
	Packages        = make(PackageSlice, 0)
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
	// Skip = fmt.Errorf("Skip this path")
	// viper = viper.New()
	// flag = pflag.NewFlagSet(os.Args[0], pflag.ExitOnError)
)

func init() {
	flags := cmd.RootCmd.Flags()
	gFlags := cmd.RootCmd.PersistentFlags()

	flags.BoolVarP(&Sync, "sync", "S", Sync, "Synchronize packages")
	flags.BoolVarP(&Remove, "remove", "R", Remove, "Remove packages")

	gFlags.StringVarP(&Source, "source", "s", CurrentDir, "Source `directory`")
	gFlags.StringVarP(&Target, "target", "t", HomeDir, "Destination `directory`")

	flags.BoolVarP(&Debug, "debug", "d", Debug, "Check mode")
	flags.BoolVarP(&ForceYes, "force", "f", ForceYes, "Force yes")
	flags.StringVarP(&ConfigFile, "config", "c", ConfigFile, "Configuration `file`")
	flags.VarP(&Packages, "package", "p", "List of packages `[name=]user/repo`")

	// Parse command line arguments
	flags.Parse(os.Args[1:])

	// 	viper.SetDefault("Source", CurrentDir)
	// 	viper.SetDefault("Target", HomeDir)

	// viper.RegisterAlias("src", "source")
	// viper.RegisterAlias("dst", "target")

	// Bind values to configuration
	viper.BindPFlags(flags)
	viper.BindPFlags(gFlags)

	if ConfigFile != "" {
		viper.SetConfigFile(ConfigFile)
	}
	configPaths := []string{HomeDir, Source}
	// if CurrentDir != Source {
	// 	configPaths = append(configPaths, ".")
	// }
	err := ReadConfigFile(viper.GetViper(), ConfigName, configPaths...)
	if err != nil {
		panic(err)
	}
}

func main() {
	configFileUsed := viper.ConfigFileUsed()
	if configFileUsed != "" {
		fmt.Printf("Used config file: %s\n", configFileUsed)
	} else {
		fmt.Printf("No config file found")
	}

	if Debug {
		fmt.Println(Source, "->", Target)
		fmt.Printf("Viper settings:\n%s\n", viper.AllSettings())
	}

	if err := cmd.RootCmd.Execute(); err != nil {
		panic(err)
	}
}

func ReadConfigFile(v *viper.Viper, name string, paths ...string) error {
	v.SetConfigName(name)
	for _, path := range paths {
		v.AddConfigPath(path)
	}
	err := v.ReadInConfig()
	if err != nil {
		return fmt.Errorf("Fatal error reading configuration: %s\n", err)
	}
	return nil
}
