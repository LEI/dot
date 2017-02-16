package cmd

import (
	"fmt"
	"github.com/LEI/dot/role"
	"github.com/spf13/viper"
	"os"
)

var (
	Packages []role.Package

	Source = viper.GetString("source")
	Target = viper.GetString("target")

	Dir = viper.GetString("dir")
	Dirs = viper.GetStringSlice("dirs")

	Link = viper.GetString("link")
	Links = viper.GetString("links")

	Line = viper.GetString("line")
	Lines = viper.GetString("lines")
)

func init() {
	err := viper.UnmarshalKey("packages", &Packages)
	if err != nil {
		er(err)
	}
}

func er(msg interface{}) {
	fmt.Fprintln(os.Stderr, "Error:",  msg)
	os.Exit(1)
}
