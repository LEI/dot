package dotfile

import (
	// "fmt"
	// "os"
	"testing"
	// "github.com/jessevdk/go-flags"
	// "github.com/LEI/dot/cmd"
)

// TestCopy ...
func TestCopy(t *testing.T) {
	task := &CopyTask{
		Source: "",
		Target: "",
	}
	task.Do("install")
	task.Install()
}
