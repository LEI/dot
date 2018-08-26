// +build doc

package main

import (
	"github.com/spf13/cobra"
	"github.com/spf13/cobra/doc"
)

var (
	mdFlag bool
)

func init() {
	flags := cmdDoc.Flags()
	flags.BoolVarP(&mdFlag, "md", "", mdFlag, "generate markdown files")
}

func genMd(cmd *cobra.Command, path string) error {
	if err := doc.GenMarkdownTree(cmd, path); err != nil {
		return err
	}
	return nil
}
