package env

import (
	"path/filepath"
	"testing"
)

func TestExpandEnvVar(t *testing.T) {
	testCases := []struct {
		in  string
		out string
		env map[string]string
	}{
		{"a", "a", nil},
		{"$a", "b", map[string]string{"a": "b"}},
		{"${a:-b}", "b", nil},
		{"${a:-$(echo c)}", "c", nil},
		{"$(echo $a)", "b", map[string]string{"a": "b"}},
	}
	for _, tc := range testCases {
		v := ExpandEnvVar("", tc.in, tc.env)
		if v != tc.out {
			t.Errorf("got %s; want %s", v, tc.out)
		}
	}
}

func removeExt(path string) string {
	// return strings.TrimSuffix(path, filepath.Ext(path))
	return path[0 : len(path)-len(filepath.Ext(path))]
}
