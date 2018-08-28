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

// Set environment variable
func Set(k, v string) error {
	// if Verbose > 0 {
	// 	fmt.Printf("%s=%s\n", k, v)
	// }
	return os.Setenv(k, v)
}

// Unset environment variable
func Unset(k, v string) error {
	return os.Unsetenv(k)
}

// GetAll environment variables
func GetAll() map[string]string {
	env := make(map[string]string, 0)
	for _, i := range os.Environ() {
		sep := strings.Index(i, "=")
		k := strings.ToUpper(i[0:sep])
		v := i[sep+1:]
		if _, ok := env[k]; !ok {
			env[k] = v
		}
	}
	return env
}

// Get environment variable
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
func Expand(k string) string {
	return os.ExpandEnv(k)
	// return os.Expand(k, Get)
}

// ExpandEnv ...
func ExpandEnv(k string, env map[string]string) string {
	expand := func(k string) string {
		if v, ok := env[k]; ok {
			return v
		}
		return Get(k)
	}
	return os.Expand(k, expand)
}

// Lookup environment variable
func Lookup(k string) (string, bool) {
	return os.LookupEnv(k)
}

// Restore environment
func Restore(env map[string]string) error {
	os.Clearenv()
	for k, v := range env {
		if err := os.Setenv(k, v); err != nil {
			return err
		}
	}
	return nil
}

// RestoreOriginal environment
func RestoreOriginal(env map[string]string) error {
	return Restore(original)
}
