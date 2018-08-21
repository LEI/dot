package main

import (
	"fmt"
	"runtime"

	"github.com/spf13/cobra"
)

var versionCmd = &cobra.Command{
	Use: "version",
	// Aliases: []string{"v"},
	Short: "Print version information",
	Long: `
The "version" command prints detailed information about the build environment
and the version of this software.
`,
	DisableAutoGenTag: true,
	Args:              cobra.NoArgs,
	Run:               runVersion,
}

func init() {
	cmdRoot.AddCommand(versionCmd)
}

func runVersion(cmd *cobra.Command, args []string) {
	ver := version
	if commit != "" {
		ver = fmt.Sprintf("%s (%s)", ver, commit)
	}
	// Add info (OS, arch and time)
	info := fmt.Sprintf("%v/%v", runtime.GOOS, runtime.GOARCH)
	if dotOpts.verbosity >= 2 {
		info += " at " + date
	}
	fmt.Printf("%s %s compiled with %v on %s\n",
		binary, ver, runtime.Version(), info)
}
