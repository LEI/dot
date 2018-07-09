package cmd

import (
	"fmt"
	// "io"
	"os"

	"github.com/LEI/dot/cfg"
	"github.com/LEI/dot/dot"

	"github.com/jessevdk/go-flags"
)

// GlobalConfig ...
var GlobalConfig *dot.Config

// DotCmd ...
type DotCmd struct {
	BaseCmd

	// Slice of bool will append 'true' each time the option
	// is encountered (can be set multiple times, like -vvv)
	Verbose []bool `short:"v" long:"verbose" description:"Show verbose debug information"`

	Version bool `short:"V" long:"version" description:"Print the version and exit"`

	// env:"DOT_CONFIG" default:".dot"
	Config func(s string) error `short:"c" long:"config" description:"Config file name"`
	IniConfig func(s string) error `short:"i" long:"ini-config" description:"INI config file" no-ini:"true"`

	// Debug bool `short:"d" long:"debug" description:""`

	Install InstallCmd `command:"install" subcommands-optional:"true" alias:"i" description:"Install"`
	Remove  RemoveCmd  `command:"remove" subcommands-optional:"true" alias:"r" description:"Remove"`

	// Role string `short:"r" long:"role" description:""`
	// URL string `short:"u" long:"url" description:""`
	// Roles []dot.Role `short:"r" long:"roles" description:""`
	// Roles map[string]string `short:"r" long:"roles" description:""`
}

// Cmd ...
// var Cmd flags.Command

// Options ...
var Options DotCmd

var parser = flags.NewParser(&Options, flags.Default)

// Execute ...
func (cmd *DotCmd) Execute(args []string) error {
	fmt.Println("exec cmd", args)

	// fmt.Println(Dot.Commands)
	// if c, ok := cmd.(*flags.Command); ok {
	// 	fmt.Println(c)
	// }

	return nil // Install.Execute(args)
}

func init() {
	// args, err := flags.Parse(&options)

	// if err != nil {
	// 	panic(err)
	// }

	// fmt.Printf("Args: %v\n", args)

	parser.SubcommandsOptional = false

	Options.Config = readConfig(GlobalConfig)
	Options.IniConfig = readIniConfig(parser)
}

// Parse ...
func Parse() ([]string, error) {
	// TODO: control (mute) output?
	remaining, err := parser.Parse()
	if err != nil {
		if flagsErr, ok := err.(*flags.Error); ok {
			switch flagsErr.Type {
			case flags.ErrHelp:
				os.Exit(1)
			case flags.ErrCommandRequired:
				// FIXME: DotCmd.Execute() never called
				// when first-level sub commands are optional
				err = Options.Install.Execute(remaining)
				remaining = []string{}
				// os.Exit(1)
			// default:
			// 	fmt.Println("Error parsing args:", err)
			// 	os.Exit(1)
			}
		}
	}

	return remaining, err
}

// WriteHelp ...
// func WriteHelp(o io.Writer) {
// 	parser.WriteHelp(o)
// }

// Help ...
func Help(rc int) {
	parser.WriteHelp(os.Stdout)
	os.Exit(rc)
}

// GetParser ...
func GetParser() *flags.Parser {
	return parser
}

func readConfig(config *dot.Config) func(s string) error {
	return func(s string) error {
		configFile, err := cfg.Load(&config, s)
		if err != nil {
			return err
		}

		if configFile != "" {
			fmt.Println("Using configuration file:", configFile)
		}

		return nil
	}
}

func readIniConfig(parser *flags.Parser) func(s string) error {
	return func(s string) error {
		ini := flags.NewIniParser(parser)
		// i.ParseAsDefaults = true
		return ini.ParseFile(s)
	}
}

// WriteIniConfig ...
func WriteIniConfig(parser *flags.Parser) {
	ini := flags.NewIniParser(parser)
	// flags.IniDefault
	ini.Write(os.Stdout, flags.IniIncludeDefaults|flags.IniCommentDefaults|flags.IniIncludeComments)
}
