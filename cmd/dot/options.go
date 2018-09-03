package main

import (
	"fmt"
	"io"
	"os"

	"github.com/LEI/dot/internal/dot"
	"github.com/LEI/dot/internal/shell"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

var (
	binary = "dot" // _, binary = filepath.Split(os.Args[0])

	// LDFlags
	version = "compiled manually"
	commit  string
	date    string

	homeDir = shell.HomeDir
)

// DotOptions hold all global options for restic.
type DotOptions struct {
	pkg bool // Enable package managment

	Source     string
	Target     string
	ConfigFile string
	RoleDir    string // name
	RoleFile   string // name
	RoleFilter []string
	DryRun     bool
	Force      bool
	Quiet      bool
	Verbose    int
	// CacheDir     string
	// NoCache      bool
	// CleanupCache bool

	// ctx      context.Context
	// private string
	// configFileUsed string
	stdout io.Writer
	stderr io.Writer

	// verbosity is set as follows:
	//  0 means: don't print any messages except errors, this is used when --quiet is specified
	//  1 is the default: print essential messages
	//  2 means: print more messages, report minor things, this is used when --verbose is specified
	//  3 means: print very detailed debug messages, this is used when --debug is specified
	verbosity uint

	// Options []string

	// extended options.Options
}

var dotOpts = DotOptions{
	stdout: os.Stdout,
	stderr: os.Stderr,
}

func init() {
	pflags := cmdRoot.PersistentFlags()
	pflags.BoolVarP(&dotOpts.Quiet, "quiet", "q", false, "do not output") // comprehensive progress report
	pflags.CountVarP(&dotOpts.Verbose, "verbose", "v", "be verbose (specify --verbose multiple times or level `n`)")
	// pflags.StringVar(&dotOpts.CacheDir, "cache-dir", "", "set the cache directory")
	// pflags.BoolVar(&dotOpts.NoCache, "no-cache", false, "do not use a local cache")
	// pflags.BoolVar(&dotOpts.CleanupCache, "cleanup-cache", false, "auto remove old cache directories")
	// pflags.StringSliceVarP(&dotOpts.Options, "option", "o", []string{}, "set extended option (`key=value`, can be specified multiple times)")
}

func addActionFlags(cmd *cobra.Command) {
	source := os.Getenv("DOT_SOURCE")
	if source == "" {
		source = homeDir
	}
	target := os.Getenv("DOT_TARGET")
	if target == "" {
		target = homeDir
	}
	// Dotfile
	envCfgFile := os.Getenv("DOT_FILE")
	if envCfgFile == "" {
		// envCfgFile = filepath.Join(homeDir, ".dotrc.yml")
		envCfgFile = ".dotrc.yml" // TODO auto ext detection
	}
	// Roles directory
	envRoleDir := os.Getenv("DOT_ROLE_DIR")
	if envRoleDir == "" {
		// envRoleDir = filepath.Join(homeDir, ".dot")
		envRoleDir = ".dot"
	}
	// Roles config name
	envRoleFile := os.Getenv("DOT_ROLE_FILE")
	if envRoleFile == "" {
		envRoleFile = ".dot.yml"
	}

	pflags := cmd.PersistentFlags()
	pflags.StringVarP(&dotOpts.Source, "source", "s", source, "`DOT_SOURCE` directory")
	pflags.StringVarP(&dotOpts.Target, "target", "t", target, "`DOT_TARGET` directory")
	pflags.StringVarP(&dotOpts.ConfigFile, "config-file", "c", envCfgFile, "main configuration `DOT_FILE`")
	pflags.StringVarP(&dotOpts.RoleDir, "role-dir", "", envRoleDir, "roles `DOT_ROLE_DIR`")
	pflags.StringVarP(&dotOpts.RoleFile, "role-file", "", envRoleFile, "roles `DOT_ROLE_FILE`")
	pflags.StringSliceVarP(&dotOpts.RoleFilter, "role-filter", "r", []string{}, "filter roles by name")
	pflags.BoolVarP(&dotOpts.DryRun, "dry-run", "d", false, "do not execute tasks")
	pflags.BoolVarP(&dotOpts.Force, "force", "F", false, "force execution")

	hidden := []string{"role-dir", "role-file"}
	if err := markHidden(pflags, hidden); err != nil {
		panic(err)
	}

	flags := cmd.Flags()
	flags.BoolVarP(&dotOpts.pkg, "packages", "P", false, "manage system packages")
}

func markHidden(f *pflag.FlagSet, in []string) error {
	for _, n := range in {
		if err := f.MarkHidden(n); err != nil {
			return err
		}
	}
	return nil
}

func setActionEnv(cmd *cobra.Command) error {
	verbose := "0"
	if dotOpts.Verbose > 0 { // verbosity > 1?
		verbose = "1"
	}
	dotPkg := "0"
	if dotOpts.pkg {
		dotPkg = "1"
	}
	vars := map[string]string{
		"DOT_ACTION": cmd.Name(),
		// "DOT_DRY_RUN": dotOpts.DryRun,
		"DOT_VERBOSE":   verbose,
		"DOT_SOURCE":    dotOpts.Source,
		"DOT_TARGET":    dotOpts.Target,
		"DOT_FILE":      dotOpts.ConfigFile, // Dotfile
		"DOT_ROLE_DIR":  dotOpts.RoleDir,    // Roles directory
		"DOT_ROLE_FILE": dotOpts.RoleFile,   // Roles config name
		"DOT_PKG":       dotPkg,
	}
	for k, v := range vars {
		// if dotOpts.verbosity >= 3 {
		// 	fmt.Fprintf(dotOpts.stdout, "%s=%q\n", k, v)
		// }
		err := os.Setenv(k, v)
		if err != nil {
			return err
		}
	}
	return nil
}

// // checkErrno returns nil when err is set to syscall.Errno(0), since this is no
// // error condition.
// func checkErrno(err error) error {
// 	e, ok := err.(syscall.Errno)
// 	if !ok {
// 		return err
// 	}

// 	if e == 0 {
// 		return nil
// 	}

// 	return err
// }

// func stdinIsTerminal() bool {
// 	return terminal.IsTerminal(int(os.Stdin.Fd()))
// }

// func stdoutIsTerminal() bool {
// 	return terminal.IsTerminal(int(os.Stdout.Fd()))
// }

// func stdoutTerminalWidth() int {
// 	w, _, err := terminal.GetSize(int(os.Stdout.Fd()))
// 	if err != nil {
// 		return 0
// 	}
// 	return w
// }

// // ClearLine creates a platform dependent string to clear the current
// // line, so it can be overwritten. ANSI sequences are not supported on
// // current windows cmd shell.
// func ClearLine() string {
// 	if runtime.GOOS == "windows" {
// 		if w := stdoutTerminalWidth(); w > 0 {
// 			return strings.Repeat(" ", w-1) + "\r"
// 		}
// 		return ""
// 	}
// 	return "\x1b[2K"
// }

// // Printf writes the message to the configured stdout stream.
// func Printf(format string, args ...interface{}) {
// 	_, err := fmt.Fprintf(dotOpts.stdout, format, args...)
// 	if err != nil {
// 		fmt.Fprintf(os.Stderr, "unable to write to stdout: %v\n", err)
// 		os.Exit(100)
// 	}
// }

// // Verbosef calls Printf to write the message when the verbose flag is set.
// func Verbosef(format string, args ...interface{}) {
// 	if dotOpts.verbosity >= 1 {
// 		Printf(format, args...)
// 	}
// }

// // PrintProgress wraps fmt.Printf to handle the difference in writing progress
// // information to terminals and non-terminal stdout
// func PrintProgress(format string, args ...interface{}) {
// 	var (
// 		message         string
// 		carriageControl string
// 	)
// 	message = fmt.Sprintf(format, args...)

// 	if !(strings.HasSuffix(message, "\r") || strings.HasSuffix(message, "\n")) {
// 		if stdoutIsTerminal() {
// 			carriageControl = "\r"
// 		} else {
// 			carriageControl = "\n"
// 		}
// 		message = fmt.Sprintf("%s%s", message, carriageControl)
// 	}

// 	if stdoutIsTerminal() {
// 		message = fmt.Sprintf("%s%s", ClearLine(), message)
// 	}

// 	fmt.Print(message)
// }

// // Warnf writes the message to the configured stderr stream.
// func Warnf(format string, args ...interface{}) {
// 	_, err := fmt.Fprintf(dotOpts.stderr, format, args...)
// 	if err != nil {
// 		fmt.Fprintf(os.Stderr, "unable to write to stderr: %v\n", err)
// 		os.Exit(100)
// 	}
// }

// // Exitf uses Warnf to write the message and then terminates the process with
// // the given exit code.
// func Exitf(exitcode int, format string, args ...interface{}) {
// 	if format[len(format)-1] != '\n' {
// 		format += "\n"
// 	}

// 	Warnf(format, args...)
// 	os.Exit(exitcode)
// }

// const maxKeys = 20

// OpenConfig ...
func OpenConfig(opts DotOptions) (*dot.Config, error) {
	if opts.ConfigFile == "" {
		// empty config file name or path
		return nil, fmt.Errorf("Please specify config file location (-c)")
	}

	cfg, err := dot.NewConfig(opts.ConfigFile, opts.RoleDir)
	if err != nil {
		return nil, err
	}
	if opts.DryRun && opts.Verbose >= 0 {
		fmt.Fprintf(dotOpts.stderr, "DRY-RUN enabled\n")
	}
	if cf := cfg.FileUsed(); cf != "" && dotOpts.verbosity >= 2 {
		fmt.Fprintf(dotOpts.stdout, "Using config file: %s\n", cf)
	}
	// Note: default flag values will always override config
	if opts.Source != "" {
		cfg.Source = opts.Source
	}
	if opts.Target != "" {
		cfg.Target = opts.Target
	}
	// if opts.RoleDir != "" {
	// 	cfg.SetDir(opts.RoleDir)
	// }
	if opts.RoleFile != "" {
		cfg.SetRoleFile(opts.RoleFile)
	}
	// s := repository.New(be)

	// opts.password, err = ReadPassword(opts, "enter password for repository: ")
	// if err != nil {
	// 	return nil, err
	// }

	// // err = s.SearchKey(opts.ctx, opts.password, maxKeys)
	// // if err != nil {
	// // 	return nil, err
	// // }

	// if stdoutIsTerminal() {
	// 	id := s.Config().ID
	// 	if len(id) > 8 {
	// 		id = id[:8]
	// 	}
	// 	Verbosef("repository %v opened successfully, password is correct\n", id)
	// }

	// if opts.NoCache {
	// 	return s, nil
	// }

	// c, err := cache.New(s.Config().ID, opts.CacheDir)
	// if err != nil {
	// 	Warnf("unable to open cache: %v\n", err)
	// 	return s, nil
	// }

	// // start using the cache
	// s.UseCache(c)

	// oldCacheDirs, err := cache.Old(c.Base)
	// if err != nil {
	// 	Warnf("unable to find old cache directories: %v", err)
	// }

	// // nothing more to do if no old cache dirs could be found
	// if len(oldCacheDirs) == 0 {
	// 	return s, nil
	// }

	// // cleanup old cache dirs if instructed to do so
	// if opts.CleanupCache {
	// 	Printf("removing %d old cache dirs from %v\n", len(oldCacheDirs), c.Base)

	// 	for _, item := range oldCacheDirs {
	// 		dir := filepath.Join(c.Base, item.Name())
	// 		err = fs.RemoveAll(dir)
	// 		if err != nil {
	// 			Warnf("unable to remove %v: %v\n", dir, err)
	// 		}
	// 	}
	// } else {
	// 	if stdoutIsTerminal() {
	// 		Verbosef("found %d old cache directories in %v, pass --cleanup-cache to remove them\n",
	// 			len(oldCacheDirs), c.Base)
	// 	}
	// }

	return cfg, nil
}
