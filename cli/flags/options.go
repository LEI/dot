package flags

import (
	"fmt"
	"os"

	cliconfig "github.com/LEI/dot/cli/config"
	"github.com/sirupsen/logrus"
	"github.com/spf13/pflag"
)

// Options are the options used to configure the cli
type Options struct {
	// LogLevel string
	Version bool

	ConfigDir     string
	DryRun, Force bool
	Verbose       int

	RoleFilter []string

	Source, Target string
}

// InstallFlags adds flags for the common options on the FlagSet
func (opts *Options) InstallFlags(flags *pflag.FlagSet) {
	// flags.StringVarP(&opts.LogLevel, "log-level", "l", "info", `Set the logging level ("debug"|"info"|"warn"|"error"|"fatal")`)
	flags.BoolVarP(&opts.Version, "version", "V", false, "Print version information and quit")
}

// InstallPersistentFlags adds persistent flags on the FlagSet
func (opts *Options) InstallPersistentFlags(flags *pflag.FlagSet) {
	flags.StringVarP(&opts.ConfigDir, "config", "c", cliconfig.Dir(), "Location of config file") // (s)
	flags.BoolVarP(&opts.DryRun, "dry-run", "d", false, "Do not execute tasks")
	flags.BoolVarP(&opts.Force, "force", "f", opts.Force, "Force execution even if the repository is dirty")
	flags.CountVarP(&opts.Verbose, "verbose", "v", "Verbosity level")

	flags.StringSliceVarP(&opts.RoleFilter, "role", "r", []string{}, "Filter role execution")
}

// InstallActionFlags adds common task falgs on the FlagSet
func (opts *Options) InstallActionFlags(flags *pflag.FlagSet) {
}

// InstallActionPersistentFlags adds common task falgs on the FlagSet
func (opts *Options) InstallActionPersistentFlags(flags *pflag.FlagSet) {
	flags.StringVarP(&opts.Source, "source", "s", "", "Source directory")
	flags.StringVarP(&opts.Target, "target", "t", "", "Target directory")
}

// InstallTaskFlags adds common task falgs on the FlagSet
func (opts *Options) InstallTaskFlags(flags *pflag.FlagSet) {
}

// InstallTaskPersistentFlags adds common task falgs on the FlagSet
func (opts *Options) InstallTaskPersistentFlags(flags *pflag.FlagSet) {
}

// SetLogLevel sets the logrus logging level
func SetLogLevel(logLevel string) {
	if logLevel != "" {
		lvl, err := logrus.ParseLevel(logLevel)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Unable to parse logging level: %s\n", logLevel)
			os.Exit(1)
		}
		logrus.SetLevel(lvl)
	} else {
		logrus.SetLevel(logrus.InfoLevel)
	}
}
