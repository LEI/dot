package cmd

import (
	"fmt"
	// "io"
	// "io/ioutil"
	"os"
	"reflect"

	"github.com/jessevdk/go-flags"
)

var (
	// ConfigName ...
	ConfigName string

	// GlobalConfig ...
	GlobalConfig *Config

	// Options ...
	Options DotCmd

	// Source ....
	Source string

	// Target ....
	Target string

	// Verbose ....
	Verbose bool
)

var parser = flags.NewParser(&Options, flags.Default)

var executedCommand flags.Commander

func init() {
	// args, err := flags.Parse(&options)

	// if err != nil {
	// 	panic(err)
	// }

	// fmt.Printf("Args: %v\n", args)

	parser.SubcommandsOptional = false

	// Options.Config = readConfig(GlobalConfig)
	Options.IniConfig = readIniConfig(parser)

	parser.CommandHandler = func(cmd flags.Commander, args []string) error {
		// executedCommand = cmd
		// return Options.CommandHandler(cmd, args)
		// fmt.Printf("----------> %+v\n", cmd)
		switch cmd.(type) {
		case *DotCmd:
		case *InstallCmd:
		case *RemoveCmd:
			break
		case *CopyCmd:
			Options.ActionFilter = append(Options.ActionFilter, "Copy")
			break
		case *LinkCmd:
			Options.ActionFilter = append(Options.ActionFilter, "Link")
			break
		case *TemplateCmd:
			Options.ActionFilter = append(Options.ActionFilter, "Template")
			break
		default:
			return fmt.Errorf("exec cmd (%+v) %+v", reflect.TypeOf(cmd).Elem(), cmd)
		}
		return cmd.Execute(args)
	}
}

// DotCmd ...
type DotCmd struct {
	BaseCmd

	// Slice of bool will append 'true' each time the option
	// is encountered (can be set multiple times, like -vvv)
	Verbose []bool `short:"v" long:"verbose" description:"Show verbose debug information"`

	Version bool `short:"V" long:"version" description:"Print the version and exit"`

	// env:"DOT_CONFIG" default:".dot"
	Config     func(s string) error `short:"c" long:"config" description:"Global config file name"`
	IniConfig  func(s string) error `short:"i" long:"ini-config" description:"INI config file" no-ini:"true"`
	ConfigName string               `short:"C" long:"config-name" description:"Config file name for roles"`

	// Debug bool `short:"D" long:"debug" description:""`

	// Ignore uncommitted changes in repository
	NoCheck bool `short:"N" long:"no-check" description:"Ignore uncommitted changes"`

	// Do not error out if unable to git clone or pull
	NoSync bool `short:"S" long:"no-sync" description:"Skip network operations"`

	Install InstallCmd `command:"install" subcommands-optional:"true" alias:"i" description:"Install"`
	Remove  RemoveCmd  `command:"remove" subcommands-optional:"true" alias:"r" description:"Remove"`

	// Role string `short:"r" long:"role" description:""`
	// URL string `short:"u" long:"url" description:""`
	// Roles []dot.Role `short:"r" long:"roles" description:""`
	// Roles map[string]string `short:"r" long:"roles" description:""`

	Source flags.Filename `short:"s" long:"source" description:"Path to source file"`
	Target flags.Filename `short:"t" long:"target" description:"Path to target link"`

	ActionFilter []string `short:"a" long:"action" description:"Filter commands by name"`
	RoleFilter   []string `short:"r" long:"role" description:"Filter roles by name"`
}

// CommandHandler ...
// func (cmd *DotCmd) CommandHandler(command flags.Commander, args []string) error {
// 	return command.Execute(args)
// }

// Execute ...
func (cmd *DotCmd) Execute(args []string) error {
	fmt.Println("exec cmd", args)
	// fmt.Println(Dot.Commands)
	// if c, ok := cmd.(*flags.Command); ok {
	// 	fmt.Println(c)
	// }
	return nil // Install.Execute(args)
}

// func readConfig(config *dot.Config) func(s string) error {
// }

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
