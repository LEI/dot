package cmd

import (
	"fmt"
	// "github.com/LEI/dot/config"
	"github.com/LEI/dot/log"
	"github.com/LEI/dot/role"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	// "io/ioutil"
	"os"
	"os/exec"
	"os/user"
	"runtime"
	// "time"
)

const OS = runtime.GOOS

var OSTYPE string

var (
	Dot        = &role.Meta{}
	Config     = viper.New()
	configFile = ""
	configName = ".dotrc"
	RolesDir   = ".dot" // Default clone parent directory relative to ~
	User       *user.User
	HomeDir    string // os.Getenv("HOME")
	debug      bool
	https      bool
	source     string
	target     string
	filter     []string
	logger     = log.New(os.Stdout, "", 0)
)

var (
	DotIgnore = []string{".git", "*.md", "*.tpl"}
	Skip      = fmt.Errorf("Skip")
)

var RootCmd = &cobra.Command{
	Use:   os.Args[0],
	Short: "Manage dotfiles",
	// Long:  ``,
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		// fmt.Printf("COMMAND: %+v\n%s\n", cmd, args)
		switch cmd {
		case versionCmd:
			return nil
		}
		return initCommand()
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) > 0 {
			logger.Warnln("Extra arguments:", args)
			return cmd.Help()
		}
		return installRoles(Dot.Source, Dot.Target, Dot.Roles)
	},
}

func Execute() error {
	// logger.SetOutput(os.Stdout)
	err := RootCmd.Execute()
	if err != nil {
		return err
	}
	return nil
}

func init() {
	User, err := user.Current()
	if err != nil {
		logger.Error(err)
	}
	HomeDir = User.HomeDir
	currentDir, err := os.Getwd() // os.Getenv("PWD")
	if err != nil {
		logger.Error(err)
	}
	cobra.OnInitialize(initConfig)
	RootCmd.PersistentFlags().StringVarP(&configFile, "config", "c", configFile, "Configuration file `path`")
	RootCmd.PersistentFlags().BoolVarP(&debug, "debug", "d", debug, "Verbose output")
	RootCmd.PersistentFlags().StringSliceVarP(&filter, "filter", "f", filter, "Filter roles by `name`")
	RootCmd.PersistentFlags().BoolVarP(&https, "https", "", https, "Default to HTTPS for git remotes")
	RootCmd.PersistentFlags().StringVarP(&source, "source", "s", currentDir, "Source `directory`")
	RootCmd.PersistentFlags().StringVarP(&target, "target", "t", HomeDir, "Destination `directory`")
	// RootCmd.PersistentFlags().BoolVarP(&version, "version", "V", false, "Print the version number")
	// RootCmd.PersistentFlags().Parse(os.Args[1:])
	// Config.BindPFlags(RootCmd.Flags())
}

func initConfig() {
	if debug {
		logger.SetLevel(log.DebugLevel)
	}
	bindPFlags := []string{"source", "target"} // Persistent flags
	bindFlags := []string{}                    // Local flags
	if configFile != "" {
		Config.SetConfigFile(configFile)
	} else {
		Config.SetConfigName(configName)
		configPath := []string{source, target, "/etc"}
		for _, p := range configPath {
			Config.AddConfigPath(p)
		}
	}
	for _, f := range bindPFlags {
		Config.BindPFlag(f, RootCmd.PersistentFlags().Lookup(f))
	}
	for _, f := range bindFlags {
		Config.BindPFlag(f, RootCmd.Flags().Lookup(f))
	}
	err := Config.ReadInConfig()
	switch err.(type) {
	case nil:
	case viper.ConfigFileNotFoundError:
		logger.Warn(err)
		// b, err := ioutil.ReadAll(os.Stdin)
		// if err != nil {
		// 	logger.Error(err)
		// 	return
		// }
		// in := string(b)
		// err := Config.ReadConfig(os.Stdin)
	default:
		logger.Error(err)
	}
}

func initCommand() error {
	err := os.Setenv("OS", OS)
	if err != nil {
		logger.Warn(err)
	}
	OSTYPE, ok := os.LookupEnv("OSTYPE")
	if !ok { // logger.Debugln("OSTYPE is not set")
		out, err := exec.Command("bash", "-c", "echo $OSTYPE").Output()
		if err != nil {
			return err
		}
		if len(out) > 0 {
			OSTYPE = string(out)
		}
	}
	if OSTYPE == "" {
		logger.Warnln("OSTYPE is not set or empty")
	}
	if Config.ConfigFileUsed() != "" {
		logger.Debugln("Using config file:", Config.ConfigFileUsed())
	}
	Dot.Source = Config.GetString("source")
	Dot.Target = Config.GetString("target")
	err = Config.UnmarshalKey("roles", &Dot.Roles)
	if err != nil {
		return err
	}
	return nil
}
