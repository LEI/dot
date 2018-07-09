package cmd_test

import (
	"os"
	"testing"

	"github.com/LEI/dot/cmd"

	// "github.com/jessevdk/go-flags"
)

// TestParse DotCmd
func TestParse(t *testing.T) {
	os.Args = []string{"dot", "-V"}
	remaining, err := cmd.Parse()
	if err != nil {
		t.Errorf("Parse error: %v", err)
	}
	if cmd.Options.Version != true {
		t.Errorf("Expected Version flag, but got %v", cmd.Options.Version)
	}
	if len(remaining) != 0 {
		t.Errorf("Expected no remaining arguments, but got %v", remaining)
	}
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
