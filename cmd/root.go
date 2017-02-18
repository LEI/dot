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
	OS = runtime.GOOS
	version = "master"
)

var ( // os.Getenv("HOME") // user.Current().HomeDir
	HOME = env("HOME")
	OSTYPE = env("OSTYPE")
	OsTypes = []string{OS, OSTYPE}
	CurrentDir string
	Source     string
	Target     string
	Debug      bool
	AssumeYes  bool
	Https      bool
	ConfigFile = ""
	ConfigName = ".dotrc"
	Config = viper.New()
	Packages   []*role.Package
	// Packages            role.PackageSlice //= make(role.PackageSlice, 0)
	// PkgConfig = make(map[string]*viper.Viper, 0)
	// flag = pflag.NewFlagSet(os.Args[0], pflag.ExitOnError)
	// Skip = fmt.Errorf("Skip this path")
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
	fmt.Printf("OS: %+v\n", OsTypes)
	// cobra.OnInitialize(initConfig)

	err := os.Setenv("OS", OS)
	if err != nil {
		fmt.Printf("Could not set env OS=%s: %s", OS, err)
	}
	CurrentDir, err := os.Getwd() // os.Getenv("PWD")
	if err != nil {
		fatal(err)
	}

	// RootCmd.Flags().BoolVarP(&Clone, "clone", "C", Clone, "Clone remote packages")
	// RootCmd.Flags().BoolVarP(&Sync, "sync", "S", Sync, "Synchronize packages")
	// RootCmd.Flags().BoolVarP(&Remove, "remove", "R", Remove, "Remove packages")

	RootCmd.PersistentFlags().StringVarP(&Source, "source", "s", CurrentDir, "Source `directory`")
	RootCmd.PersistentFlags().StringVarP(&Target, "target", "t", HOME, "Destination `directory`")

	RootCmd.PersistentFlags().BoolVarP(&Debug, "debug", "d", Debug, "Check mode")
	RootCmd.PersistentFlags().BoolVarP(&AssumeYes, "assume-yes", "y", AssumeYes, "Force yes")

	RootCmd.PersistentFlags().StringVarP(&ConfigFile, "config", "c", ConfigFile, "Configuration `file`")
	RootCmd.PersistentFlags().StringVarP(&ConfigName, "config-name", "", ConfigName, "Configuration `file`")
	// RootCmd.PersistentFlags().VarP(&Packages, "package", "p", "List of packages `[name=]user/repo`")

	RootCmd.PersistentFlags().BoolVarP(&Https, "https", "", false, "Force HTTPS for git clone")

	Config.BindPFlags(RootCmd.PersistentFlags())

	// 	viper.SetDefault("Source", CurrentDir)
	// 	viper.SetDefault("Target", HOME)

	// viper.RegisterAlias("src", "source")
	// viper.RegisterAlias("dst", "target")

	RootCmd.PersistentFlags().Parse(os.Args[1:])

	initConfig()
}

var RootCmd = &cobra.Command{
	Use:   "dot",
	Short: "Manage dotfiles",
	Long:  ``,
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		Source = Config.GetString("source")
		Target = Config.GetString("target")
		// Debug = Config.GetString("debug")
		// AssumeYes = Config.GetString("force-yes")
		// ConfigFile = Config.GetString("config")
		// ConfigName = Config.GetString("config-name")
		err := Config.UnmarshalKey("packages", &Packages)
		if err != nil {
			fatal(err)
		}
	},
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) > 0 {
			fmt.Printf("Extra arguments: %s", args)
		}
		err := cmd.Help()
		if err != nil {
			fatal(err)
		}
	},
}

func flagToCmd(short string, long string) error {
	var new = []string{os.Args[0]}
	re := regexp.MustCompile("^-[a-zA-Z]*" + short + "[a-zA-Z]*$")
	for i, arg := range os.Args {
		if i == 0 {
			continue
		}
		if arg == "-"+short || arg == "--"+long {
			os.Args[i] = long
			return nil
		}
		matched := re.MatchString(arg)
		if matched {
			new = append(new, long)
			os.Args[i] = strings.Replace(arg, short, "", 1)
			for _, a := range os.Args[1:] {
				new = append(new, a)
			}
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

func env(key string) string {
	val, ok := os.LookupEnv(key)
	if !ok {
		fmt.Printf("%s is not set!\n", val)
	}
	return val
}

func initConfig() {
	err := readConfig(Config)
	if err != nil && os.IsExist(err) {
		fatal(err)
	}
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
	}
}

func readConfig(c *viper.Viper) error {
	var configPaths []string
	if ConfigFile != "" {
		c.SetConfigFile(ConfigFile)
	} else { // ConfigName != ""
		c.SetConfigName(ConfigName)
		configPaths = []string{Source}
		for _, dir := range []string{HOME, CurrentDir} {
			if Source != dir {
				configPaths = append(configPaths, dir)
			}
		}
		for _, path := range configPaths {
			c.AddConfigPath(path)
		}
	}
	// Read in environment variables that match
	c.AutomaticEnv()
	err := c.ReadInConfig()
	configUsed := c.ConfigFileUsed()
	if configUsed != "" {
		fmt.Printf("Using: %s\n", configUsed)
		// if Debug {
		// 	fmt.Printf("%s >>> %+v\n", ConfigName, configPaths)
		// }
	} else if ConfigFile == "" {
		fmt.Println(ConfigName, "not found in", configPaths)
	}
	return err
}

func fatal(msg interface{}) {
	// log.Fatal*
	fmt.Fprintln(os.Stderr, "Error:", msg)
	os.Exit(1)
}
