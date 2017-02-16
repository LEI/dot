package main

import (
	"fmt"
	"github.com/LEI/dot/cmd"
	// "github.com/LEI/dot/role"
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
	// viper = viper.New()
	// flag = pflag.NewFlagSet(os.Args[0], pflag.ExitOnError)
	// Skip = fmt.Errorf("Skip this path")
	HomeDir, CurrentDir string
	Source, Target      string
	Debug, ForceYes     bool
	ConfigFile          = ""
	ConfigName          = ".dotrc"
	IgnoreFiles         = []string{".git", ".*\\.md"}
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
	HomeDir = os.Getenv("HOME")   // user.Current().HomeDir
	CurrentDir, err := os.Getwd() // os.Getenv("PWD")
	if err != nil {
		panic(err)
	}

	// flags := RootCmd.Flags()
	flags := cmd.RootCmd.PersistentFlags()

	// flags.BoolVarP(&Sync, "sync", "S", Sync, "Synchronize packages")
	// flags.BoolVarP(&Remove, "remove", "R", Remove, "Remove packages")

	flags.StringVarP(&Source, "source", "s", CurrentDir, "Source `directory`")
	flags.StringVarP(&Target, "target", "t", HomeDir, "Destination `directory`")

	flags.BoolVarP(&Debug, "debug", "d", Debug, "Check mode")
	flags.BoolVarP(&ForceYes, "force", "f", ForceYes, "Force yes")
	flags.StringVarP(&ConfigFile, "config", "c", ConfigFile, "Configuration `file`")
	// flags.VarP(&Packages, "package", "p", "List of packages `[name=]user/repo`")

	// 	viper.SetDefault("Source", CurrentDir)
	// 	viper.SetDefault("Target", HomeDir)

	// viper.RegisterAlias("src", "source")
	// viper.RegisterAlias("dst", "target")

	viper.BindPFlags(flags)
	flags.Parse(os.Args[1:])
}

func main() {
	if ConfigFile != "" {
		viper.SetConfigFile(ConfigFile)
	}

	configPaths := []string{HomeDir, Source}
	// if CurrentDir != Source {
	// 	configPaths = append(configPaths, ".")
	// }
	err := readConfigFile(viper.GetViper(), ConfigName, configPaths...)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	// Read in environment variables that match
	viper.AutomaticEnv()

	if Debug {
		fmt.Println(Source, "->", Target)
		fmt.Printf("Viper settings:\n%s\n", viper.AllSettings())
	}

	if err := cmd.RootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(-1)
	}
}

func readConfigFile(v *viper.Viper, name string, paths ...string) error {
	v.SetConfigName(name)
	for _, path := range paths {
		v.AddConfigPath(path)
	}
	err := v.ReadInConfig()
	if err != nil && os.IsExist(err) {
		return err
	}

	configFileUsed := viper.ConfigFileUsed()
	if configFileUsed != "" {
		fmt.Printf("Using config file: %s\n", configFileUsed)
	} else {
		fmt.Println("No config file used")
	}

	return nil
}
