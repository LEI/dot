package parsers

import (
	"testing"
	"reflect"
)

var sliceTests = []struct{
	in interface{}
	out []string
}{
	{"a", []string{"a"}},
	{[]string{"a"}, []string{"a"}},
}

func TestSlice(t *testing.T) {
	for _, tt := range sliceTests {
		slice := NewSlice(tt.in)
		if !reflect.DeepEqual(slice.Value(), tt.out) {
			t.Errorf("in: %+v out: %+v expected: %+v", tt.in, slice.Value(), tt.out)
		}
	}
}
