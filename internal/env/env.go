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
		k, v := Split(i)
		if _, ok := env[k]; !ok {
			env[k] = v
		}
	}
	return env
}

// Split "key=value" into two variables.
func Split(s string) (string, string) {
	sep := strings.Index(s, "=")
	k := strings.ToUpper(s[0:sep])
	v := s[sep+1:]
	return k, v
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
func Expand(s string) string {
	// return os.ExpandEnv(s)
	return os.Expand(s, Get)
}

// ExpandEnv ...
func ExpandEnv(s string, env map[string]string) string {
	expand := func(k string) string {
		if v, ok := env[k]; ok {
			return v
		}
		return Get(k)
	}
	return os.Expand(s, expand)
}

// Lookup environment variable
func Lookup(k string) (string, bool) {
	return os.LookupEnv(k)
}

// Clear environment
func Clear() {
	os.Clearenv()
}

// Restore environment
func Restore(env map[string]string) error {
	Clear()
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
