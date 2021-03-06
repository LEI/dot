package main

import (
	"fmt"
	"io/ioutil"
	"log"
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

	// Debug level
	Debug *log.Logger
	// Info level
	Info *log.Logger
	// Warn level
	Warn *log.Logger
	// Error level
	Error *log.Logger
)

var dotConfig = &dot.Config{}

// cmdRoot is the base command when no other command has been specified.
var cmdRoot = &cobra.Command{
	Use:   binary,
	Short: "dot files manager",
	Long:  binary + ` is yet another dot files manager program.`,
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
		cmd.Name() == "doc",
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
	// log.SetOutput(dotOpts.stdout)

	// set verbosity, default is one
	dotOpts.verbosity = 1
	if dotOpts.Quiet && dotOpts.Verbose > 1 {
		return fmt.Errorf("--quiet and --verbose cannot be specified at the same time")
	}

	logFlag := 0 // log.Ldate|log.Ltime|log.Lshortfile
	Debug = log.New(ioutil.Discard, "DEBUG: ", logFlag)
	Info = log.New(dotOpts.stdout, "INFO: ", logFlag)
	Warn = log.New(dotOpts.stderr, "WARN: ", logFlag)
	Error = log.New(dotOpts.stderr, "ERROR: ", logFlag)

	switch {
	case dotOpts.Verbose >= 2:
		dotOpts.verbosity = 3
		Debug.SetOutput(dotOpts.stderr)
	case dotOpts.Verbose > 0:
		dotOpts.verbosity = 2
		Debug.SetOutput(dotOpts.stderr)
	case dotOpts.Quiet:
		dotOpts.verbosity = 0
		Info.SetOutput(ioutil.Discard)
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
	dot.Stdout = opts.stdout
	dot.Stderr = opts.stderr
	// dot.Stdin = os.Stdin
	if !opts.HTTPS {
		git.Scheme = "ssh"
		git.User = url.User("git")
	}
	git.DryRun = opts.DryRun
	git.Quiet = opts.Quiet // opts.verbosity == 0
	git.Verbose = opts.Verbose
	pkg.DryRun = opts.DryRun
	pkg.Stdout = opts.stdout
	pkg.Stderr = opts.stderr
	// dot.Stdin = os.Stdin
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
	// if err := cfg.ParseRoles(); err != nil {
	if err := cfg.PrepareRoles(); err != nil {
		return err
	}
	dotConfig = cfg
	return nil
}

func setupConfigRoles(cfg *dot.Config) (err error) {
	roles := cfg.Roles[:0] // []*dot.Role{}
	for platform, rs := range cfg.Platforms {
		// if !host.HasOS(platform) {
		// 	continue
		// }
		for _, r := range rs {
			r.OS = append(r.OS, platform)
			cfg.Roles = append(cfg.Roles, r)
		}
	}
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
		msg := fmt.Sprintf("nothing to do with %d roles", len(cfg.Roles))
		if len(dotOpts.RoleFilter) > 0 {
			msg += fmt.Sprintf(" and filter: %s", dotOpts.RoleFilter)
		}
		// fmt.Fprintln(os.Stderr, msg)
		return fmt.Errorf(msg)
	}
	cfg.Roles = roles
	return nil
}

// filterRoles by name
func filterRoles(roles []*dot.Role, names []string) ([]*dot.Role, error) {
	matched := []*dot.Role{} // FIXME roles[:0]

	// First match the indexes
	matchedIndexes := []int{}
	for _, s := range names {
		c := len(matchedIndexes)
		for i, r := range roles {
			if s == r.Name {
				matchedIndexes = append(matchedIndexes, i)
				break
			}
		}
		if c == len(matchedIndexes) {
			return roles, fmt.Errorf("%s: role not found", s)
		}
	}
	// Apppend matches in the original order
	for i, r := range roles {
		for _, j := range matchedIndexes {
			if i == j {
				matched = append(matched, r)
			}
		}
	}

	// for _, s := range names {
	// 	c := len(matched)
	// 	for _, r := range roles {
	// 		if s == r.Name {
	// 			matched = append(matched, r)
	// 			break
	// 		}
	// 	}
	// 	if c == len(matched) {
	// 		return roles, fmt.Errorf("%s: role not found", s)
	// 	}
	// }

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
		if Error != nil {
			Error.Println(err)
		} else {
			fmt.Fprintf(os.Stderr, "err: %+v\n", err)
		}

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
