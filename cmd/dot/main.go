package main

import (
	"fmt"
	"net/url"
	"os"
	"runtime"
	"strings"

	"github.com/LEI/dot/internal/dot"
	"github.com/LEI/dot/internal/git"
	"github.com/LEI/dot/internal/host"
	"github.com/LEI/dot/internal/pkg"
	"github.com/spf13/cobra"
)

type appError struct {
	Err  error
	Msg  string
	Code int
}

func (e *appError) Error() string {
	// s := fmt.Sprintf("exited with code %d, and message '%s'", e.Code, e.Msg)
	s := e.Msg
	if e.Err != nil {
		s += fmt.Sprintf(": %s", e.Err)
	}
	return s
}

var (
	// Root command options
	flagVersion bool
	flagRelease bool

	// Extra environment variables
	extraEnv = map[string]string{
		"OS": runtime.GOOS,
		// DOT_...
	}
)

var dotConfig = &dot.Config{}

// cmdRoot is the base command when no other command has been specified.
var cmdRoot = &cobra.Command{
	Use:   binary,
	Short: "Manage files",
	Long:  binary + ` is yet another file manager program.`,
	Args:  cobra.NoArgs,
	// Cannot be overriden
	PersistentPreRunE: persistentPreRunDot,
	RunE:              runDot,
	SilenceErrors:     true,
	SilenceUsage:      true,
	DisableAutoGenTag: true,
	TraverseChildren:  true,
}

// var logBuffer = bytes.NewBuffer(nil)

func init() {
	// // install custom global logger into a buffer, if an error occurs
	// // we can show the logs
	// log.SetOutput(logBuffer)

	flags := cmdRoot.Flags()
	flags.BoolVarP(&flagVersion, "version", "V", flagVersion, "Print short version information")
	flags.BoolVarP(&flagRelease, "release", "", flagRelease, "Print parsed release information and OS types")

	// cmdRoot.SetVersionTemplate("dot version {{.Version}}\n")
}

func shouldLoadConfig(cmd *cobra.Command) bool {
	// Conditions
	not := []bool{
		cmd.Name() == "version",
		flagVersion,
		flagRelease,
	}
	for _, n := range not {
		if n {
			return false
		}
	}
	return true
}

func persistentPreRunDot(cmd *cobra.Command, args []string) error {
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
	// if dotOps.DryRun {
	// 	os.Setenv("DOT_DRY_RUN", "1")
	// }
	if err := setupGlobalOptions(&dotOpts); err != nil {
		return err
	}
	// parse extended options
	// opts, err := options.Parse(dotOpts.Options)
	// if err != nil {
	// 	return err
	// }
	// dotOpts.extended = opts
	if !shouldLoadConfig(cmd) {
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
		fmt.Fprintf(dotOpts.stdout, "%s", str)
		return nil
	}
	if flagRelease {
		ostypes := host.GetOSTypes()
		fmt.Fprintf(dotOpts.stdout, "%s\n", ostypes)
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
	if !opts.HTTPS {
		git.Scheme = "ssh"
		git.User = url.User("git")
	}
	git.DryRun = opts.DryRun
	git.Quiet = opts.verbosity == 0
	git.Verbose = opts.Verbose
	pkg.DryRun = opts.DryRun
	// pkg.Update = true
	return nil
}

func setupGlobalConfig(cfg *dot.Config) error {
	// Filter roles according to env and flags
	if err := setupConfigRoles(cfg); err != nil {
		return err
	}
	if cfg.Git != nil {
		if cfg.Git.Scheme != "" {
			git.Scheme = cfg.Git.Scheme
		}
		if cfg.Git.Host != "" {
			git.Host = cfg.Git.Host
		}
		if cfg.Git.User != nil &&
			cfg.Git.User.String() != "" {
			git.User = cfg.Git.User
		}
	}
	if err := cfg.ParseRoles(); err != nil {
		return err
	}
	dotConfig = cfg
	return nil
}

func setupConfigRoles(cfg *dot.Config) (err error) {
	roles := cfg.Roles[:0] // []*dot.Role{}
	// Filter roles by platform
	for _, r := range cfg.Roles {
		if len(r.OS) > 0 && !host.HasOS(r.OS...) {
			continue
		}
		roles = append(roles, r)
	}
	// Filter roles by name
	if len(dotOpts.RoleFilter) > 0 {
		roles, err = filterRoles(roles, dotOpts.RoleFilter)
		if err != nil {
			return err
		}
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
	return nil
}

// filterRoles by name
func filterRoles(roles []*dot.Role, names []string) ([]*dot.Role, error) {
	matched := []*dot.Role{} // FIXME roles[:0]
	for _, s := range names {
		c := len(matched)
		for _, r := range roles {
			if s == r.Name {
				matched = append(matched, r)
				break
			}
		}
		if c == len(matched) {
			return roles, fmt.Errorf("%s: role not found", s)
		}
	}
	return matched, nil
}

func main() {
	// Setup execution environment
	for k, v := range extraEnv {
		o := os.Getenv(k)
		if o == v {
			continue
		}
		// TODO: !ok os.Unsetenv?
		defer os.Setenv(k, o)
		if err := os.Setenv(k, v); err != nil {
			panic(err)
		}
	}

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
		if aerr, ok := err.(*appError); ok {
			exitCode = aerr.Code
		} else {
			exitCode = 1
		}
	}

	os.Exit(exitCode)
}
