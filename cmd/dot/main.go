package main

import (
	"fmt"
	"os"

	"github.com/LEI/dot/internal/dot"
	"github.com/spf13/cobra"
)

var flagVersion bool

var globalConfig *dot.Config

// cmdRoot is the base command when no other command has been specified.
var cmdRoot = &cobra.Command{
	Use:   binary,
	Short: "Manage files",
	Long: `
dot is yet another file manager program.
`,
	SilenceErrors:     true,
	SilenceUsage:      true,
	DisableAutoGenTag: true,
	// TraverseChildren:  true,
	// Args:              cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		if flagVersion {
			fmt.Printf("%s %s\n", binary, version)
			return
		}
		// if err := cmd.Usage(); err != nil {
		// 	return err
		// }
		// cmd.SetOutput(os.Stderr)
		// cmd.HelpFunc()(cmd, args)
		cmd.Help()
	},
	PersistentPreRunE: func(c *cobra.Command, args []string) error {
		// set verbosity, default is one
		dotOpts.verbosity = 1
		if dotOpts.Quiet && dotOpts.Verbose > 1 {
			return fmt.Errorf("--quiet and --verbose cannot be specified at the same time")
		}

		switch {
		case dotOpts.Verbose >= 2:
			dotOpts.verbosity = 3
		case dotOpts.Verbose > 0:
			dotOpts.verbosity = 2
		case dotOpts.Quiet:
			dotOpts.verbosity = 0
		}

		// parse extended options
		// opts, err := options.Parse(dotOpts.Options)
		// if err != nil {
		// 	return err
		// }
		// dotOpts.extended = opts
		if c.Name() == "version" {
			return nil
		}
		// pwd, err := resolvePassword(dotOpts, "RESTIC_PASSWORD")
		// if err != nil {
		// 	fmt.Fprintf(os.Stderr, "Resolving password failed: %v\n", err)
		// 	Exit(1)
		// }
		// dotOpts.password = pwd

		cfg, err := OpenConfig(dotOpts)
		if err != nil {
			return err
		}
		globalConfig = cfg
		if len(dotOpts.RoleFilter) > 0 {
			roles := []*dot.Role{}
			for _, r := range globalConfig.Roles {
				for _, s := range dotOpts.RoleFilter {
					if s == r.Name {
						roles = append(roles, r)
						break
					}
				}
			}
			globalConfig.Roles = roles
		}
		if len(globalConfig.Roles) == 0 {
			return fmt.Errorf("no roles to execute")
		}
		if err := globalConfig.Parse(); err != nil {
			return err
		}
		// if err := globalConfig.Load(); err != nil {
		// 	return err
		// }

		// // run the debug functions for all subcommands (if build tag "debug" is
		// // enabled)
		// if err := runDebug(); err != nil {
		// 	return err
		// }

		return nil
	},
}

// var logBuffer = bytes.NewBuffer(nil)

func init() {
	// // install custom global logger into a buffer, if an error occurs
	// // we can show the logs
	// log.SetOutput(logBuffer)

	flags := cmdRoot.Flags()
	flags.BoolVarP(&flagVersion, "version", "V", flagVersion, "Print short version information")
}

func main() {
	// debug.Log("main %#v", os.Args)
	// debug.Log("restic %s compiled with %v on %v/%v",
	// 	version, runtime.Version(), runtime.GOOS, runtime.GOARCH)
	err := cmdRoot.Execute()

	switch {
	// case restic.IsAlreadyLocked(errors.Cause(err)):
	// 	fmt.Fprintf(os.Stderr, "%v\nthe `unlock` command can be used to remove stale locks\n", err)
	// case errors.IsFatal(errors.Cause(err)):
	// 	fmt.Fprintf(os.Stderr, "%v\n", err)
	case err != nil:
		fmt.Fprintf(os.Stderr, "%+v\n", err)

		// if logBuffer.Len() > 0 {
		// 	fmt.Fprintf(os.Stderr, "also, the following messages were logged by a library:\n")
		// 	sc := bufio.NewScanner(logBuffer)
		// 	for sc.Scan() {
		// 		fmt.Fprintln(os.Stderr, sc.Text())
		// 	}
		// }
	}

	var exitCode int
	if err != nil {
		exitCode = 1
	}

	os.Exit(exitCode)
}
