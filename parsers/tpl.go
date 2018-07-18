package parsers

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"reflect"
	"strings"

	"gopkg.in/yaml.v2"
)

// Tpl ...
type Tpl struct {
	Source, Target string
	Ext            string `default:"tpl"`
	Env            map[string]string
	Vars           map[string]interface{}
	IncludeVars    string `yaml:"include_vars"`
	// Data           interface{}
}

// Templates ...
type Templates []*Tpl

// func (p *Templates) String() string {
// 	s := ""
// 	for _, v := range *p {
// 		s+= fmt.Sprintf("%+v", v)
// 	}
// 	return s
// }

// Append ...
func (t *Templates) Append(tpl *Tpl) *Templates {
	if tpl.Target != "" && tpl.Ext != "" && strings.HasSuffix(tpl.Target, "."+tpl.Ext) {
		tpl.Target = strings.TrimSuffix(tpl.Target, "."+tpl.Ext)
	}
	*t = append(*t, tpl)
	return t
}

// Add ...
func (t *Templates) Add(i interface{}) error {
	tpl := &Tpl{}
	if i == nil {
		return fmt.Errorf("Trying to add nil to tmpls: %+v", t)
	}
	if val, ok := i.(string); ok {
		tpl.Source = val
	} else if val, ok := i.(Tpl); ok {
		*tpl = val
	} else if val, ok := i.(map[interface{}]interface{}); ok {
		// Get source
		src, ok := val["source"].(string)
		if !ok {
			return fmt.Errorf("Missing tpl source: %+v", val)
		}
		tpl.Source = src
		dst, ok := val["target"].(string)
		if ok {
			tpl.Target = dst
		}
		if ext, ok := val["ext"].(string); ok {
			tpl.Ext = ext
		}
		if env, ok := val["env"].(map[string]string); ok {
			// tpl.Env = NewSlice(env.(*Slice))
			tpl.Env = env
		} else {
			tpl.Env = make(map[string]string, 0)
		}
		if vars, ok := val["vars"].(map[string]interface{}); ok {
			tpl.Vars = vars
			// tpl.Data = data
		} else {
			tpl.Vars = make(map[string]interface{}, 0)
		}
		// Included variables override existing tpl.Vars keys
		if file, ok := val["include_vars"].(string); ok {
			// tpl.IncludeVars = file // os.ExpandEnv(file)
			inclVars, err := parseTemplate(file)
			if err != nil {
				return err
			}
			for k, v := range inclVars {
				// if w, ok := tpl.Vars[k]; ok { return fmt.Errorf... }
				tpl.Vars[k] = v
			}
		}
		// } else if val, ok := i.(*Tpl); ok {
		// 	tpl = val
		// } else if val, ok := i.([]string); ok {
		// 	fmt.Println("MS", val)
		// } else if val, ok := i.([]interface{}); ok {
		// 	// tpl.OS = *NewSlice(val["os"])
		// 	fmt.Println("IS", val)
		// } else if val, ok := i.(map[string]string); ok {
		// 	fmt.Println("MSS", val, i)
		// } else if val, ok := i.(map[string]interface{}); ok {
		// 	fmt.Println("MSI", val, i)
		// } else if val, ok := i.(interface{}); ok {
		// 	fmt.Println("II", val, i)
	} else {
		return fmt.Errorf("Unable to assert Tpl: %+v", i)
	}
	*t = append(*t, tpl)
	return nil
}

// UnmarshalYAML ...
func (t *Templates) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var i interface{}
	if err := unmarshal(&i); err != nil {
		return err
	}
	switch val := i.(type) {
	case map[string]string:
		fmt.Println("Unmarshal tpl -> map[string]string", val)
		for k, v := range val {
			if k != "" {
				return fmt.Errorf("Unexpected key: %s", k)
			}
			if err := t.Add(v); err != nil {
				return err
			}
		}
	case []interface{}:
		for _, v := range val {
			if err := t.Add(v); err != nil {
				return err
			}
		}
	default:
		t := reflect.TypeOf(val)
		T := t.Elem()
		if t.Kind() == reflect.Map {
			T = reflect.MapOf(t.Key(), t.Elem())
		}
		return fmt.Errorf("Unable to unmarshal templates (%s) into struct: %+v", T, val)
	}
	return nil
}

func parseTemplate(file string) (vars map[string]interface{}, err error) {
	if strings.HasPrefix(file, "~/") {
		file = filepath.Join(os.Getenv("HOME"), file[2:])
	}
	bytes, err := ioutil.ReadFile(file)
	if err != nil {
		return vars, err
	}
	if err := yaml.Unmarshal(bytes, &vars); err != nil {
		return vars, err
	}
	return vars, nil
}
