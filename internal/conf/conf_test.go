package conf

import (
	"reflect"
	"testing"
)

func TestRead(t *testing.T) {
	testCases := []struct {
		path string
		in   []byte
		out  map[string]interface{}
	}{
		{"test.json", []byte(`{"a": 1}`), map[string]interface{}{"a": float64(1)}},
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
}
