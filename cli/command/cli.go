package command

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	// "reflect"
	"runtime"
	"strings"

	"github.com/LEI/dot/cli/config"
	"github.com/LEI/dot/cli/config/tasks"
	cliconfig "github.com/LEI/dot/cli/config"
	cliflags "github.com/LEI/dot/cli/flags"
	"github.com/spf13/cobra"
)

var (
	homeDir = os.Getenv("HOME")
	roleDir = ".dot"
)

// Streams is an interface which exposes the standard input and output streams
type Streams interface {
	In() io.ReadCloser // *InStream
	Out() io.Writer // *OutStream
	Err() io.Writer
}

// Cli represents the dot command line client.
type Cli interface {
	In() io.ReadCloser // *InStream
	Out() io.Writer // *OutStream
	Err() io.Writer
	Config() *config.Config
}

// DotCli is an instance the dot command line client.
// Instances of the client can be returned from NewDotCli.
type DotCli struct {
	config *config.Config
	in     io.ReadCloser // *InStream
	out    io.Writer // *OutStream
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
	// DOT_SOURCE
	if cli.config.Source == "" {
		cli.config.Source = filepath.Join(homeDir, roleDir)
	}
	// DOT_TARGET
	if cli.config.Target == "" {
		cli.config.Target = "/tmp/todo" // homeDir
	}

	// opts

	// err := cli.config.Parse(&config.DotConfig)
	// if err != nil {
	// 	return err
	// }

	// dotConfig := cli.config.Get("roles")
	// fmt.Println(reflect.TypeOf(dotConfig))
	// fmt.Println("=", reflect.TypeOf(cli.config.GetAll()))

	return nil
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
			fmt.Println("TODO FILTER ROLE:", filter)
			// matched := false
			// for _, f := range filter {
			// 	if f == role.Name {
			// 		matched = true
			// 		break
			// 	}
			// }
			// if !matched {
			// 	continue
			// }
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
	if err := cli.config.Load(role); err != nil {
		fmt.Fprintf(os.Stderr, "WARNING: Error loading role config file: %v\n", err)
	}
	// Prepare tasks paths
	links := &tasks.Links{}
	for _, l := range *role.Link {
		os.Setenv("OS", runtime.GOOS)
		src := os.ExpandEnv(l.Source)
		dst := os.ExpandEnv(l.Target)
		if !filepath.IsAbs(src) {
			src = filepath.Join(role.Dir, src)
		}
		*links = append(*links, l)
		if strings.Contains(src, "*") {
			// fmt.Println("*", src, dst)
			glob, err := filepath.Glob(src)
			if err != nil {
				return role, err
			}
		// GLOB:
			for _, s := range glob {
				// // Extract source file name
				// _, n := filepath.Split(s)
				// for _, i := range ignore {
				// 	// Check for ignored patterns
				// 	matched, err := filepath.Match(i, n)
				// 	if err != nil {
				// 		return err
				// 	}
				// 	if matched {
				// 		continue GLOB
				// 	}
				// }
				// fmt.Println("PREPARE GLOB", s, "/", dst)
				t, err := prepareTarget(cli.config.Target, s, dst)
				if err != nil {
					return role, err
				}
				l.Source = s
				l.Target = t
				*links = append(*links, l)
			}
		} else {
			t, err := prepareTarget(cli.config.Target, src, dst)
			if err != nil {
				return role, err
			}
			l.Source = src
			l.Target = t
			*links = append(*links, l)
		}
	}
	*role.Link = *links
	return role, nil
}

func prepareTarget(dir, src, dst string) (string, error) {
	//fmt.Println("+", src, dst)
	_, f := filepath.Split(src)
	if f == "" {
		return "", fmt.Errorf("error (no source file name) while parsing: %s / %s", src, dst)
	}
	if !filepath.IsAbs(dst) {
		dst = filepath.Join(dir, dst)
	}
	// if _, err := dotfile.CreateDir(baseDir); err != nil {
	// 	return baseDir, err
	// }
	t := filepath.Join(dst, f)
	return t, nil
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
	configFile := config.FileUsed()
	if configFile != "" { // debug
		fmt.Printf("Using config file: %s\n", configFile)
	}
	// if !config.ContainsAuth() {
	// 	credentials.DetectDefaultStore(config)
	// }
	return config
}
