package tasks

import (
	"fmt"
	"strings"

	"github.com/LEI/dot/cli/config/types"
	"github.com/LEI/dot/system"
	"github.com/mitchellh/mapstructure"
)

var (
	defaultTemplateExt = "tpl"
)

// Template task
type Template struct {
	Task
	Source, Target string
	Ext            string    // `default:"tpl"`
	Env            types.Map // map[string]string
	// Vars           map[string]interface{}
	Vars        types.Map
	IncludeVars string `mapstructure:"include_vars"`
	// backup bool
	// overwrite bool
}

func (t *Template) String() string {
	return fmt.Sprintf("template[%s:%s]", t.Source, t.Target)
}

// Data template
func (t *Template) Data() (map[string]interface{}, error) {
	data := make(map[string]interface{}, 0) // types.Map{}
	for k, v := range t.Env {               // GetEnv()
		data[k] = v
	}
	// for k, v := range t.Env {
	// 	// k = strings.ToTitle(k)
	// 	v, err := TemplateEnv(k, v)
	// 	if err != nil {
	// 		return err
	// 	}
	// 	// fmt.Println("$ ENV", k, "=", v)
	// 	data[k] = v
	// }
	if t.IncludeVars != "" {
		inclVars, err := system.ParseTemplate(t.IncludeVars)
		if err != nil {
			return data, err
		}
		for k, v := range inclVars {
			// if w, ok := t.Vars[k]; ok { return fmt.Errorf... }
			data[k] = v
		}
	}
	for k, v := range t.Vars {
		data[k] = v
	}
	return data, nil
}

// Check template task
func (t *Template) Check() error {
	if t.Source == "" {
		return fmt.Errorf("template: empty source")
	}
	data, err := t.Data()
	if err != nil {
		return err
	}
	err = system.CheckTemplate(t.Source, t.Target, data)
	switch err {
	case system.ErrTemplateAlreadyExist:
		t.Done()
	default:
		return err
	}
	return nil
}

// Install template task
func (t *Template) Install() error {
	str := fmt.Sprintf("gotpl %s %s", t.Source, t.Target)
	if !t.ShouldInstall() {
		if Verbose > 0 {
			fmt.Fprintf(Stdout, "# %s\n", str)
		}
		return ErrSkip
	}
	fmt.Fprintf(Stdout, "$ %s\n", str)
	data, err := t.Data()
	if err != nil {
		return err
	}
	// fmt.Println("tpl data", data)
	return system.CreateTemplate(t.Source, t.Target, data)
}

// Remove template task
func (t *Template) Remove() error {
	str := fmt.Sprintf("rm %s", t.Target)
	if !t.ShouldRemove() {
		if Verbose > 0 {
			fmt.Fprintf(Stdout, "# %s\n", str)
		}
		return ErrSkip
	}
	fmt.Fprintf(Stdout, "$ %s\n", str)
	data, err := t.Data()
	if err != nil {
		return err
	}
	// fmt.Println("tpl data", data)
	return system.RemoveTemplate(t.Source, t.Target, data)
}

// Templates task slice
type Templates []*Template

func (templates *Templates) String() string {
	// s := ""
	// for i, t := range *templates {
	// 	s += fmt.Sprintf("%s", t)
	// 	if i > 0 {
	// 		s += "\n"
	// 	}
	// }
	// return s
	return fmt.Sprintf("%s", *templates)
}

// Parse template tasks
func (templates *Templates) Parse(i interface{}) error {
	tt := &Templates{}
	m, err := types.NewMap(i, "source")
	if err != nil {
		return err
	}
	for k, v := range *m {
		t := &Template{}
		switch val := v.(type) {
		case string:
			t.Source = k
			t.Target = val
		case *types.Map:
			mapstructure.Decode(val, &t)
		case map[interface{}]interface{}:
			mapstructure.Decode(val, &t)
		case interface{}:
			t = val.(*Template)
		default:
			return fmt.Errorf("invalid template map: %+v", val)
		}
		tt.Add(*t)
		// *tt = append(*tt, t)
	}
	*templates = *tt
	return nil
}

// Add a template
func (templates *Templates) Add(t Template) {
	if t.Ext == "" {
		t.Ext = defaultTemplateExt
	}
	if t.Target != "" && t.Ext != "" && strings.HasSuffix(t.Target, "."+t.Ext) {
		t.Target = strings.TrimSuffix(t.Target, "."+t.Ext)
	}
	*templates = append(*templates, &t)
}
