package dot

import (
	"log"
	"strings"

	"github.com/LEI/dot/internal/env"
)

// Env map
type Env map[string]string

// NewEnv vars
func NewEnv(i interface{}) *Env {
	e := &Env{}
	switch in := i.(type) {
	case string:
		k, v := env.Split(in)
		(*e)[k] = v
	case []string:
		for _, n := range in {
			k, v := env.Split(n)
			(*e)[k] = v
		}
	default:
		log.Fatalf("invalid env: %s\n", i)
	}
	return e
}

func parseEnviron(e *Env) (*Env, error) {
	m := &Env{}
	for k, v := range *e {
		k = strings.ToUpper(k)
		ev, err := buildTplEnv(k, v, *e)
		if err != nil {
			return m, err
		}
		// fmt.Printf("$ export %s=%q\n", k, ev)
		(*m)[k] = ev
	}
	return m, nil
}

// Vars map
type Vars map[string]interface{}

func parseVars(e *Env, vars *Vars, incl ...string) (*Vars, error) {
	data := &Vars{}
	// Parse extra variables, already merged with role vars
	for k, v := range *vars {
		// if k == "Env" ...
		if val, ok := v.(string); ok && val != "" {
			// Parse go template
			ev, err := buildTplEnv(k, val, *e)
			if err != nil {
				return data, err
			}
			// Expand resulting environment variables
			v = env.ExpandEnvVar(k, ev, *e)
			// expand := func(s string) string {
			// 	if v, ok := e[s]; ok {
			// 		return v
			// 	}
			// 	return env.Get(s) // os.ExpandEnv(s)
			// }
			// v = os.Expand(ev, expand)
		}
		// fmt.Printf("# var %s = %+v\n", k, v)
		(*data)[k] = v
	}
	// Included variables override existing vars
	for _, v := range incl {
		inclVars, err := includeVars(v) // os.ExpandEnv?
		if err != nil {
			return data, err
		}
		for k, v := range inclVars {
			(*data)[k] = v
		}
	}
	return data, nil
}
