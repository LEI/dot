package configfile

import (
	// "fmt"
	// "encoding/base64"
	"encoding/json"
	"io"
	// "io/ioutil"
	// "os"
	// "path/filepath"
	// "strings"
)

const (
)

// ConfigFile ~/.dotrc.yml file info
type ConfigFile struct {
	// AuthConfigs          map[string]types.AuthConfig `json:"auths"`
	// HTTPHeaders          map[string]string           `json:"HttpHeaders,omitempty"`
	// PsFormat             string                      `json:"psFormat,omitempty"`
	Filename             string                      `json:"-"` // Note: for internal use only
}

// LoadFromReader reads the configuration data given
func (configFile *ConfigFile) LoadFromReader(configData io.Reader) error {
	if err := json.NewDecoder(configData).Decode(&configFile); err != nil {
		return err
	}
	// var err error
	// for addr, ac := range configFile.AuthConfigs {
	// 	ac.Username, ac.Password, err = decodeAuth(ac.Auth)
	// 	if err != nil {
	// 		return err
	// 	}
	// 	ac.Auth = ""
	// 	ac.ServerAddress = addr
	// 	configFile.AuthConfigs[addr] = ac
	// }
	return nil
}

// func (configFile *ConfigFile) SaveToWriter(writer io.Writer) error {}
// func (configFile *ConfigFile) Save() error {}
