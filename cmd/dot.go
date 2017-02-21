package cmd

import (
	"fmt"
	// "github.com/LEI/dot/config"
	"github.com/LEI/dot/log"
	"github.com/LEI/dot/role"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"os"
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
	RolesDir   = ".dot" // Default clone parent directory
	HomeDir    = os.Getenv("HOME")
	debug      bool
	https      bool
	source     string
	target     string
	filter     []string
	logger     = log.New(os.Stdout, "", 0)
)

var (
	IgnoreNames = []string{".git", ".*\\.md"}
	Skip = fmt.Errorf("Skip")
)

var DotCmd = &cobra.Command{
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
		return installCommand(Dot.Source, Dot.Target, Dot.Roles)
	},
}

func Execute() error {
	// logger.SetOutput(os.Stdout)
	err := DotCmd.Execute()
	if err != nil {
		return err
	}
	return nil
}

func init() {
	currentDir, err := os.Getwd() // os.Getenv("PWD")
	if err != nil {
		logger.Error(err)
	}
	cobra.OnInitialize(initConfig)
	DotCmd.PersistentFlags().StringVarP(&configFile, "config", "c", configFile, "Configuration file `path`")
	DotCmd.PersistentFlags().BoolVarP(&debug, "debug", "d", debug, "Verbose output")
	DotCmd.PersistentFlags().StringSliceVarP(&filter, "filter", "f", filter, "Filter roles by `name`")
	DotCmd.PersistentFlags().BoolVarP(&https, "https", "", https, "Default to HTTPS for git remotes")
	DotCmd.PersistentFlags().StringVarP(&source, "source", "s", currentDir, "Source `directory`")
	DotCmd.PersistentFlags().StringVarP(&target, "target", "t", HomeDir, "Destination `directory`")
	// DotCmd.PersistentFlags().BoolVarP(&version, "version", "V", false, "Print the version number")
	// DotCmd.PersistentFlags().Parse(os.Args[1:])
	// Config.BindPFlags(DotCmd.Flags())
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
		Config.BindPFlag(f, DotCmd.PersistentFlags().Lookup(f))
	}
	for _, f := range bindFlags {
		Config.BindPFlag(f, DotCmd.Flags().Lookup(f))
	}
	err := Config.ReadInConfig()
	if err != nil {
		logger.Warn(err)
	}
}

func initCommand() error {
	if Config.ConfigFileUsed() != "" {
		logger.Debugln("Using config file:", Config.ConfigFileUsed())
	}
	Dot.Source = Config.GetString("source")
	Dot.Target = Config.GetString("target")
	err := Config.UnmarshalKey("roles", &Dot.Roles)
	if err != nil {
		return err
	}
	err = os.Setenv("OS", OS)
	if err != nil {
		logger.Warn(err)
	}
	OSTYPE, ok := os.LookupEnv("OSTYPE")
	if !ok {
		logger.Debugln("OSTYPE is not set")
	} else if OSTYPE == "" {
		logger.Debugln("OSTYPE is empty")
	}
	return nil
}
