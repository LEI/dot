package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/LEI/dot/internal/dot"
	"github.com/LEI/dot/internal/git"
	"github.com/LEI/dot/internal/ostype"
	"github.com/LEI/dot/internal/pkg"
	"github.com/spf13/cobra"
)

var flagVersion bool

var dotConfig *dot.Config = &dot.Config{}

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
	TraverseChildren:  true,
	Args:              cobra.NoArgs,
	// Cannot be overriden
	PersistentPreRunE: persistentPreRunDot,
	RunE:              runDot,
}

// var logBuffer = bytes.NewBuffer(nil)

func init() {
	// // install custom global logger into a buffer, if an error occurs
	// // we can show the logs
	// log.SetOutput(logBuffer)

	flags := cmdRoot.Flags()
	flags.BoolVarP(&flagVersion, "version", "V", flagVersion, "Print short version information")

	// cmdRoot.SetVersionTemplate("dot version {{.Version}}\n")
}

func persistentPreRunDot(c *cobra.Command, args []string) error {
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
	if err := setupGlobalOptions(&dotOpts); err != nil {
		return err
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
	if err := setupGlobalConfig(cfg); err != nil {
		return err
	}
	// if err := dotConfig.Load(); err != nil {
	// 	return err
	// }

	// // run the debug functions for all subcommands (if build tag "debug" is
	// // enabled)
	// if err := runDebug(); err != nil {
	// 	return err
	// }

	return nil
}

func runDot(cmd *cobra.Command, args []string) error {
	if flagVersion {
		str := fmt.Sprintf("%s %s", binary, version)
		// if commit != "" {
		// 	str += fmt.Sprintf(" (%s)", commit)
		// }
		fmt.Println(str)
		return nil
	}
	// if err := cmd.Usage(); err != nil {
	// 	return err
	// }
	// cmd.SetOutput(os.Stderr)
	// cmd.HelpFunc()(cmd, args)
	usage := strings.TrimRight(cmd.UsageString(), "\n")
	return fmt.Errorf("%s", usage)
	// return fmt.Errorf("%s", cmd.UseLine())
}

func setupGlobalOptions(opts *DotOptions) error {
	dot.DecodeErrorUnused = !opts.Force
	git.DryRun = opts.DryRun
	git.Quiet = opts.verbosity == 0
	git.Verbose = opts.Verbose
	pkg.DryRun = opts.DryRun
	// pkg.Update = true
	return nil
}

func setupGlobalConfig(cfg *dot.Config) error {
	roles := cfg.Roles[:0] // []*dot.Role{}
	// Filter roles by platform
	for _, r := range cfg.Roles {
		if len(r.OS) > 0 && !ostype.Has(r.OS...) {
			continue
		}
		roles = append(roles, r)
	}
	// Filter roles by name
	if len(dotOpts.RoleFilter) > 0 {
		roles = filterRoles(roles, dotOpts.RoleFilter)
		// tmp := roles[:0] // []*dot.Role{}
		// for _, r := range roles {
		// 	for _, s := range dotOpts.RoleFilter {
		// 		if s == r.Name {
		// 			tmp = append(tmp, r)
		// 			break
		// 		}
		// 	}
		// }
		// roles = tmp
	}
	if len(roles) == 0 {
		msg := "nothing to do"
		msg += fmt.Sprintf(" with %d roles", len(roles))
		if len(dotOpts.RoleFilter) > 0 {
			msg += fmt.Sprintf(" and filter %s", dotOpts.RoleFilter)
		}
		return fmt.Errorf(msg)
	}
	cfg.Roles = roles
	if err := cfg.ParseRoles(); err != nil {
		return err
	}
	dotConfig = cfg
	return nil
}

// filterRoles by name
func filterRoles(roles []*dot.Role, names []string) []*dot.Role {
	res := roles[:0]
	for _, r := range roles {
		for _, s := range names {
			if s == r.Name {
				res = append(res, r)
				break
			}
		}
	}
	return res
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
