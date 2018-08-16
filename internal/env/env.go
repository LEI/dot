package env

import (
	"os"
	"strings"
)

var (
	original map[string]string
)

func init() {
	original = GetAll()
}

// // Init ...
// func Init() error {
// 	for k, v := range baseEnv {
// 		if err := os.Setenv(k, v); err != nil {
// 			return err
// 		}
// 	}
// 	return nil
// }

// Set ...
func Set(k, v string) error {
	// if Verbose > 0 {
	// 	fmt.Printf("%s=%s\n", k, v)
	// }
	return os.Setenv(k, v)
}

// GetAll ...
func GetAll() map[string]string {
	env := make(map[string]string, 0)
	for _, i := range os.Environ() {
		sep := strings.Index(i, "=")
		k := i[0:sep]
		v := i[sep+1:]
		if _, ok := env[k]; !ok {
			env[k] = v
		}
	}
	return env
}

// Get ...
func Get(k string) string {
	env := GetAll()
	// v, ok := env[k]
	// if !ok {
	// 	return v
	// }
	v := env[k]
	// v, err := TemplateEnv(k, v)
	// if err != nil {
	// 	return err
	// }
	return v
}

// Expand ...
func Expand(k string, env map[string]string) string {
	return os.Expand(k, Get)
}

// ExpandEnv ...
func ExpandEnv(k string) string {
	return Expand(k, GetAll())
	// return os.ExpandEnv(k)
}

// Restore ...
func Restore(env map[string]string) error {
	os.Clearenv()
	for k, v := range env {
		if err := os.Setenv(k, v); err != nil {
			return err
		}
	}
	return nil
}
