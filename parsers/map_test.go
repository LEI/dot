package parsers

import (
	"reflect"
	"testing"
	// "github.com/LEI/dot/dotfile"
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

// func TestUnmarshalYAML(t *testing.T) {
// 	in := `---
// role:
//   # os: ["*"]
//   link:
//     - "a:b"` // a: "b"
// 	// out := &dotfile.Role{}
// 	// role := &dotfile.Role{}
// 	// _, err := role.ReadConfig(tmp)
// 	paths := make(Map)
// 	if err := dotfile.PreparePaths(paths); err != nil {
// 		t.Errorf(err)
// 	}
// }
