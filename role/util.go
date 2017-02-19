package role

import (
	"fmt"
	"os"
)

func fatal(msg interface{}) {
	fmt.Fprintf(os.Stderr, "Error while parsing link: %s", msg)
	os.Exit(64)
}

func fataln(msg interface{}) {
	fatal(fmt.Sprintf("%s\n", msg))
}
