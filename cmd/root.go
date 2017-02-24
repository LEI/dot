package cmd

import (
	// "github.com/LEI/dot/config"
	dot "github.com/LEI/dot/dotfile"
	"github.com/LEI/dot/git"
	"github.com/LEI/dot/log"
	"github.com/LEI/dot/role"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	// "io/ioutil"
	"os"
	"os/exec"
	"os/user"
	"path"
	"runtime"
	"strings"
	// "time"
)

const OS = runtime.GOOS

var (
	OsFamily   = []string{OS}
	Dot        = &role.Meta{}
	DotIgnore  = []string{".git", "*.md", "*.tpl"}
	Config     = viper.New()
	configFile = ""
	configName = ".dotrc"
	RolesDir   = ".dot" // Default clone parent directory relative to ~
	User       *user.User
	HomeDir    string // os.Getenv("HOME")
	currentDir string // os.Getenv("PWD")
	roleFilter []string
	debug      bool
	dryrun     bool //= true
	https      bool
	logger     = log.New(os.Stdout, "", 0)
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
		return installRoles()
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
	err := os.Setenv("OS", OS)
	if err != nil {
		logger.Warn(err)
	}
	osType, ok := os.LookupEnv("OSTYPE")
	logger.Debugf("OSTYPE='%s' (%v)", osType, ok)
	if !ok || osType == "" {
		out, err := exec.Command("bash", "-c", "printf '%s' \"$OSTYPE\"").Output()
		if err != nil {
			logger.Error(err)
		}
		if len(out) > 0 {
			osType = string(out)
			ot := strings.Split(osType, ".")
			if len(ot) > 0 {
				OsFamily = append(OsFamily, ot[0])
			}
			OsFamily = append(OsFamily, osType)
		}
	}
	if osType != "" {
	} else {
		logger.Warnln("OSTYPE is not set or empty")
	}
	User, err = user.Current()
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
	RootCmd.PersistentFlags().BoolVarP(&dryrun, "dry-run", "D", dryrun, "Check-mode")
	RootCmd.PersistentFlags().StringSliceVarP(&roleFilter, "filter", "f", roleFilter, "Filter roles by `name`")
	RootCmd.PersistentFlags().BoolVarP(&https, "https", "", https, "Default to HTTPS for git remotes")
	RootCmd.PersistentFlags().StringVarP(&Dot.Source, "source", "s", currentDir, "Dot.Source `directory`")
	RootCmd.PersistentFlags().StringVarP(&Dot.Target, "target", "t", HomeDir, "Destination `directory`")
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
		configPath := []string{Dot.Source, Dot.Target, "/etc"}
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
	if Config.ConfigFileUsed() != "" {
		logger.Debugln("Using config file:", Config.ConfigFileUsed())
	}
	Dot.Source = Config.GetString("source")
	Dot.Target = Config.GetString("target")
	err := Config.UnmarshalKey("roles", &Dot.Roles)
	if err != nil {
		return err
	}
	if len(Dot.Roles) == 0 {
		logger.Warnln("No role found")
	}
	return nil
}

func apply(handlers ...func(*role.Role) error) error {
ROLES:
	for _, r := range Dot.Roles {
		r, err := r.New(Dot.Source, Dot.Target)
		if err != nil {
			return err
		}
		for _, f := range handlers {
			err := f(r)
			if err != nil {
				switch err {
				case dot.Skip:
					continue ROLES
				}
				return err
			}
		}
	}
	return nil
}

func validateRole(r *role.Role) error {
	title := strings.Title(r.Name)
	// Check platform
	ok := r.IsOs(OsFamily)
	if !ok {
		logger.Debugf("## %s (only for %s)\n", title, strings.Join(r.Os, ", "))
		return dot.Skip
	}
	// Filter by name
	skip := len(roleFilter) > 0
	for _, roleName := range roleFilter {
		if roleName == r.Name {
			skip = false
			break
		}
	}
	if skip {
		logger.Debugf("## %s (skipped)\n", title)
		return dot.Skip
	}
	// logger.SetPrefix(r.Name+": ") // ctx.Value("role")
	logger.Infof("## %s\n", title)
	return nil
}

func initGitRepo(r *role.Role) error {
	dir := path.Join(r.Target, RolesDir, r.Name) // git.DefaultPath
	git.Https = https
	repo, err := git.New(r.Origin, dir)
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
	return nil
}

func initRoleConfig(r *role.Role) error {
	if r.Config == nil {
		r.Config = viper.New()
	}
	r.Config.SetConfigName(configName)
	r.Config.AddConfigPath(r.Source)
	err := r.Config.ReadInConfig()
	if err != nil { // && !os.IsNotExist(err)
		return err
	}
	cfgUsed := r.Config.ConfigFileUsed()
	if cfgUsed != "" {
		logger.Debugln("Using role config file:", cfgUsed)
	}
	return nil
}

func only(t string) dot.FileHandler {
	return func(path string, fi os.FileInfo) error {
		switch {
		// case fi == nil:
		// 	logger.Debugf("# ignore %s (does not exist)\n", fi.Name(), t)
		// 	return dot.Skip
		case t == "directory" && !fi.IsDir(),
			 t == "file" && fi.IsDir():
			logger.Debugf("# ignore %s (not a %s)\n", fi.Name(), t)
			return dot.Skip
		// case t == "":
		// 	logger.Errorf("! Invalid type: %s", t)
		// default:
		// 	logger.Debugln("default", fi.Name(), t)
		}
		return nil
	}
}

func filterIgnored(path string, fi os.FileInfo) error {
	ignore, err := dot.Match(fi.Name(), DotIgnore...)
	if err != nil {
		return err
	}
	if ignore {
		logger.Debugf("Ignore %s\n", fi.Name())
		return dot.Skip
	}
	return nil
}
