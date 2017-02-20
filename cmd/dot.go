package cmd

import (
	"fmt"
	// "github.com/LEI/dot/config"
	"github.com/LEI/dot/fileutil"
	"github.com/LEI/dot/git"
	"github.com/LEI/dot/log"
	"github.com/LEI/dot/role"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
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
	Config     = viper.New()
	configFile = ""
	configName = ".dotrc"
	RolesDir   = ".dot" // Default clone parent directory
	HomeDir    = os.Getenv("HOME")
	debug      bool
	source     string
	target     string
	logger        = log.New(os.Stdout, "", 0)
)

var DotCmd = &cobra.Command{
	Use:   "dot",
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
		// if len(args) > 0 {
		// 	return cmd.Help()
		// }
		return syncCommand(args)
	},
}

func Execute() error {
	// logger.SetOutput(os.Stdout)
	err := DotCmd.Execute()
	if err != nil {
		logger.Fatal(err)
	}
	return nil
}

func init() {
	currentDir, err := os.Getwd() // os.Getenv("PWD")
	if err != nil {
		logger.Error(err)
	}
	cobra.OnInitialize(initConfig)
	DotCmd.PersistentFlags().StringVarP(&configFile, "config", "c", configFile, "Configuration `file`")
	DotCmd.PersistentFlags().BoolVarP(&debug, "debug", "d", debug, "Verbose output")
	DotCmd.PersistentFlags().BoolVarP(&git.Https, "https", "", git.Https, "Force HTTPS for git clone")
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
	bindPFlags := []string{"source", "target"}
	bindFlags := []string{}
	// switch {
	// case version:
	// 	os.Args = []string{os.Args[0], "version"}
	// 	err := versionCmd.Execute()
	// 	if err != nil {
	// 		logger.Fatal("Error:", err)
	// 	}
	// 	os.Exit(0)
	// }
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
		logger.Fatal(err)
	}
}

func initCommand() error {
	if Config.ConfigFileUsed() != "" {
		logger.Debug("Using config file: " + Config.ConfigFileUsed())
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
		logger.Debug("OSTYPE is not set")
	} else if OSTYPE == "" {
		logger.Debug("OSTYPE is empty")
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
	// log := log.WithFields(logrus.Fields{
	// 	"role": r.Name,
	// })
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
	if r.Config == nil {
		r.Config = viper.New()
	}
	r.Config.SetConfigName(configName)
	r.Config.AddConfigPath(r.Source)
	err = r.Config.ReadInConfig()
	if err != nil { // && !os.IsNotExist(err)
		return err
	}
	cfgUsed := r.Config.ConfigFileUsed()
	if cfgUsed != "" {
		logger.Debug("Using role config file: " + cfgUsed)
	}

	for _, d := range r.Dirs() {
		logger.Info("- Create", d.Path)
		d.Path = os.ExpandEnv(d.Path)
		d.Path = path.Join(r.Target, d.Path)
		err := fileutil.MakeDir(d.Path) // <- fileutil.MakeDir
		if err != nil {
			return err
		}
	}
	for _, l := range r.Links() {
		logger.Info("- Symlink", l.Pattern)
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
		logger.Info("- Line in", l.File)
		l.File = os.ExpandEnv(l.File)
		l.File = path.Join(r.Target, l.File)
		err := fileutil.LineInFile(l.File, l.Line) // <- fileutil.LineInFile
		if err != nil {
			return err
		}
	}
	return nil
}
