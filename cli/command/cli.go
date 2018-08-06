package command

import (
	"fmt"
	"io"
	// "os"
	// "path/filepath"
	// "runtime"

	cliconfig "github.com/LEI/dot/cli/config"
	cliflags "github.com/LEI/dot/cli/flags"
	"github.com/LEI/dot/cli/config/configfile"
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
	ConfigFile() *configfile.ConfigFile
}

// DotCli is an instance the dot command line client.
// Instances of the client can be returned from NewDotCli.
type DotCli struct {
	configFile     *configfile.ConfigFile
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

// // ConfigFile returns the ConfigFile
// func (cli *DotCli) ConfigFile() *configfile.ConfigFile {
// 	return cli.configFile
// }

// Initialize the dotCli runs initialization that must happen after command
// line flags are parsed.
func (cli *DotCli) Initialize(opts *cliflags.Options) error {
	cli.configFile = LoadDefaultConfigFile(cli.err)

	return nil
}

// NewDotCli returns a DotCli instance with IO output and error streams set by in, out and err.
func NewDotCli(in io.ReadCloser, out, err io.Writer) *DotCli {
	// in: NewInStream(in), out: NewOutStream(out), err: err
	return &DotCli{in: in, out: out, err: err}
}

// LoadDefaultConfigFile attempts to load the default config file and returns
// an initialized ConfigFile struct if none is found.
func LoadDefaultConfigFile(err io.Writer) *configfile.ConfigFile {
	configFile, e := cliconfig.Load(cliconfig.Dir())
	if e != nil {
		fmt.Fprintf(err, "WARNING: Error loading config file:%v\n", e)
	}
	// if !configFile.ContainsAuth() {
	// 	credentials.DetectDefaultStore(configFile)
	// }
	return configFile
}
