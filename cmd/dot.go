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
	// Options ...
	Options DotCmd

	// ConfigName ...
	ConfigName string

	// Action (install/remove)
	Action string
	// Only (copy, link, template)
	// TODO? --filter=*
	Only []string

	// Verbose mode
	Verbose bool
)

var executedCommand flags.Commander

func init() {
	// args, err := flags.Parse(&options)

	// if err != nil {
	// 	panic(err)
	// }

	// fmt.Printf("Args: %v\n", args)

	parser.SubcommandsOptional = true

	// Options.Config = readConfig(config)
	Options.IniConfig = readIniConfig(parser)

	parser.CommandHandler = func(cmd flags.Commander, args []string) error {
		// executedCommand = cmd
		// return Options.CommandHandler(cmd, args)
		// fmt.Printf("----------> %+v\n", cmd)
		Action = "install"
		switch cmd.(type) {
		case nil:
			cmd = &DotCmd{}
			// cmd = &InstallCmd{}
		case *DotCmd:
		case *InstallCmd:
			break
		case *RemoveCmd:
			Action = "remove"
			break
		case *CopyCmd:
			Only = append(Only, "copy")
			break
		// case *Exec:
		// 	Only = append(Only, "exec")
		// 	break
		case *LineCmd:
			Only = append(Only, "line")
			break
		case *LinkCmd:
			Only = append(Only, "link")
			break
		case *TemplateCmd:
			Only = append(Only, "template")
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

	Source flags.Filename `short:"s" long:"source" description:"Path to source file"`
	Target flags.Filename `short:"t" long:"target" description:"Path to target link" default:"$HOME"`

	// env:"DOT_CONFIG" default:".dot"
	Config    func(s string) error `short:"c" long:"config" description:"Global config file name" default:".dot.yml"`
	IniConfig func(s string) error `short:"i" long:"ini" description:"INI config file" no-ini:"true"`

	RoleDir    string   `short:"R" long:"role-dir" description:"Role directory name" default:".dot"`
	RoleFilter []string `short:"r" long:"role" description:"Filter roles by name"`

	// Debug bool `short:"D" long:"debug" description:""`

	// Ignore uncommitted changes in repository
	Force bool `short:"f" long:"force" description:"Ignore uncommitted changes"`

	// Do not error out if unable to git clone or pull
	NoSync bool `short:"n" long:"no-sync" description:"Skip network operations"`

	DryRun bool `short:"d" long:"dry-run" description:"Do not execute tasks"`

	// Slice of bool will append 'true' each time the option
	// is encountered (can be set multiple times, like -vvv)
	Verbose []bool `short:"v" long:"verbose" description:"Show verbose debug information"`

	Version bool `short:"V" long:"version" description:"Print the version and exit"`

	Install InstallCmd `command:"install" subcommands-optional:"true" alias:"i" description:"Install"`
	Remove  RemoveCmd  `command:"remove" subcommands-optional:"true" alias:"r" description:"Remove"`
}

// CommandHandler ...
// func (cmd *DotCmd) CommandHandler(command flags.Commander, args []string) error {
// 	return command.Execute(args)
// }

// Execute ...
func (cmd *DotCmd) Execute(args []string) error {
	// fmt.Println("exec dot cmd", args)
	return cmd.Install.Execute(args)
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
