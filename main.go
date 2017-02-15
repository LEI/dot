package main

import (
	"fmt"
	"log"
	"github.com/spf13/pflag"
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
)

var flag = pflag.NewFlagSet(os.Args[0], pflag.ExitOnError)

func init() {
	flag.BoolVarP(&Sync, "sync", "S", Sync, "Synchronize packages")
	flag.BoolVarP(&Remove, "remove", "R", Remove, "Remove packages")

	flag.StringVarP(&Source, "source", "s", CurrentDir, "Source `directory`")
	flag.StringVarP(&Target, "target", "t", HomeDir, "Destination `directory`")

	flag.BoolVarP(&Debug, "debug", "d", Debug, "Check mode")
	flag.BoolVarP(&ForceYes, "force", "f", ForceYes, "Force yes")

	flag.StringVarP(&ConfigFile, "config", "c", ConfigFile, "Configuration `file`")

	flag.VarP(&Packages, "package", "p", "List of packages: `[path=]user/repo`")

	flag.Parse(os.Args[1:])
}

func main() {
	// 	viper.SetDefault("Source", CurrentDir)
	// 	viper.SetDefault("Target", HomeDir)

	viper.RegisterAlias("src", "source")
	viper.RegisterAlias("dst", "target")

	viper.BindPFlags(flag)

	viper.SetConfigName(ConfigName)
	viper.AddConfigPath(HomeDir)
	viper.AddConfigPath(Source)
	if CurrentDir != Source {
		viper.AddConfigPath(".")
	}

	viper.SetConfigFile(ConfigFile)
	err := viper.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("Fatal error reading configuration: %s\n", err))
	}

	configFile := viper.ConfigFileUsed()
	fmt.Printf("Used config file: %s\n", configFile)

	// Packages = viper.Get("packages").(PackageSlice)

	fmt.Println(Source, "->", Target)
	if Debug {
		fmt.Printf("%s", viper.AllSettings())
	}
	fmt.Printf("%s: %+v\n", "Packages", viper.Get("packages"))
	fmt.Println(viper.Get("dir"), viper.Get("dirs"))
	fmt.Println(viper.Get("link"), viper.Get("links"))
	fmt.Println(viper.Get("line"), viper.Get("lines"))
}
