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
	fmt.Printf("%s %s (%s) compiled with %v on %v/%v\n",
		binary, version, commit, runtime.Version(), runtime.GOOS, runtime.GOARCH)
}
