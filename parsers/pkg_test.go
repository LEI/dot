package parsers

import (
	"testing"
	"reflect"
)

var pkgTests = []struct{
	in interface{}
	out *Pkg
}{
	{"a", &Pkg{Name: "a"}},
}

func TestPkg(t *testing.T) {
	for _, tt := range pkgTests {
		p, err := NewPkg(tt.in)
		if err != nil {
			t.Error(err)
		}
		if !reflect.DeepEqual(p, tt.out) {
			t.Errorf("in: %+v out: %+v expected: %+v", tt.in, p, tt.out)
		}
	}
}
