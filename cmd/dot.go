package cmd

import (
	"fmt"
	"github.com/LEI/dot/config"
	"github.com/LEI/dot/fileutil"
	"github.com/LEI/dot/git"
	logger "github.com/LEI/dot/log"
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
	source     string
	target     string
	log        *logger.Logger
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
		// 	return cmd.Help()
		// }
		return syncCommand(args)
	},
}

func Execute() error {
	// logrus.SetFormatter(&logrus.TextFormatter{DisableTimestamp: true})
	// log.Out = os.Stdout
	// log.WithFields(logrus.Fields{
	// 	"command": "dot",
	// }).Info("Executing command...")
	err := DotCmd.Execute()
	if err != nil {
		log.Fatal(err)
	}
	return nil
}

func init() {
	currentDir, err := os.Getwd() // os.Getenv("PWD")
	if err != nil {
		log.Error(err)
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
	// if debug {
	// 	log.Level = logrus.DebugLevel
	// }
	bindPFlags := []string{"source", "target"}
	bindFlags := []string{}
	// switch {
	// case version:
	// 	os.Args = []string{os.Args[0], "version"}
	// 	err := versionCmd.Execute()
	// 	if err != nil {
	// 		log.Fatal("Error:", err)
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
		log.Fatal(err)
	}
	if configUsed != "" {
		log.Debug("Using config file: %s" + configUsed)
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
		log.Warn(err)
	}
	OSTYPE, ok := os.LookupEnv("OSTYPE")
	if !ok || OSTYPE == "" {
		log.Warn("OSTYPE is not set")
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
	cfgUsed, err := r.ReadConfig(configName, []string{r.Source})
	if err != nil && !os.IsNotExist(err) {
		return err
	}
	if cfgUsed != "" {
		log.Debug("Using role config file: " + cfgUsed)
	}

	for _, d := range r.Dirs() {
		log.Info("- Create", d.Path)
		d.Path = os.ExpandEnv(d.Path)
		d.Path = path.Join(r.Target, d.Path)
		err := fileutil.MakeDir(d.Path) // <- fileutil.MakeDir
		if err != nil {
			return err
		}
	}
	for _, l := range r.Links() {
		log.Info("- Symlink", l.Pattern)
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
		log.Info("- Line in", l.File)
		l.File = os.ExpandEnv(l.File)
		l.File = path.Join(r.Target, l.File)
		err := fileutil.LineInFile(l.File, l.Line) // <- fileutil.LineInFile
		if err != nil {
			return err
		}
	}
	return nil
}
