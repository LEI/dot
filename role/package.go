package role

import (
	"fmt"
	"os"
)

// type Package map[string]interface{}
type Package struct {
	Dir       *Dir
	Dirs      []*Dir
	Link      *Link   // interface{}
	Links     []*Link // interface{}
	Line      *Line
	Lines     []*Line
	Template  *Template
	Templates []*Template
	// Template interface{}
	PreInstall  string `mapstructure:"pre_install"`
	PostInstall string `mapstructure:"post_install"`
	PreRemove   string `mapstructure:"pre_remove"`
	PostRemove  string `mapstructure:"post_remove"`
	Env         map[string]string
}

func getMapInterfaceKey(val map[interface{}]interface{}, key string) interface{} {
	v, ok := val[interface{}(key)]
	if !ok {
		fmt.Fprintf(os.Stderr, "Missing key '%s' in %+v\n", key, val)
		os.Exit(1)
	}
	return v
}

func getMapStringKey(val map[string]interface{}, key string) interface{} {
	v, ok := val[key]
	if !ok {
		fmt.Fprintf(os.Stderr, "Missing key '%s' in %+v\n", key, val)
		os.Exit(1)
	}
	return v
}
