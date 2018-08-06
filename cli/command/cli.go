package command

import (
	"fmt"
	"io"
	// "os"
	// "path/filepath"
	// "runtime"

	"github.com/LEI/dot/cli/config"
	cliconfig "github.com/LEI/dot/cli/config"
	cliflags "github.com/LEI/dot/cli/flags"
	"github.com/spf13/cobra"
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
	config     *config.Config
	in             io.ReadCloser // *InStream
	out            io.Writer // *OutStream
	err            io.Writer
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

// Initialize the dotCli runs initialization that must happen after command
// line flags are parsed.
func (cli *DotCli) Initialize(opts *cliflags.Options) error {
	cli.config = LoadDefaultConfig(cli.err)

	return nil
}

// NewDotCli returns a DotCli instance with IO output and error streams set by in, out and err.
func NewDotCli(in io.ReadCloser, out, err io.Writer) *DotCli {
	// in: NewInStream(in), out: NewOutStream(out), err: err
	return &DotCli{in: in, out: out, err: err}
}

// LoadDefaultConfig attempts to load the default config file and returns
// an initialized Config struct if none is found.
func LoadDefaultConfig(err io.Writer) *config.Config {
	config, e := cliconfig.Load(cliconfig.Dir())
	if e != nil {
		fmt.Fprintf(err, "WARNING: Error loading config file:%v\n", e)
	}
	// if !config.ContainsAuth() {
	// 	credentials.DetectDefaultStore(config)
	// }
	return config
}
