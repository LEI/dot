package parsers

import (
	"reflect"
	"testing"
)

var mapTests = []struct {
	in struct {
		key, val string
	}
	out map[string]string
}{
	{struct{ key, val string }{"a", ""}, map[string]string{"a": ""}},
	// {[]string{"a"}, map[string]string{"a": ""}},
}

func TestMap(t *testing.T) {
	for _, tt := range mapTests {
		m := &Map{}
		m.Add(tt.in.key, tt.in.val)
		if !reflect.DeepEqual(m.Value(), tt.out) {
			t.Errorf("in: %+v out: %+v expected: %+v", tt.in, m.Value(), tt.out)
		}
	}
}
