package cmd

import (
	"fmt"
	"os"
)

func er(msg interface{}) {
	fmt.Fprintln(os.Stderr, "Error:",  msg)
	os.Exit(1)
}
