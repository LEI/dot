package cmd

import (
	"fmt"
	"github.com/LEI/dot/config"
	"github.com/LEI/dot/fileutil"
	"github.com/LEI/dot/git"
	"github.com/LEI/dot/role"
	"github.com/spf13/cobra"
	"os"
	"path"
	"runtime"
	"strings"
	// "time"
)

const OS = runtime.GOOS

var OSTYPE string

var (
	Dot        = &role.Meta{}
	Config     *config.Configuration
	configFile = ""
	configName = ".dotrc"
	RolesDir   = ".dot" // Default clone parent directory
	HomeDir    = os.Getenv("HOME")
	debug      bool
	verbose    bool
	source     string
	target     string
)

var DotCmd = &cobra.Command{
	Use:   "dot",
	Short: "Manage dotfiles",
	// Long:  ``,
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

func Execute() error {
	err := DotCmd.Execute()
	if err != nil {
		if verbose {
			fmt.Fprintf(os.Stderr, "Fatal error: %+v\n", err)
		}
		os.Exit(1)
	}
	return nil
}

func init() {
	currentDir, err := os.Getwd() // os.Getenv("PWD")
	if err != nil {
		fmt.Fprintf(os.Stderr, "os.Getcwd() Error: %s\n", err)
	}
	cobra.OnInitialize(initConfig)
	DotCmd.PersistentFlags().StringVarP(&configFile, "config", "c", configFile, "Configuration `file`")
	DotCmd.PersistentFlags().BoolVarP(&debug, "dry-run", "d", debug, "Enable check-mode")
	DotCmd.PersistentFlags().BoolVarP(&git.Https, "https", "", git.Https, "Force HTTPS for git clone")
	DotCmd.PersistentFlags().StringVarP(&source, "source", "s", currentDir, "Source `directory`")
	DotCmd.PersistentFlags().StringVarP(&target, "target", "t", HomeDir, "Destination `directory`")
	DotCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "Verbose output")
	// DotCmd.PersistentFlags().BoolVarP(&version, "version", "V", false, "Print the version number")
	// DotCmd.PersistentFlags().Parse(os.Args[1:])
	// Config.BindPFlags(DotCmd.Flags())
}

func initConfig() {
	bindPFlags := []string{"source", "target"}
	bindFlags := []string{}
	// switch {
	// case version:
	// 	os.Args = []string{os.Args[0], "version"}
	// 	err := versionCmd.Execute()
	// 	if err != nil {
	// 		fmt.Println("Error:", err)
	// 		os.Exit(1)
	// 	}
	// 	os.Exit(0)
	// }
	if configFile != "" {
		Config = config.NewFile(configFile)
	} else {
		configPath := []string{source, target, "/etc"}
		Config = config.New(configName, configPath)
	}
	for _, f := range bindPFlags {
		Config.BindPFlag(f, DotCmd.PersistentFlags().Lookup(f))
	}
	for _, f := range bindFlags {
		Config.BindPFlag(f, DotCmd.Flags().Lookup(f))
	}
	configUsed, err := Config.Read()
	if err != nil {
		fmt.Println("read error", os.IsNotExist(err))
		fmt.Fprintf(os.Stderr, "Read config error: %s\n", err)
		os.Exit(1)
	}
	if verbose && configUsed != "" {
		fmt.Printf("Using config file: %s\n", configUsed)
	}
}

func initCommand() error {
	Dot.Source = Config.GetString("source")
	Dot.Target = Config.GetString("target")
	err := Config.UnmarshalKey("roles", &Dot.Roles)
	if err != nil {
		return err
	}
	err = os.Setenv("OS", OS)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Warning: could not set OS to %s: %s\n", OS, err)
	}
	OSTYPE, ok := os.LookupEnv("OSTYPE")
	if !ok || OSTYPE == "" {
		fmt.Fprintln(os.Stderr, "Warning: OSTYPE is not set")
	}
	return nil
}

func syncCommand(args []string) error {
	for _, r := range Dot.Roles {
		r, err := r.New(Dot.Source, Dot.Target)
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
		err = syncRole(r)
		if err != nil {
			return err
		}
	}
	return nil
}

func syncRole(r *role.Role) error {
	fmt.Printf("--- Role %s\n", r.Name)
	// defer fmt.Printf("---\n")

	// TODO func (r *Role) NewRepo() error?
	git.DefaultPath = path.Join(r.Target, RolesDir)
	repo, err := git.NewRepository(r.Origin)
	if err != nil {
		return err
	}
	repo.Name = r.Name
	err = repo.CloneOrPull()
	if err != nil {
		return err
	}
	if repo.Path != r.Source {
		r.Source = repo.Path
	}
	cfgUsed, err := r.ReadConfig(configName, []string{r.Source})
	if err != nil && !os.IsNotExist(err) {
		return err
	}
	if verbose && cfgUsed != "" {
		fmt.Printf("Using role config file: %s\n", cfgUsed)
	}

	for _, d := range r.Dirs() {
		fmt.Println("- Create", d.Path)
		if verbose {
			fmt.Println("->", d)
		}
		d.Path = os.ExpandEnv(d.Path)
		d.Path = path.Join(r.Target, d.Path)
		err := fileutil.MakeDir(d.Path) // <- fileutil.MakeDir
		if err != nil {
			return err
		}
	}
	for _, l := range r.Links() {
		fmt.Println("- Symlink", l.Pattern)
		if verbose {
			fmt.Println("->", l)
		}
		l.Pattern = os.ExpandEnv(l.Pattern)
		paths, err := l.GlobFiles(r.Source) // <- role.Link.GlobFiles(src string)
		if err != nil {
			return err
		}
		for _, src := range paths {
			dst := strings.Replace(src, r.Source, r.Target, 1)
			err := fileutil.Symlink(src, dst) // <- fileutil.Symlink
			if err != nil {
				return err
			}
		}
	}
	for _, l := range r.Lines() {
		fmt.Println("- Line in", l.File)
		if verbose {
			fmt.Println("->", l)
		}
		l.File = os.ExpandEnv(l.File)
		l.File = path.Join(r.Target, l.File)
		err := fileutil.LineInFile(l.File, l.Line) // <- fileutil.LineInFile
		if err != nil {
			return err
		}
	}
	return nil
}
