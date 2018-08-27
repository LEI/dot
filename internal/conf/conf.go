package conf

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"path/filepath"

	// Use "github.com/BurntSushi/toml" over
	// toml "github.com/pelletier/go-toml"
	// because it allows unmarshalling into
	// map[string]interface{}
	"github.com/BurntSushi/toml"
	yaml "gopkg.in/yaml.v2"
)

type configType struct {
	name   string
	alt    []string // Alternative extensions
	decode func([]byte, interface{}) error
}

var configFileTypes = []configType{
	{"toml", []string{}, toml.Unmarshal},
	{"yaml", []string{"yml"}, yaml.Unmarshal},
	{"json", []string{}, json.Unmarshal},
	/* {"ini", []string{}, func(b []byte, i interface{}) error {
		// if err := ini.MapTo(&i, b); err != nil {
		// 	return err
		// }
		return ini.MapTo(&i, b)
	}}, */
}

// ReadFile config file
func ReadFile(path string) (map[string]interface{}, error) {
	b, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}
	return Read(path, b)
}

// Read detects the config file type base on its extension or content.
func Read(path string, b []byte) (map[string]interface{}, error) {
	var data map[string]interface{}
	fileTypes := configFileTypes
FT:
	// Check file extension
	for _, ft := range configFileTypes {
		exts := append([]string{ft.name}, ft.alt...)
		for _, e := range exts {
			if e == filepath.Ext(path) {
				fileTypes = []configType{ft}
				break FT
			}
		}
	}
	// Attempt to decode config
	for i, ft := range fileTypes {
		err := ft.decode(b, &data)
		if err != nil {
			// Last or single file type
			if i == len(fileTypes)-1 {
				return data, fmt.Errorf("%s error: %s", ft.name, err)
			}
			// if Verbose > 1 {
			// 	fmt.Fprintf(os.Stderr, "failed to decode as %s: %s\n", ft.name, err)
			// }
			continue
		}
		// if err == nil {
		// 	fmt.Printf("%s: decoded as %s config file\n", path, ft.name)
		// }
		break
	}
	return data, nil // fmt.Errorf("%s: unknown config file type", path)
}
