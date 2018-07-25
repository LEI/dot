package parsers

import (
	"reflect"
	"testing"
)

var sliceTests = []struct {
	in  interface{}
	out []string
}{
	{"a", []string{"a"}},
	{[]string{"a"}, []string{"a"}},
}

func TestSlice(t *testing.T) {
	for _, tt := range sliceTests {
		slice, err := NewSlice(tt.in)
		if err != nil {
			t.Error(err)
		}
		if !reflect.DeepEqual(slice.Value(), tt.out) {
			t.Errorf("in: %+v out: %+v expected: %+v", tt.in, slice.Value(), tt.out)
		}
	}
}
