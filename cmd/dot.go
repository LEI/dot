package cmd

import (
	"fmt"
	// "io"
	// "io/ioutil"
	"os"
	"reflect"

	"github.com/jessevdk/go-flags"

	"github.com/LEI/dot/dotfile"
)

// DotCmd ...
type DotCmd struct {
	BaseCmd

	List    ListCmd    `command:"list" alias:"ls" description:"List"`
	Install InstallCmd `command:"install" subcommands-optional:"true" alias:"i" description:"Install"`
	Remove  RemoveCmd  `command:"remove" subcommands-optional:"true" alias:"rm" description:"Remove"`

	Source flags.Filename `short:"s" long:"source" description:"Path to source file"`
	Target flags.Filename `short:"t" long:"target" description:"Path to target link" default:"$HOME"`

	RoleDir string `long:"role-dir" description:"Role directory name relative to source" default:".dot"`

	// env:"DOT_CONFIG" default:".dot"
	Config    func(s string) error `short:"c" long:"config" description:"Config file name" default:".dot.yml"`
	IniConfig func(s string) error `hidden:"true" short:"i" long:"ini-config" description:"Path to INI config file" no-ini:"true"`

	RoleFilter      []string `short:"n" long:"name" description:"Filter roles by name"`
	IgnoreDeps      bool     `short:"I" long:"ignore-dependencies" description:"Honor the role filter"`
	RemoveEmptyDirs bool     `short:"R" long:"remove-empty-dirs" description:"Remove modified directories if empty"`

	Packages bool `short:"P" long:"packages" description:"Install or remove required packages"`
	Sudo     bool `short:"S" long:"sudo" description:"Use sudo to execute package manager"`

	// Do not error out if unable to git clone or pull
	NoSync bool `short:"N" long:"no-sync" description:"Skip network operations"`

	// Ignore uncommitted changes in repository
	Force bool `short:"f" long:"force" description:"Ignore uncommitted changes"`

	// Debug bool `short:"D" long:"debug" description:""`
	DryRun bool `short:"d" long:"dry-run" description:"Do not execute tasks"`

	// Slice of bool will append 'true' each time the option
	// is encountered (can be set multiple times, like -vvv)
	Verbose []bool `short:"v" long:"verbose" description:"Show verbose debug information"`

	Version bool `short:"V" long:"version" description:"Print the version and exit"`
}

// Parse ...
func Parse() ([]string, error) {
	remaining, err := parser.Parse()
	if err != nil {
		if flagsErr, ok := err.(*flags.Error); ok {
			switch flagsErr.Type {
			case flags.ErrHelp:
				parser.WriteHelp(os.Stdout)
				os.Exit(0)
				// break
			case flags.ErrCommandRequired:
				// err = Options.Install.Execute(remaining)
				// remaining = []string{}
				err = nil
				break
			default:
				fmt.Fprintf(os.Stderr, "Error parsing args: %s\n", err)
				parser.WriteHelp(os.Stdout)
				os.Exit(1)
			}
		}
	}
	// Update variables
	source = dotfile.ExpandEnv(string(Options.Source))
	target = dotfile.ExpandEnv(string(Options.Target))
	dotfile.DryRun = Options.DryRun
	Verbose = len(Options.Verbose)
	dotfile.Verbose = Verbose
	dotfile.RemoveEmptyDirs = Options.RemoveEmptyDirs
	// WriteIniConfig(parser)
	return remaining, err
}

// CommandHandler ...
// func (cmd *DotCmd) CommandHandler(command flags.Commander, args []string) error {
// 	return command.Execute(args)
// }

// Execute default action (list roles)
func (cmd *DotCmd) Execute(args []string) error {
	return cmd.List.Execute(args)
	// return cmd.Install.Execute(args)
}

var (
	// Options ...
	Options DotCmd

	// ConfigName ...
	ConfigName string

	// Action (install/remove)
	Action string

	// RunOnly task filter
	RunOnly []string

	// Verbose level (verbosity)
	Verbose int
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
		// Note: cannot fallthrough in type switch
		switch cmd.(type) {
		case nil:
			cmd = &ListCmd{}
			// Options.Force = true
			// Options.NoSync = true
			RunOnly = append(RunOnly, "list")
			break
		case *DotCmd:
			RunOnly = append(RunOnly, "list")
			break
		case *ListCmd:
			RunOnly = append(RunOnly, "list")
			break
		// case nil:
		// 	cmd = &DotCmd{}
		// 	// cmd = &InstallCmd{}
		// 	fallthrough
		// case *DotCmd:
		// 	fallthrough
		case *InstallCmd:
			// Default action: install all
			break
		case *RemoveCmd:
			Action = "remove"
			// Run all
			break
		case *CopyCmd:
			RunOnly = append(RunOnly, "copy")
			break
		// case *Exec:
		// 	RunOnly = append(RunOnly, "exec")
		// 	break
		case *LineCmd:
			RunOnly = append(RunOnly, "line")
			break
		case *LinkCmd:
			RunOnly = append(RunOnly, "link")
			break
		case *TemplateCmd:
			RunOnly = append(RunOnly, "template")
			break
		case *PackageCmd:
			RunOnly = append(RunOnly, "package")
			break
		default:
			return fmt.Errorf("# Unhandled command (%+v): %+v", reflect.TypeOf(cmd).Elem(), cmd)
		}
		// return fmt.Errorf("Execute %+v: %+v\n %+v", reflect.TypeOf(cmd).Elem(), RunOnly, cmd)
		return cmd.Execute(args)
	}
}

// func readConfig(config *dot.Config) func(s string) error {
// }

func readIniConfig(parser *flags.Parser) func(s string) error {
	return func(s string) error {
		ini := flags.NewIniParser(parser)
		// ini.ParseAsDefaults = true
		return ini.ParseFile(s)
	}
}

// WriteIniConfig ...
func WriteIniConfig(parser *flags.Parser) {
	ini := flags.NewIniParser(parser)
	// flags.IniDefault
	ini.Write(os.Stdout, flags.IniIncludeDefaults|flags.IniCommentDefaults|flags.IniIncludeComments)
}
