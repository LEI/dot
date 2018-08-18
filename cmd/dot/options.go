package main

import (
	"fmt"
	"io"
	"os"

	"github.com/LEI/dot/internal/dot"
)

var (
	binary  = "dot" // _, binary = filepath.Split(os.Args[0])
	version = "compiled manually"
)

// GlobalOptions hold all global options for restic.
type GlobalOptions struct {
	Source     string
	Target     string
	ConfigFile string
	RoleDir    string
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

var dotOpts = GlobalOptions{
	stdout: os.Stdout,
	stderr: os.Stderr,
}

func init() {
	source := os.Getenv("DOT_SOURCE")
	if source == "" {
		source = os.Getenv("HOME")
	}
	target := os.Getenv("DOT_TARGET")
	if target == "" {
		target = os.Getenv("HOME")
	}
	// Dotfile
	envCfgFile := os.Getenv("DOT_FILE")
	if envCfgFile == "" {
		// envCfgFile = os.ExpandEnv("$HOME/.dotrc.yml")
		envCfgFile = ".dotrc.yml"
	}
	// Roles directory
	envRoleDir := os.Getenv("DOT_ROLE_DIR")
	if envRoleDir == "" {
		envRoleDir = os.ExpandEnv("$HOME/.dot")
		envRoleDir = ".dot"
	}

	f := cmdRoot.PersistentFlags()
	f.StringVarP(&dotOpts.Source, "source", "s", source, "`DOT_SOURCE` directory")
	f.StringVarP(&dotOpts.Target, "target", "t", target, "`DOT_TARGET` directory")
	f.StringVarP(&dotOpts.ConfigFile, "config-file", "c", envCfgFile, "global configuration `DOT_FILE`")
	f.StringVarP(&dotOpts.RoleDir, "role-dir", "", envRoleDir, "roles `DOT_ROLE_DIR`")
	f.StringSliceVarP(&dotOpts.RoleFilter, "role-filter", "r", []string{}, "filter roles by name")
	f.BoolVarP(&dotOpts.DryRun, "dry-run", "d", false, "do not execute tasks")
	f.BoolVarP(&dotOpts.Force, "force", "F", false, "force execution")
	f.BoolVarP(&dotOpts.Quiet, "quiet", "q", false, "do not output") // comprehensive progress report
	f.CountVarP(&dotOpts.Verbose, "verbose", "v", "be verbose (specify --verbose multiple times or level `n`)")
	// f.StringVar(&dotOpts.CacheDir, "cache-dir", "", "set the cache directory")
	// f.BoolVar(&dotOpts.NoCache, "no-cache", false, "do not use a local cache")
	// f.BoolVar(&dotOpts.CleanupCache, "cleanup-cache", false, "auto remove old cache directories")
	// f.StringSliceVarP(&dotOpts.Options, "option", "o", []string{}, "set extended option (`key=value`, can be specified multiple times)")
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
func OpenConfig(opts GlobalOptions) (*dot.Config, error) {
	if opts.ConfigFile == "" {
		return nil, fmt.Errorf("Please specify config file location (-c)")
	}

	cfg, err := dot.NewConfig(opts.ConfigFile)
	if err != nil {
		return nil, err
	}
	if cfg.Source == "" {
		cfg.Source = opts.Source
	}
	if cfg.Target == "" {
		cfg.Target = opts.Target
	}
	if cfg.DirName == "" {
		cfg.DirName = ".dot"
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
