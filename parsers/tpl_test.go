package parsers

import (
	"testing"
	"reflect"
)

var tplTests = []struct{
	in interface{}
	out *Tpl
}{
	{"a", &Tpl{Source: "a"}},
}

func TestTpl(t *testing.T) {
	for _, tt := range tplTests {
		p, err := NewTpl(tt.in)
		if err != nil {
			t.Error(err)
		}
		if !reflect.DeepEqual(p, tt.out) {
			t.Errorf("in: %+v out: %+v expected: %+v", tt.in, p, tt.out)
		}
	}
}
