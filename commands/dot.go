package commands

import (
	"fmt"
	"github.com/LEI/dot/conf"
	"github.com/spf13/cobra"
	"os"
	"runtime"
)

const OS = runtime.GOOS

var (
	HOME = env("HOME")
	OSTYPE = env("OSTYPE")
)

var (
	CurrentDir string
	Config *conf.Config
	ConfigFile string
	ConfigName = ".dotrc"
	source, target string
)

var DotCmd = &cobra.Command{
	Use:   "dot",
	Short: "Manage dotfiles",
	Long:  ``,
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		source = Config.GetString("source")
		target = Config.GetString("target")
		// err := Config.UnmarshalKey("packages", &Packages)
		// if err != nil {
		// 	fatal(err)
		// }
		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) > 0 {
			fmt.Printf("Extra arguments: %s", args)
		}
		fmt.Println(source, target)
		return cmd.Help()
		// if err != nil {
		// 	return err
		// }
		// return nil
	},
}

func Execute() {
	if err := DotCmd.Execute(); err != nil {
		fatal(err)
	}
}

func init() {
	CurrentDir, err := os.Getwd() // os.Getenv("PWD")
	if err != nil {
		fatal(err)
	}

	DotCmd.Flags().StringVarP(&ConfigFile, "config", "c", ConfigFile, "Configuration `file`")
	DotCmd.Flags().StringVarP(&source, "source", "s", CurrentDir, "Source `directory`")
	DotCmd.Flags().StringVarP(&target, "target", "t", HOME, "Destination `directory`")
	// DotCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "verbose output")
	// DotCmd.PersistentFlags().Parse(os.Args[1:])

	Config = conf.New(ConfigName, []string{source, CurrentDir, HOME})
	if ConfigFile != "" {
		Config.SetFile(ConfigFile)
	}
	Config.BindPFlags(DotCmd.Flags())
	configUsed, err := Config.Read()
	if err != nil && !os.IsNotExist(err) {
		fatal(err)
	}
	if configUsed != "" {
		fmt.Printf("Using config file: %s\n", configUsed)
	}
}

func env(key string) string {
	val, ok := os.LookupEnv(key)
	if !ok {
		fmt.Fprintf(os.Stderr, "Warning: %s is not set\n", key)
	}
	return val
}

func fatal(msg interface{}) {
	// log.Fatal*
	fmt.Fprintf(os.Stderr, "Error: %s\n", msg)
	os.Exit(1)
}
