package command

import (
	"fmt"
	"io"
	"os"
	"path/filepath"

	// "reflect"
	"runtime"

	"github.com/LEI/dot/cli/config"
	cliconfig "github.com/LEI/dot/cli/config"
	"github.com/LEI/dot/cli/config/tasks"
	cliflags "github.com/LEI/dot/cli/flags"
	"github.com/LEI/dot/pkg/git"
	"github.com/LEI/dot/pkg/homedir"
	"github.com/LEI/dot/pkg/ostype"
	"github.com/LEI/dot/pkg/prompt"
	"github.com/LEI/dot/system"

	// "github.com/LEI/dot/pkg/sliceutils"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

var (
	// Options ...
	Options *cliflags.Options = &cliflags.Options{}

	homeDir = homedir.Get()

	noConfirmFile = ".dotnc"
)

func init() {
	// osTypes := ostype.Get()
	// fmt.Println("OS TYPES", osTypes)
}

// Streams is an interface which exposes the standard input and output streams
type Streams interface {
	In() io.ReadCloser // *InStream
	Out() io.Writer    // *OutStream
	Err() io.Writer
}

// Cli represents the dot command line client.
type Cli interface {
	In() io.ReadCloser // *InStream
	Out() io.Writer    // *OutStream
	Err() io.Writer
	Config() *config.Config
}

// DotCli is an instance the dot command line client.
// Instances of the client can be returned from NewDotCli.
type DotCli struct {
	config *config.Config
	in     io.ReadCloser // *InStream
	out    io.Writer     // *OutStream
	err    io.Writer
}

// In returns the reader used for stdin
func (cli *DotCli) In() io.ReadCloser { // *InStream {
	return cli.in
}

// Out returns the writer used for stdout
func (cli *DotCli) Out() io.Writer { // *OutStream {
	return cli.out
}

// Err returns the writer used for stderr
func (cli *DotCli) Err() io.Writer {
	return cli.err
}

// ShowHelp shows the command help.
func (cli *DotCli) ShowHelp(cmd *cobra.Command, args []string) error {
	cmd.SetOutput(cli.err)
	cmd.HelpFunc()(cmd, args)
	return nil
}

// AddCommands ...
func (cli *DotCli) AddCommands(cmd *cobra.Command) {
	cmd.AddCommand(
		NewDirCommand(cli),
		NewCopyCommand(cli),
		NewLinkCommand(cli),
		NewTemplateCommand(cli),
		NewLineCommand(cli),
	)
}

// Config returns the configuration
func (cli *DotCli) Config() *config.Config {
	return cli.config
}

// ParseConfig unmarshals the configuration
func (cli *DotCli) ParseConfig(i *interface{}) error {
	return cli.config.Parse(&i)
}

// Initialize the dotCli runs initialization that must happen after command
// line flags are parsed.
func (cli *DotCli) Initialize(opts *cliflags.Options) error {
	cli.config = LoadDefaultConfig(cli.err)
	if configFile := cli.config.FileUsed(); configFile != "" && opts.Verbose > 0 {
		fmt.Fprintf(cli.Out(), "> Using config file: %s\n", configFile)
	}
	// DOT_SOURCE
	if opts.Source == "" { // opts.Source
		opts.Source = homeDir // filepath.Join(homeDir, config.RoleConfigDir)
	}
	// DOT_TARGET
	if opts.Target == "" { // opts.Target
		opts.Target = homeDir
	}

	if opts.Target == "" {
		fmt.Println("fatal: target not specified")
		os.Exit(1)
	}
	if opts.Target == homeDir {
		noconfirm := filepath.Join(homeDir, noConfirmFile)
		msg := fmt.Sprintf("%s: use homedir?", homeDir)
		if !system.Exists(noconfirm) && !prompt.AskConfirmation(msg) {
			fmt.Println("abort homedir")
			os.Exit(1)
		}
	}
	// fmt.Println("SOURCE", opts.Source)
	// fmt.Println("TARGET", opts.Target)

	git.Force = opts.Force
	git.Stdout = cli.Out()
	git.Stderr = cli.Err()

	system.DryRun = opts.DryRun

	// cli.config.Verbosity = len(cli.config.Verbose)
	tasks.Verbose = opts.Verbose
	tasks.Stdout = cli.Out()
	tasks.Stderr = cli.Err()

	// opts.Action

	// err := cli.config.Parse(&config.DotConfig)
	// if err != nil {
	// 	return err
	// }

	// dotConfig := cli.config.Get("roles")
	// fmt.Println(reflect.TypeOf(dotConfig))
	// fmt.Println("=", reflect.TypeOf(cli.config.GetAll()))

	return nil
}

// InitializeAction adds common action falgs on the FlagSet
func (cli *DotCli) InitializeAction(opts *cliflags.Options, flags *pflag.FlagSet) {
	flags.StringVarP(&opts.Source, "source", "s", "", "Source directory")
	flags.StringVarP(&opts.Target, "target", "t", "", "Target directory")
}

// Parse roles
func (cli *DotCli) Parse(opts *cliflags.Options) error {
	var roles []*config.Role
	cliConfig := cli.Config()
	cfgroles := cliConfig.Get("roles")
	if cfgroles == nil {
		// fmt.Fprintf(os.Stderr, "no roles in: %+v\n", cliConfig.GetAll())
		return fmt.Errorf("no roles")
	}
	configRoles := cfgroles.([]interface{})
	for _, r := range configRoles {
		role, err := config.NewRole(r)
		if err != nil {
			return err
		}
		if len(role.OS) > 0 && !ostype.Has(role.OS...) {
			continue
		}
		roles = append(roles, role)
	}
	for _, role := range roles {
		if err := cli.ParseRole(opts, role); err != nil {
			return err
		}
		// Check OS a second time, in case it is specified in role config only
		if len(role.OS) > 0 && !ostype.Has(role.OS...) {
			continue
		}
		if len(role.Deps) > 0 {
		DEPS:
			for _, d := range role.Deps {
				for _, r := range roles {
					if d == r.Name {
						continue DEPS
					}
				}
				return fmt.Errorf(role.Name, "required", role.Deps, "but", d, "is not found")
			}
		}
		// if err := cli.Prepare(role); err != nil {
		// 	return err
		// }
		if len(opts.Filter) > 0 {
			matched := false
			for _, f := range opts.Filter {
				if f == role.Name {
					matched = true
					break
				}
			}
			if !matched {
				continue
			}
		}
		cli.config.Roles = append(cli.config.Roles, role)
	}
	if len(cli.config.Roles) == 0 {
		return fmt.Errorf("no roles (total: %d) matching filter: %+v", len(configRoles), opts.Filter)
	}
	return nil
}

// Roles ...
func (cli *DotCli) Roles() []*config.Role {
	return cli.config.Roles
}

// ParseRole ...
func (cli *DotCli) ParseRole(opts *cliflags.Options, role *config.Role) error {
	if role.Path == "" {
		role.Path = filepath.Join(opts.Source, config.RoleConfigDir, role.Name)
		// role.Path = filepath.Join(opts.Source, role.Name)
	}
	role.Ignore = opts.Ignore
	role.Target = opts.Target
	if err := cli.config.LoadRole(role); err != nil {
		if opts.Verbose > 0 {
			fmt.Fprintf(os.Stderr, "WARNING: Error loading role config file: %v\n", err)
		}
		return nil // err
	}
	// TODO init env
	os.Setenv("OS", runtime.GOOS)
	if err := role.Prepare(); err != nil {
		return err
	}
	return nil
}

// CheckRole ...
func (cli *DotCli) CheckRole(action string, role *config.Role) error {
	if err := tasks.Check(role.Links); err != nil {
		return err
	}
	return nil
}

// ExecRole ...
func (cli *DotCli) ExecRole(action string, role *config.Role) error {
	if len(role.OS) > 0 && !ostype.Has(role.OS...) {
		return fmt.Errorf("bad os: %s", role.OS)
	}
	// if len(role.If) > 0 {
	// 	return nil
	// }
	if tasks.Verbose > 0 {
		var s string
		switch action {
		case "install":
			s = "Installing"
		case "remove":
			s = "Removing"
			// default:
			// 	s = strings.ToTitle(s)
		}
		fmt.Printf("## %s %s...\n", s, role.Name)
	}
	tasks.ExecDir = role.Path
	// if err := tasks.Check(role.Install); err != nil {
	// 	return err
	// }
	// if err := tasks.Check(role.Remove); err != nil {
	// 	return err
	// }
	// if err := tasks.Check(role.Packages); err != nil {
	// 	return err
	// }
	if err := tasks.Check(role.Dirs); err != nil {
		return err
	}
	if err := tasks.Check(role.Files); err != nil {
		return err
	}
	if err := tasks.Check(role.Links); err != nil {
		return err
	}
	if err := tasks.Check(role.Templates); err != nil {
		return err
	}
	if err := tasks.Check(role.Lines); err != nil {
		return err
	}
	// if err := tasks.Check(role.PostInstall); err != nil {
	// 	return err
	// }
	// if err := tasks.Check(role.PostRemove); err != nil {
	// 	return err
	// }
	switch action {
	case "install":
		if err := tasks.Install(role.Install); err != nil {
			return err
		}
		if err := tasks.Install(role.Packages); err != nil {
			return err
		}
		if err := tasks.Install(role.Dirs); err != nil {
			return err
		}
		if err := tasks.Install(role.Files); err != nil {
			return err
		}
		if err := tasks.Install(role.Links); err != nil {
			return err
		}
		if err := tasks.Install(role.Templates); err != nil {
			return err
		}
		if err := tasks.Install(role.Lines); err != nil {
			return err
		}
		if err := tasks.Install(role.PostInstall); err != nil {
			return err
		}
	case "remove":
		if err := tasks.Remove(role.Remove); err != nil {
			return err
		}
		if err := tasks.Remove(role.Files); err != nil {
			return err
		}
		if err := tasks.Remove(role.Links); err != nil {
			return err
		}
		if err := tasks.Remove(role.Dirs); err != nil {
			return err
		}
		if err := tasks.Remove(role.Templates); err != nil {
			return err
		}
		if err := tasks.Remove(role.Lines); err != nil {
			return err
		}
		if err := tasks.Remove(role.Packages); err != nil {
			return err
		}
		if err := tasks.Remove(role.PostRemove); err != nil {
			return err
		}
	default:
		return fmt.Errorf("%s: not implemented", action)
	}
	tasks.ExecDir = ""
	return nil
}

// NewDotCli returns a DotCli instance with IO output and error streams set by in, out and err.
func NewDotCli(in io.ReadCloser, out, err io.Writer) *DotCli {
	// in: NewInStream(in), out: NewOutStream(out), err: err
	return &DotCli{in: in, out: out, err: err}
}

// LoadDefaultConfig attempts to load the default config file and returns
// an initialized Config struct if none is found.
func LoadDefaultConfig(stderr io.Writer) *config.Config {
	config, err := cliconfig.Load(cliconfig.Dir())
	if err != nil {
		fmt.Fprintf(stderr, "WARNING: Error loading config file: %v\n", err)
	}
	// if !config.ContainsAuth() {
	// 	credentials.DetectDefaultStore(config)
	// }
	return config
}

type execOptions struct {
	action string
}

func runExec(dotCli *DotCli, opts execOptions) error {
	roles := dotCli.Roles()
	if len(roles) == 0 {
		return fmt.Errorf("no roles to exec")
	}
	for _, r := range roles {
		switch opts.action {
		case "install":
			if err := tasks.Install(r.Install); err != nil {
				return err
			}
		case "remove":
			if err := tasks.Remove(r.Remove); err != nil {
				return err
			}
		default:
			return fmt.Errorf("%s: not implemented", opts.action)
		}
	}
	return nil
}

func runPostExec(dotCli *DotCli, opts execOptions) error {
	roles := dotCli.Roles()
	if len(roles) == 0 {
		return fmt.Errorf("no roles to exec")
	}
	for _, r := range roles {
		switch opts.action {
		case "install":
			if err := tasks.Install(r.PostInstall); err != nil {
				return err
			}
		case "remove":
			if err := tasks.Remove(r.PostRemove); err != nil {
				return err
			}
		default:
			return fmt.Errorf("%s: not implemented", opts.action)
		}
	}
	return nil
}
