//go:build integration
// +build integration

package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"testing"

	"github.com/spf13/cobra"
)

// TODO: test main no confirm

func testCobra(cmd *cobra.Command, args []string) error {
	var buf bytes.Buffer
	cmd.SetArgs(args)
	cmd.SetOutput(buf)
	err := cmd.Execute()
	if err != nil {
		fmt.Print(buf)
	}
	return err
}

// func TestRunCmd(t *testing.T) {
// 	args := []string{}
// 	if err := testCobra(cmdRoot, args); err != nil {
// 		t.Fatal(err)
// 	}
// }

func TestRunListCmd(t *testing.T) {
	args := []string{"list"} // "--format", "{{.Name}}"
	if err := testCobra(cmdRoot, args); err != nil {
		t.Fatalf("list command failed: %#v", err)
		return
	}
}

func TestInstallLinkCmd(t *testing.T) {
	name := "link"
	tmpDir, err := ioutil.TempDir("", name)
	if err != nil {
		t.Fatalf("create %s tempdir: %s", name, err)
		return
	}
	defer os.RemoveAll(tmpDir)
	args := []string{"install", "--target", tmpDir, "link", "--dry-run"}
	if err := testCobra(cmdRoot, args); err != nil {
		t.Fatalf("install %s command failed: %#v", name, err)
		return
	}
}
