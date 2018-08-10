package main

import (
	// "errors"
	"fmt"
	"os"
	"path/filepath"
	// "strings"

	"github.com/LEI/dot/cli"
	"github.com/LEI/dot/cli/command"
	cliconfig "github.com/LEI/dot/cli/config"
	cliflags "github.com/LEI/dot/cli/flags"
	"github.com/LEI/dot/system"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

// Version number
const Version = "0.0.0"

var (
	binary string // "dot"

	opts = command.Options
)

func init() {
	_, binary = filepath.Split(os.Args[0])
}

func newDotCommand(dotCli *command.DotCli) *cobra.Command {
	// opts := &cliflags.Options{}
	var flags *pflag.FlagSet

	cmd := &cobra.Command{
		Version:          Version,
		Use:              os.Args[0] + " [OPTIONS] COMMAND [ARG...]",
		Short:            "Dotfiles and system manager",
		SilenceUsage:     true,
		SilenceErrors:    true,
		TraverseChildren: true,
		Args:             noArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			// if opts.Version {
			// 	fmt.Fprintf(dotCli.Out(), version())
			// 	return nil
			// }
			return dotCli.ShowHelp(cmd, args)
		},
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			if cmd.Name() != "list" {
				// Create cache store directory
				if err := system.Init(); err != nil {
					return err
				}
			}
			dotPreRun(opts)
			if err := dotCli.Initialize(opts); err != nil {
				return err
			}
			if err := dotCli.Parse(opts.RoleFilter...); err != nil {
				return err
			}
			return nil // isSupported(cmd, dotCli)
		},
	}

	cli.SetupRootCommand(cmd)

	opts.InstallPersistentFlags(cmd.PersistentFlags())

	flags = cmd.Flags()
	opts.InstallFlags(flags)

	// setFlagErrorFunc(dotCli, cmd, flags, opts)

	// setHelpFunc(dotCli, cmd, flags, opts)

	cmd.SetOutput(dotCli.Out())

	cmd.AddCommand(
		command.NewListCommand(dotCli),
		command.NewSyncCommand(dotCli),
		command.NewInstallCommand(dotCli),
		command.NewRemoveCommand(dotCli),
	)

	// setValidateArgs(dotCli, cmd, flags, opts)

	return cmd
}

// func setFlagErrorFunc(dockerCli *command.DockerCli, cmd *cobra.Command, flags *pflag.FlagSet, opts *cliflags.Options) {
// 	// When invoking `docker stack --nonsense`, we need to make sure FlagErrorFunc return appropriate
// 	// output if the feature is not supported.
// 	// As above cli.SetupRootCommand(cmd) have already setup the FlagErrorFunc, we will add a pre-check before the FlagErrorFunc
// 	// is called.
// 	flagErrorFunc := cmd.FlagErrorFunc()
// 	cmd.SetFlagErrorFunc(func(cmd *cobra.Command, err error) error {
// 		initializeDockerCli(dockerCli, flags, opts)
// 		if err := isSupported(cmd, dockerCli); err != nil {
// 			return err
// 		}
// 		return flagErrorFunc(cmd, err)
// 	})
// }

// func setHelpFunc(dockerCli *command.DockerCli, cmd *cobra.Command, flags *pflag.FlagSet, opts *cliflags.Options) {
// 	cmd.SetHelpFunc(func(ccmd *cobra.Command, args []string) {
// 		initializeDockerCli(dockerCli, flags, opts)
// 		if err := isSupported(ccmd, dockerCli); err != nil {
// 			ccmd.Println(err)
// 			return
// 		}

// 		hideUnsupportedFeatures(ccmd, dockerCli)

// 		if err := ccmd.Help(); err != nil {
// 			ccmd.Println(err)
// 		}
// 	})
// }

// func setValidateArgs(dockerCli *command.DockerCli, cmd *cobra.Command, flags *pflag.FlagSet, opts *cliflags.Options) {
// 	// The Args is handled by ValidateArgs in cobra, which does not allows a pre-hook.
// 	// As a result, here we replace the existing Args validation func to a wrapper,
// 	// where the wrapper will check to see if the feature is supported or not.
// 	// The Args validation error will only be returned if the feature is supported.
// 	visitAll(cmd, func(ccmd *cobra.Command) {
// 		// if there is no tags for a command or any of its parent,
// 		// there is no need to wrap the Args validation.
// 		if !hasTags(ccmd) {
// 			return
// 		}

// 		if ccmd.Args == nil {
// 			return
// 		}

// 		cmdArgs := ccmd.Args
// 		ccmd.Args = func(cmd *cobra.Command, args []string) error {
// 			initializeDockerCli(dockerCli, flags, opts)
// 			if err := isSupported(cmd, dockerCli); err != nil {
// 				return err
// 			}
// 			return cmdArgs(cmd, args)
// 		}
// 	})
// }

// func initializeDockerCli(dockerCli *command.DockerCli, flags *pflag.FlagSet, opts *cliflags.Options) {
// 	if dockerCli.Client() == nil { // when using --help, PersistentPreRun is not called, so initialization is needed.
// 		// flags must be the top-level command flags, not cmd.Flags()
// 		opts.Common.SetDefaultOptions(flags)
// 		dotPreRun(opts)
// 		dockerCli.Initialize(opts)
// 	}
// }

// // visitAll will traverse all commands from the root.
// // This is different from the VisitAll of cobra.Command where only parents
// // are checked.
// func visitAll(root *cobra.Command, fn func(*cobra.Command)) {
// 	for _, cmd := range root.Commands() {
// 		visitAll(cmd, fn)
// 	}
// 	fn(root)
// }

func noArgs(cmd *cobra.Command, args []string) error {
	if len(args) == 0 {
		return nil
	}
	return fmt.Errorf(
		"%s: '%s' is not a %s command.\nSee '%s --help'",
		binary,
		args[0],
		binary,
		binary,
	)
}

func main() {
	stdIn := os.Stdin
	stdOut := os.Stdout
	stdErr := os.Stderr

	logrus.SetOutput(stdErr)

	dotCli := command.NewDotCli(stdIn, stdOut, stdErr)
	cmd := newDotCommand(dotCli)

	if err := cmd.Execute(); err != nil {
		if sterr, ok := err.(cli.StatusError); ok {
			if sterr.Status != "" {
				fmt.Fprintln(stdErr, sterr.Status)
			}
			// StatusError should only be used for errors, and all errors should
			// have a non-zero exit status, so never exit with 0
			if sterr.StatusCode == 0 {
				os.Exit(1)
			}
			os.Exit(sterr.StatusCode)
		}
		fmt.Fprintln(stdErr, err)
		os.Exit(1)
	}
}

func dotPreRun(opts *cliflags.Options) {
	// cliflags.SetLogLevel(opts.LogLevel)
	if opts.ConfigDir != "" {
		cliconfig.SetDir(opts.ConfigDir)
	}
	// if opts.Debug {
	// 	debug.Enable()
	// }
}

// type versionDetails interface {
// 	Client() client.APIClient
// 	ServerInfo() command.ServerInfo
// }

// func hideUnsupportedFeatures(cmd *cobra.Command, details versionDetails) {
// 	clientVersion := details.Client().ClientVersion()
// 	osType := details.ServerInfo().OSType
// 	hasExperimental := details.ServerInfo().HasExperimental

// 	cmd.Flags().VisitAll(func(f *pflag.Flag) {
// 		// hide experimental flags
// 		if !hasExperimental {
// 			if _, ok := f.Annotations["experimental"]; ok {
// 				f.Hidden = true
// 			}
// 		}

// 		// hide flags not supported by the server
// 		if !isOSTypeSupported(f, osType) || !isVersionSupported(f, clientVersion) {
// 			f.Hidden = true
// 		}
// 	})

// 	for _, subcmd := range cmd.Commands() {
// 		// hide experimental subcommands
// 		if !hasExperimental {
// 			if _, ok := subcmd.Tags["experimental"]; ok {
// 				subcmd.Hidden = true
// 			}
// 		}

// 		// hide subcommands not supported by the server
// 		if subcmdVersion, ok := subcmd.Tags["version"]; ok && versions.LessThan(clientVersion, subcmdVersion) {
// 			subcmd.Hidden = true
// 		}
// 	}
// }

// func isSupported(cmd *cobra.Command, details versionDetails) error {
// 	clientVersion := details.Client().ClientVersion()
// 	osType := details.ServerInfo().OSType
// 	hasExperimental := details.ServerInfo().HasExperimental

// 	// Check recursively so that, e.g., `docker stack ls` returns the same output as `docker stack`
// 	for curr := cmd; curr != nil; curr = curr.Parent() {
// 		if cmdVersion, ok := curr.Tags["version"]; ok && versions.LessThan(clientVersion, cmdVersion) {
// 			return fmt.Errorf("%s requires API version %s, but the Docker daemon API version is %s", cmd.CommandPath(), cmdVersion, clientVersion)
// 		}
// 		if _, ok := curr.Tags["experimental"]; ok && !hasExperimental {
// 			return fmt.Errorf("%s is only supported on a Docker daemon with experimental features enabled", cmd.CommandPath())
// 		}
// 	}

// 	errs := []string{}

// 	cmd.Flags().VisitAll(func(f *pflag.Flag) {
// 		if f.Changed {
// 			if !isVersionSupported(f, clientVersion) {
// 				errs = append(errs, fmt.Sprintf("\"--%s\" requires API version %s, but the Docker daemon API version is %s", f.Name, getFlagAnnotation(f, "version"), clientVersion))
// 				return
// 			}
// 			if !isOSTypeSupported(f, osType) {
// 				errs = append(errs, fmt.Sprintf("\"--%s\" requires the Docker daemon to run on %s, but the Docker daemon is running on %s", f.Name, getFlagAnnotation(f, "ostype"), osType))
// 				return
// 			}
// 			if _, ok := f.Annotations["experimental"]; ok && !hasExperimental {
// 				errs = append(errs, fmt.Sprintf("\"--%s\" is only supported on a Docker daemon with experimental features enabled", f.Name))
// 			}
// 		}
// 	})
// 	if len(errs) > 0 {
// 		return errors.New(strings.Join(errs, "\n"))
// 	}

// 	return nil
// }

// func getFlagAnnotation(f *pflag.Flag, annotation string) string {
// 	if value, ok := f.Annotations[annotation]; ok && len(value) == 1 {
// 		return value[0]
// 	}
// 	return ""
// }

// func isVersionSupported(f *pflag.Flag, clientVersion string) bool {
// 	if v := getFlagAnnotation(f, "version"); v != "" {
// 		return versions.GreaterThanOrEqualTo(clientVersion, v)
// 	}
// 	return true
// }

// func isOSTypeSupported(f *pflag.Flag, osType string) bool {
// 	if v := getFlagAnnotation(f, "ostype"); v != "" && osType != "" {
// 		return osType == v
// 	}
// 	return true
// }

// // hasTags return true if any of the command's parents has tags
// func hasTags(cmd *cobra.Command) bool {
// 	for curr := cmd; curr != nil; curr = curr.Parent() {
// 		if len(curr.Tags) > 0 {
// 			return true
// 		}
// 	}

// 	return false
// }
