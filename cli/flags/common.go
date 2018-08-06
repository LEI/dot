package cli

import (
	"fmt"
	"os"

	"github.com/sirupsen/logrus"
	"github.com/spf13/pflag"
)

// Options are the options used to configure the cli
type Options struct { // (cliflags.ClientOptions)
	ConfigDir string
	Debug bool
	LogLevel string
	Version bool
}

// InstallFlags adds flags for the common options on the FlagSet
func (cmdOpts *Options) InstallFlags(flags *pflag.FlagSet) {
	flags.BoolVarP(&cmdOpts.Debug, "debug", "D", false, "Enable debug mode")
	flags.StringVarP(&cmdOpts.LogLevel, "log-level", "l", "info", `Set the logging level ("debug"|"info"|"warn"|"error"|"fatal")`)
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
