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
	"github.com/LEI/dot/pkg/ostype"
	"github.com/LEI/dot/system"
	// "github.com/LEI/dot/pkg/sliceutils"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

var (
	// Options ...
	Options *cliflags.Options = &cliflags.Options{}

	homeDir = os.Getenv("HOME")
	roleDir = ".dot"
)

// func init() {
// 	osTypes := ostype.Get()
// 	fmt.Println("OS TYPES", osTypes)
// }

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
	if cli.config.Source == "" { // opts.Source
		cli.config.Source = filepath.Join(homeDir, roleDir)
	}
	// DOT_TARGET
	if cli.config.Target == "" { // opts.Target
		cli.config.Target = opts.Target
	}

	if cli.config.Target == "" {
		fmt.Println("owait empty target")
		os.Exit(1)
	}
	if cli.config.Target == homeDir {
		fmt.Println("owait", cli.config)
		os.Exit(2)
	}

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
func (cli *DotCli) Parse(filter ...string) error {
	var roles []*config.Role
	configRoles := cli.Config().Get("roles").([]interface{})
	for _, r := range configRoles {
		role, err := cli.ParseRole(r)
		if err != nil {
			return err
		}
		// if err := cli.Prepare(r); err != nil {
		// 	return err
		// }
		if len(filter) > 0 {
			matched := false
			for _, f := range filter {
				if f == role.Name {
					matched = true
					break
				}
			}
			if !matched {
				continue
			}
		}
		roleOS := role.OS
		if len(roleOS) > 0 && !ostype.Has(roleOS...) {
			// fmt.Fprintf(cli.Out(), "Skip role %s (OS: %s)\n", role.Name, roleOS)
			continue
		}
		roles = append(roles, role)
	}
	if len(roles) == 0 {
		return fmt.Errorf("no roles (total: %d) matching filter: %+v", len(configRoles), filter)
	}
	cli.config.Roles = roles
	return nil
}

// Roles ...
func (cli *DotCli) Roles() []*config.Role {
	return cli.config.Roles
}

// ParseRole ...
func (cli *DotCli) ParseRole(i interface{}) (*config.Role, error) {
	role, err := config.NewRole(i)
	if err != nil {
		return role, err
	}
	if role.Path == "" {
		role.Path = filepath.Join(cli.config.Target, ".dot", role.Name)
	}
	if err := cli.config.LoadRole(role); err != nil {
		fmt.Fprintf(os.Stderr, "WARNING: Error loading role config file: %v\n", err)
	}
	// TODO init
	os.Setenv("OS", runtime.GOOS)
	if err := role.Prepare(cli.config.Target); err != nil {
		return role, err
	}
	return role, nil
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
