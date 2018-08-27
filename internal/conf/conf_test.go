package conf

import (
	"path/filepath"
	"reflect"
	"testing"
)

func TestRead(t *testing.T) {
	testCases := []struct {
		path string
		in   []byte
		out  map[string]interface{}
	}{
		{
			"test.toml",
			[]byte(`a = "b"`),
			map[string]interface{}{
				"a": "b",
			},
		},
		{
			"test.yaml",
			[]byte(`a: b`),
			map[string]interface{}{
				"a": "b",
			},
		},
		{
			"test.json",
			[]byte(`{"a": "b"}`),
			map[string]interface{}{
				"a": "b",
			},
		},
	}
	for _, tc := range testCases {
		data, err := Read(tc.path, tc.in)
		if err != nil {
			t.Fatalf("could not read file: %s", err)
		}
		if !reflect.DeepEqual(data, tc.out) {
			t.Errorf("got %s; want %s", data, tc.out)
		}
	}
	// JSON actually gets decoded as YAML w/o extension
	for _, tc := range testCases {
		tc.path = removeExt(tc.path)
		data, err := Read(tc.path, tc.in)
		if err != nil {
			t.Fatalf("could not read file w/o extension: %s", err)
		}
		if !reflect.DeepEqual(data, tc.out) {
			t.Errorf("got %s; want %s", data, tc.out)
		}
	}
}

func removeExt(path string) string {
	return path[0 : len(path)-len(filepath.Ext(path))]
}
