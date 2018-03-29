package cmd

import (
	"os"
	"os/exec"
	"os/user"
	"path"
	// "reflect"
	"runtime"
	"strings"

	dot "github.com/LEI/dot/dotfile"
	"github.com/LEI/dot/git"
	"github.com/LEI/dot/log"
	"github.com/LEI/dot/role"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const (
	OS = runtime.GOOS
	shell = "bash"
)

var (
	User       *user.User
	HomeDir    string // os.Getenv("HOME")
	CurrentDir string // os.Getenv("PWD")
)

var (
	Dot        = &role.Meta{}
	DotIgnore  = []string{".git", "*.md", "*.tpl", "*.yml", "*.yaml"}
	// , "*.json" (TODO: allow overwrite ignore, e.g. karabiner.json)
	DotExclude = []string{}
	Config     = viper.New()
	configFile = ""
	configName = ".dotrc"
	RolesDir   = ".dot" // Default clone parent directory relative to ~
	// User       *user.User
	// HomeDir    string // os.Getenv("HOME")
	// CurrentDir string // os.Getenv("PWD")
	roleFilter []string
	debug      bool
)

var logger = log.New(os.Stdout, "", 0)

var DotCmd = &cobra.Command{
	Use:   os.Args[0],
	Short: "Manage dotfiles",
	// Long:  ``,
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		switch cmd {
		case versionCmd:
			return nil
		}
		return initRoles()
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

	AddCommands()

	err := DotCmd.Execute()
	if err != nil {
		return err
	}
	return nil
}

func AddCommands() {
	DotCmd.AddCommand(installCmd)
	DotCmd.AddCommand(removeCmd)
	DotCmd.AddCommand(versionCmd)
}

func init() {
	err := os.Setenv("OS", OS)
	if err != nil {
		logger.Warn(err)
	}
	cobra.OnInitialize(initConfig)
	initPersistentFlags(DotCmd)
}

func initPersistentFlags(cmd *cobra.Command) {
	User, err := user.Current()
	if err != nil {
		logger.Error(err)
	}
	HomeDir = User.HomeDir
	CurrentDir, err := os.Getwd() // os.Getenv("PWD")
	if err != nil {
		logger.Error(err)
	}
	DotCmd.PersistentFlags().StringVarP(&configFile, "config", "c", configFile, "Configuration file `path`")
	DotCmd.PersistentFlags().BoolVarP(&debug, "debug", "d", debug, "Verbose output")
	DotCmd.PersistentFlags().BoolVarP(&dot.DryRun, "dry-run", "D", dot.DryRun, "Check-mode")
	DotCmd.PersistentFlags().StringSliceVarP(&DotExclude, "exclude", "x", DotExclude, "Exclude `pattern`")
	DotCmd.PersistentFlags().StringSliceVarP(&roleFilter, "filter", "f", roleFilter, "Filter roles by `name`")
	DotCmd.PersistentFlags().BoolVarP(&git.Https, "https", "", git.Https, "Default to HTTPS for git remotes")
	DotCmd.PersistentFlags().StringVarP(&Dot.Source, "source", "s", CurrentDir, "Source `directory`")
	DotCmd.PersistentFlags().StringVarP(&Dot.Target, "target", "t", HomeDir, "Destination `directory`")
	// DotCmd.PersistentFlags().BoolVarP(&version, "version", "V", false, "Print the version number")
	// DotCmd.PersistentFlags().Parse(os.Args[1:])
	// Config.BindPFlags(DotCmd.Flags())
}

func initConfig() {
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
		Config.BindPFlag(f, DotCmd.PersistentFlags().Lookup(f))
	}
	for _, f := range bindFlags {
		Config.BindPFlag(f, DotCmd.Flags().Lookup(f))
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

func osTypes() []string {
	types := []string{OS}
	OSTYPE, ok := os.LookupEnv("OSTYPE")
	logger.Debugf("OSTYPE='%s' (%v)\n", OSTYPE, ok)
	if !ok || OSTYPE == "" {
		out, err := exec.Command(shell, "-c", "printf '%s' \"$OSTYPE\"").Output()
		if err != nil {
			logger.Error(err)
		}
		if len(out) > 0 {
			OSTYPE = string(out)
			ost := strings.Split(OSTYPE, ".")
			if len(ost) > 0 {
				types = append(types, ost[0])
			}
			types = append(types, OSTYPE)
		}
	}
	if OSTYPE == "" {
		logger.Warnln("OSTYPE is not set or empty")
	}
	return types
}

func initRoles() error {
	configUsed := Config.ConfigFileUsed()
	if configUsed != "" {
		logger.Debugln("Using config file:", configUsed)
	}
	Dot.Source = Config.GetString("source")
	Dot.Target = Config.GetString("target")
	err := Config.UnmarshalKey("roles", &Dot.Roles)
	if err != nil {
		return err
	}
	if len(Dot.Roles) == 0 {
		logger.Errorf("No roles found in %s\n", configUsed)
	}
	roles := make([]*role.Role, 0)
	for _, r := range Dot.Roles {
		r, err := r.Init(Dot.Source, Dot.Target)
		if err != nil {
			return err
		}

		// Check platform
		ok := r.IsOs(osTypes())
		if !ok {
			logger.Debugf("## %s (only for %s)\n", r.Title(), strings.Join(r.Os, ", "))
			continue // return dot.Skip
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
			logger.Debugf("## %s (skipped)\n", r.Title())
			continue // return dot.Skip
		}
		// logger.SetPrefix(r.Name+": ") // ctx.Value("role")
		err = initRoleRepo(r)
		if err != nil {
			return err
		}

		roles = append(roles, r)
	}
	Dot.Roles = roles
	return nil
}

// Init Git repository
func initRoleRepo(r *role.Role) error {
	name, dir, url, err := role.ParseOrigin(r.Origin)
	if err != nil {
		return err
	}
	if dir == "" {
		dir = path.Join(r.Target, RolesDir, r.Name)
	}
	repo, err := git.NewRepo(name, dir, url)
	if err != nil {
		return err
	}
	repo.Name = r.Name
	if repo.Path != r.Source {
		r.Source = repo.Path
	}
	r.Repo = repo
	return nil
}

// Init role configuration
func initRoleConfig(r *role.Role) error {

	cfg := viper.New()
	cfg.SetConfigName(configName)
	cfg.AddConfigPath(r.Source)
	err := cfg.ReadInConfig()
	switch err.(type) {
	case nil:
	case viper.ConfigFileNotFoundError:
		logger.Warn(err)
	default:
		// logger.Error(err)
		return err
	}
	cfgUsed := cfg.ConfigFileUsed()
	if cfgUsed != "" {
		logger.Debugln("Using role config file:", cfgUsed)
	}

	// r.Config = cfg

	// logger.Infof(" ROLE --------------------- %+v\n", r)
	// logger.Infof(" ROLE CONFIG --------------------- %+v\n", r.Config)

	/*
	// p := &role.Package{}
	// err = Config.UnmarshalKey("roles.git", &p)
	// if err != nil {
	// 	return err
	// }
	// logger.Infof(" PACKAGE --------------------- %+v\n", p)
	logger.Infof("\n!!!!!!!!! Config %+v\n\n", Config)
	// logger.Infof("\nRole %+v\n\n", r)
	// logger.Infof("= %+v\n", Config.Get("roles"))
	logger.Infof("!!!!!!!! Dot.Roles %+v\n", Dot.Roles)
	*/

	// p := &role.Package{}
	// err = cfg.UnmarshalKey("role", &p)
	// if err != nil {
	// 	return err
	// }

	p := cfg.Sub("role")

	r.Config = p
	// v := reflect.ValueOf(p)
	// logger.Infof("--------------------v %+v\n", v)
	// for i, key := range v.MapKeys() {
	// 	// key := v.Type().Field(i).Name
	// 	val := v.Field(i) // .Elem()
	// 	logger.Printf("\nSET %s: %+v\n\n", key, val)
	// }

	/*
	base := &role.Role{}
	for _, b := range Dot.Roles {
		if b.Name == r.Name {
			base = b
			logger.Infof("BASE ROLE %+v\n", base)
			break
		}
	}

	if r.Config == nil {
		r.Config = &viper.Viper{}
	}

	if base.Config != nil {
		logger.Infof("-------_> BASE %+v\n", base.Config)

		r.Config = base.Config
	}

	// baseRole := &role.Role{}
	// err = baseRole.Config.UnmarshalKey("role", &baseRole)
	// if err != nil {
	// 	return err
	// }

	// values := make(map[string]interface{}, v.NumField())

	//r.Config = reflect.Indirect(reflect.ValueOf(Dot.Roles)).Interface()
	//r.Config = make(map[string]interface{}, len(Dot.Roles))
	//copy(r.Config, Config)


	logger.Infof("BEFORE ===================== %+v\n", base.Config.GetString("template"))
	// r.Config.Set("name", "xxx")

	for i := 0; i < v.NumField(); i++ {
		key := v.Type().Field(i).Name
		val := v.Field(i) // .Elem()

		logger.Printf("\nSET %s: %+v\n\n", key, val)

		fv := v.FieldByName(key)

		fv.Set(val)
		// r.Config.Set(key, val)
		// r.Config[key] = val

		// ptr? := v.Field(i).Interface()
		// values[i] = v.Field(i).Interface()
	}

	logger.Infof("AFTER ===================== %+v\n", base.Config.GetString("template"))
	// r.Config = cfg // .Get("role")
	// for k, v := range v {
	// 	logger.Infoln(k, v)
	// 	// a[k] = v
	// }
	logger.Infoln("r.Config", r.Config)
	*/

	// r.Package = &role.Package{}
	// err = cfg.UnmarshalKey("package", &r.Package)
	// if err != nil {
	// 	return err
	// }
	// logger.Infof(" r --------------------- %+v\n", r)
	// logger.Infof(" PACKAGE --------------------- %+v\n", r.Package)

	return nil
}
