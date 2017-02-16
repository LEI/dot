package cmd

import (
	"fmt"
	// "github.com/LEI/dot/role"
	// "github.com/spf13/viper"
	"os"
)

func fatal(msg interface{}) {
	// log.Fatal*
	fmt.Fprintln(os.Stderr, "Error:", msg)
	os.Exit(1)
}
