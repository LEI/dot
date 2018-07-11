package cmd_test

import (
	// "fmt"
	// "os"
	"testing"

	// "github.com/LEI/dot/cmd"

	// "github.com/jessevdk/go-flags"
)

// TestInstallCmd ...
func TestInstallLinkCmd(t *testing.T) {
	/*
	// cmd := &cmd.InstallCmd{
	// 	// Link: cmd.LinkCmd{
	// 	// 	&cmd.LinkArg{
	// 	// 		Source: "a",
	// 	// 		Target: "b",
	// 	// 	},
	// 	// },
	// }
	args := []string{"install", "link", "-s", "a", "-t", "b"}
	parser := flags.NewParser(&cmd.Options, flags.Default)
	_, err := parser.ParseArgs(args)
	// err := cmd.Execute(args)
	if err != nil {
		t.Errorf("Command error: %v", err)
	}
	if cmd.Options.Install.Link.Source != "a" {
		t.Errorf("Expected a, but got %v", cmd.Options.Install.Link.Source)
	}
	if cmd.Options.Install.Link.Target != "b" {
		t.Errorf("Expected b, but got %v", cmd.Options.Install.Link.Target)
	}
	if len(dot.CacheStore.Link) != 1 {
		t.Errorf("Expected 1, but got %v", len(dot.CacheStore.Link))
	}
	// fmt.Println(dot.CacheStore.Link)
	*/
}

/* FIXME os.Exit
// TestHelp ...
func TestHelp(t *testing.T) {
	os.Args = []string{"dot", "--help"}
	_, err := cmd.Parse()
	if err != nil {
		t.Errorf("%v", err)
	}

	flagsErr, ok := err.(*flags.Error)
	if !ok {
		t.Errorf("Expected flag error, but got %v", err)
	} else if flagsErr.Type != flags.ErrHelp {
		t.Errorf("Expected ErrHelp, but got %v", flagsErr.Type)
	}
}
*/
