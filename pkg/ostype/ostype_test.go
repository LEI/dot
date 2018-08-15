package ostype

import (
	// "fmt"
	// "gopkg.in/go-ini/ini.v1"
	// "os"
	// "os/exec"
	// "path/filepath"
	// "regexp"
	// "runtime"
	// "strconv"
	// "strings"
	"testing"
)

var osTests = []struct {
	os  string
	in  []string
	out bool
}{
	{"a", []string{}, false},
	{"a", []string{"a"}, true},
	{"a", []string{"a", "b"}, true},
	{"a", []string{"b", "a"}, true},
	{"a", []string{"b"}, false},
	{"a", []string{"!a"}, false},
	{"a", []string{"!a"}, false},
	{"a", []string{"!b"}, true},
	{"a", []string{"a", "!b"}, true},
	{"a", []string{"!a", "b"}, false},
	{"a", []string{"!a", "!b"}, false},
	{"a", []string{"!b", "!a"}, false},
	{"a", []string{"!b", "!c"}, true},
}

func TestHas(t *testing.T) {
	for _, tt := range osTests {
		List = []string{tt.os}
		matched := Has(tt.in...)
		if matched != tt.out {
			t.Fatalf("%s should be %v with %+v", tt.os, tt.out, tt.in)
		}
	}
}
