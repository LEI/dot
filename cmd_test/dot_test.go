package cmd_test

import (
	"os"
	"testing"

	"github.com/LEI/dot/cmd"
	// "github.com/jessevdk/go-flags"
)

/* https://joeshaw.org/testing-with-os-exec-and-testmain/
// TestMain ...
func TestMain(m *testing.M) {
	switch os.Getenv("GO_TEST_MODE") {
	case "":
		// Normal test mode
		os.Exit(m.Run())

	case "echo":
		iargs := []interface{}{}
		for _, s := range os.Args[1:] {
			iargs = append(iargs, s)
		}
		fmt.Println(iargs...)
	}
}

cmd := exec.Command(os.Args[0], "hello", "world")
cmd.Env = []string{"GO_TEST_MODE=echo"}
output, err := cmd.Output()
if err != nil {
	t.Errorf("echo: %v", err)
}
if g, e := string(output), "hello world\n"; g != e {
	t.Errorf("echo: want %q, got %q", e, g)
}
*/

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
