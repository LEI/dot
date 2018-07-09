package cfg

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v2"
)

var (
	// Debug ...
	Debug bool
)

// Load ...
func Load(c interface{}, s string) (string, error) {
	cfgPath, err := findConfig(s)
	if err != nil {
		return cfgPath, err
	}

	if cfgPath == "" {
	    return "", nil
	}

	cfg, err := readConfig(cfgPath)
	if err != nil {
	    fmt.Println("xx", err)
	    return cfgPath, err
	}

	err = yaml.Unmarshal(cfg, &c)
	if Debug {
	    fmt.Printf("+++\n%v\n+++\n", string(cfg))
	}
	if err != nil {
	    return cfgPath, err
	}
	if Debug {
	    fmt.Printf("---\n%v\n---\n", c)
	}

	return cfgPath, nil
}

func findConfig(s string) (string, error) {
	if s == "" {
	    return "", nil
	}

	paths := []string {
		s, // Current working directory
		filepath.Join(getConfigDir(), s),
	}

	for _, p := range paths {
		if isFile(p) {
			return p, nil
		}
	}

	return "", nil
}

// shibukawa/configdir
func getConfigDir() string {
	dir := ""
	if os.Getenv("XDG_CONFIG_HOME") != "" {
		dir = os.Getenv("XDG_CONFIG_HOME")
	} else {
		dir = filepath.Join(os.Getenv("HOME"), ".config")
	}
	// XDG_CONFIG_DIRS /etc/xdg
	return dir
}

func readConfig(s string) ([]byte, error) {
	bytes, err := read(s)
	// str := string(bytes)
	// if err != nil {
	// 	return str, err
	// }
	return bytes, err
}

// func exists(s string) bool {
// 	_, err := os.Stat(s)
// 	return !os.IsNotExist(err)
// }

func isFile(s string) bool {
	fi, err := os.Stat(s)

	return !os.IsNotExist(err) && !fi.IsDir()
}

func read(s string) ([]byte, error) {
	return ioutil.ReadFile(s)
}
