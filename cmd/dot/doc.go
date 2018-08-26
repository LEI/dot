// +build doc

// This package contains the code for the dot executable.
package main

import (
	"fmt"

	"github.com/spf13/cobra"
)

var defaultDocPath = "/tmp"

var cmdDoc = &cobra.Command{
	Use: "doc",
	// Aliases: []string{},
	Short: "Generate documentation",
	Long:  ``,
	// DisableAutoGenTag: true,
	Args: cobra.MaximumNArgs(1), // cobra.ExactArgs(1),
	RunE: genDoc,
}

func init() {
	cmdRoot.AddCommand(cmdDoc)

	// flags := cmdDoc.Flags()
	// flags.StringVarP(&docFlag, "path", "", docPathFlag, "path to target directory")
}

// Generate documentation files.
func genDoc(cmd *cobra.Command, args []string) error {
	path := defaultDocPath
	if len(args) == 1 {
		path = args[0]
	}
	// if docPathFlag == "" {
	// 	return fmt.Errorf("--path is required")
	// }
	switch {
	case manFlag:
	case mdFlag:
	default:
		manFlag = true
		mdFlag = true
		// return fmt.Errorf("Please specify something to generate (--man, --md or both)")
	}
	if manFlag {
		if err := genMan(cmd, path); err != nil {
			return err
		}
		fmt.Printf("Generated man pages in: %s\n", path)
	}
	if mdFlag {
		if err := genMd(cmd, path); err != nil {
			return err
		}
		fmt.Printf("Generated markdown pages in: %s\n", path)
	}
	return nil
}
