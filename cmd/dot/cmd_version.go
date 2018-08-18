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
	v := version
	if commit != "" {
		v = fmt.Sprintf("%s (%s)", v, commit)
	}
	info := fmt.Sprintf("%v/%v", runtime.GOOS, runtime.GOARCH)
	if dotOpts.Verbose > 0 {
		info += " at " + timestamp
	}
	fmt.Printf("%s %s compiled with %v on %s\n",
		binary, v, runtime.Version(), info)
}
