// +build doc

package main

import (
	"github.com/spf13/cobra"
	"github.com/spf13/cobra/doc"
)

var (
	manFlag bool
)

func init() {
	flags := cmdDoc.Flags()
	flags.BoolVarP(&manFlag, "man", "", manFlag, "generate man pages")
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
