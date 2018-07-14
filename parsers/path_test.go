package parsers

import (
	"fmt"
	"os"
	"testing"

	// "github.com/jessevdk/go-flags"

	// "github.com/LEI/dot/cmd"
)

var pathtests = []struct {
	in  string
	out string
}{
	{"$HOME", "$HOME"},
}

// TestCopy ...
func TestPathsAdd(t *testing.T) {
	for _, pt := range pathtests {
		paths := &Paths{}
		// for _, p := range paths {}
		paths.Add(pt.in)
		res := os.ExpandEnv(pt.out)
		for src, dst := range *paths {
			fmt.Printf("in:%s src:%s dst:%s want:%s\n", pt.in, src, dst, res)
			if src != res {
				t.Errorf("got %+v, want %+v", src, res)
			}
		}
	}

}
