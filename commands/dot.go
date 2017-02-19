package commands

import (
	"fmt"
	"github.com/LEI/dot/conf"
	"github.com/LEI/dot/git"
	"github.com/LEI/dot/role"
	"github.com/spf13/cobra"
	"os"
	"path"
	"runtime"
)

const OS = runtime.GOOS

var (
	HOME   = env("HOME")
	OSTYPE = env("OSTYPE")
)

var (
	Dot        *role.Meta
	Config     *conf.Conf
	ConfigName = ".dotrc"
	debug, https, verbose bool
	defaultDotDir = ".dot"
)

var DotCmd = &cobra.Command{
	Use:   "dot",
	Short: "Manage dotfiles",
	Long:  ``,
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		return initCommand()
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		// if len(args) > 0 {
		// 	fmt.Printf("Args: %s\n", args)
		// 	return cmd.Help()
		// }
		return syncCommand(args)
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
	DotCmd.Flags().BoolVarP(&debug, "dry-run", "d", debug, "Enable check-mode")
	DotCmd.Flags().BoolVarP(&git.Https, "https", "", git.Https, "Force HTTPS for git clone")
	DotCmd.Flags().StringVarP(&source, "source", "s", currentDir, "Source `directory`")
	DotCmd.Flags().StringVarP(&target, "target", "t", HOME, "Destination `directory`")
	DotCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "verbose output")
	// DotCmd.PersistentFlags().Parse(os.Args[1:])

	if configFile != "" {
		Config = conf.NewFile(configFile)
	} else {
		configPaths = []string{source}
		for _, p := range []string{HOME, currentDir} {
			if source != p {
				configPaths = append(configPaths, p)
			}
		}
		Config = conf.New(ConfigName, configPaths)
	}
	Config.BindPFlags(DotCmd.Flags())
	configUsed, err := Config.Read()
	if err != nil && !os.IsNotExist(err) {
		fatal(err)
	}
	if verbose && configUsed != "" {
		fmt.Printf("Using config file: %s\n", configUsed)
	}
}

func initCommand() error {
	Dot = &role.Meta{
		Config.GetString("source"),
		Config.GetString("target"),
		make([]*role.Role, 0),
	}
	err := Config.UnmarshalKey("roles", &Dot.Roles)
	if err != nil {
		return err
	}
	return nil
}

func syncCommand(args []string) error {
	err := os.Setenv("OS", OS)
	if err != nil {
		fmt.Printf("Could not set env OS=%s: %s", OS, err)
	}

	for _, r := range Dot.Roles {
		r, err = r.New(Dot.Source, Dot.Target)
		if err != nil {
			return err
		}
		// Check platform
		ok := r.IsOs([]string{OS, OSTYPE})
		if !ok {
			continue
		}
		// Filter roles by name
		skip := len(args) > 0
		for _, arg := range args {
			if arg == r.Name {
				skip = false
				break
			}
		}
		if skip {
			continue
		}

		fmt.Printf("### Role %s\n", r.Name)
		if verbose {
			fmt.Println(" ->", r)
		}

		// TODO func (r *Role) NewRepo() error?
		git.DefaultPath = path.Join(r.Target, defaultDotDir)
		repo, err := git.NewRepository(r.Origin)
		if err != nil {
			return err
		}
		fmt.Println("Changing", repo.Name, "to", r.Name, "--", repo)
		repo.Name = r.Name
		err = repo.CloneOrPull()
		if err != nil {
			return err
		}
		if repo.Path != r.Source {
			r.Source = repo.Path
		}
		cfgUsed, err := r.ReadConfig(ConfigName, []string{r.Source})
		if err != nil && !os.IsNotExist(err) {
			return err
		}
		if verbose && cfgUsed != "" {
			fmt.Printf("Using role config file: %s\n", cfgUsed)
		}

		for _, dir := range r.Dirs() {
			fmt.Printf("## Create %s\n", dir)
			err := dir.Create(r.Target)
			if err != nil {
				return err
			}
		}
		for _, link := range r.Links() {
			fmt.Printf("## Symlink %s\n", link)
			err := link.Sync(r.Source, r.Target)
			if err != nil {
				return err
			}
		}
		for _, line := range r.Lines() {
			fmt.Printf("## Line in file %s\n", line)
			err := line.InFile(r.Target)
			if err != nil {
				return err
			}
		}
	}
	return nil
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
