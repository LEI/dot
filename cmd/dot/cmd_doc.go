//go:build doc
// +build doc

// This package contains the code for the dot executable.
package main

// https://godoc.org/github.com/spf13/cobra/doc

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/cobra/doc"
)

var (
	docEnv = map[string]string{
		"DOT_SOURCE": "$HOME",
		"DOT_TARGET": "$HOME",
	}

	manFlag string // ./man
	mdFlag  string // ./docs
)

var cmdDoc = &cobra.Command{
	Use: "doc",
	// Aliases: []string{"docs"},
	Short: "Generate documentation",
	Long:  ``,
	Args:  cobra.NoArgs,
	// cobra.MaximumNArgs(1),
	// cobra.ExactArgs(1),
	Hidden:            true,
	RunE:              runDoc,
	DisableAutoGenTag: true,
}

func init() {
	cmdRoot.AddCommand(cmdDoc)

	flags := cmdDoc.Flags()
	flags.StringVarP(&manFlag, "man-page", "", manFlag, "generate man pages")
	flags.StringVarP(&mdFlag, "markdown", "", mdFlag, "generate markdown files")

	for k, v := range docEnv {
		os.Setenv(k, v)
	}
}

// Generate documentation files.
func runDoc(cmd *cobra.Command, args []string) error {
	c := cmdRoot // cmd.Parent()
	if manFlag == "" && mdFlag == "" {
		return fmt.Errorf("Please specify --man-page and/or --markdown")
	}
	if manFlag != "" {
		if err := genMan(c, manFlag); err != nil {
			return err
		}
		fmt.Fprintf(dotOpts.stdout, "Generated man pages in: %s\n", manFlag)
	}
	if mdFlag != "" {
		if err := genMd(c, mdFlag); err != nil {
			return err
		}
		fmt.Fprintf(dotOpts.stdout, "Generated markdown pages in: %s\n", mdFlag)
	}
	return nil
}

// cmd *cobra.Command, args []string
func genMan(cmd *cobra.Command, path string) error {
	if err := genManTree(cmd, path); err != nil {
		return err
	}
	return nil
}

func genManTree(cmd *cobra.Command, path string) error {
	// header := &doc.GenManHeader{
	// 	Title:   strings.ToUpper(binary),
	// 	Section: "1",
	// }
	// err := doc.GenManTree(cmd, header, "/tmp")
	err := doc.GenManTree(cmd, nil, path)
	if err != nil {
		return err
	}
	return nil
}

func genMd(cmd *cobra.Command, path string) error {
	if err := doc.GenMarkdownTree(cmd, path); err != nil {
		return err
	}
	return nil
}
