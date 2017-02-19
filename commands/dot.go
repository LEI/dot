package commands

import (
	"fmt"
	"github.com/LEI/dot/conf"
	"github.com/LEI/dot/git"
	"github.com/LEI/dot/role"
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
	Config *conf.Config
	ConfigName = ".dotrc"
	Dot *role.Meta
	Https bool
)

var DotCmd = &cobra.Command{
	Use:   "dot",
	Short: "Manage dotfiles",
	Long:  ``,
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		git.Https = Https
		roles := make([]*role.Role, 0)
		Dot = &role.Meta{
			Config.GetString("source"),
			Config.GetString("target"),
			roles,
		}
		err := Config.UnmarshalKey("roles", &roles)
		if err != nil {
			return err
		}
		for _, r := range roles {
			err = r.Init(Dot.Source, Dot.Target)
			if err != nil {
				return err
			}
			ok := r.IsOs([]string{OS, OSTYPE})
			if !ok {
				continue
			}
			fmt.Println("ORIGIN", r.Origin)
			// repo, err := git.NewRepository(r.Origin)
			// if err != nil {
			// 	return err
			// }
			// r.Repo = repo
			Dot.Roles = append(Dot.Roles, r)
		}
		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) > 0 {
			fmt.Printf("Warning: extra arguments %s\n", args)
			return cmd.Help()
		}

		fmt.Println(Dot)
		for _, r := range Dot.Roles {
			// viper.Sub()
			fmt.Printf("Role %+v\n", r)
			// r.Repo.Clone()
			// r.Repo.Update()
		}
		// if err != nil {
		// 	return err
		// }
		return nil
	},
}

func Execute() {
	if err := DotCmd.Execute(); err != nil {
		fatal(err)
	}
}

func init() {
	var source, target string
	var currentDir, configFile string
	var configPaths = []string{}

	currentDir, err := os.Getwd() // os.Getenv("PWD")
	if err != nil {
		fatal(err)
	}

	DotCmd.Flags().StringVarP(&configFile, "config", "c", configFile, "Configuration `file`")
	DotCmd.Flags().BoolVarP(&Https, "https", "", Https, "Force HTTPS for git clone")
	DotCmd.Flags().StringVarP(&source, "source", "s", currentDir, "Source `directory`")
	DotCmd.Flags().StringVarP(&target, "target", "t", HOME, "Destination `directory`")
	// DotCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "verbose output")
	// DotCmd.PersistentFlags().Parse(os.Args[1:])

	if configFile != "" {
		Config.SetFile(configFile)
	} else {
		configPaths = []string{source}
		for _, p := range []string{HOME, currentDir} {
			if source != p {
				configPaths = append(configPaths, p)
			}
		}
	}
	Config = conf.New(ConfigName, configPaths)
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
