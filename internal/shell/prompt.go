package shell

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// AskConfirmation ...
func AskConfirmation(s string) (ret bool) {
	if noConfirm() {
		fmt.Println(s)
		return true
	}
	reader := bufio.NewReader(Stdin)
	for {
		// fmt.Printf("%s [y/n]: ", s)
		fmt.Printf("%s [y/n]:\n", s)
		res, err := reader.ReadString('\n')
		if err != nil {
			fmt.Fprintf(Stderr, "Could not read input from stdin: %s\n", err)
			os.Exit(1)
		}
		res = strings.ToLower(strings.TrimSpace(res))
		if res == "y" || res == "yes" {
			ret = true
			break
		} else if res == "n" || res == "no" {
			ret = false
			break
		}
	}
	// FIXME: no new line if enter is pressed before the last fmt.Printf
	// fmt.Printf("\n")
	return
}

func noConfirm() bool {
	ncfile := filepath.Join(HomeDir, ".dotnc")
	_, err := os.Stat(ncfile)
	exists := err == nil || os.IsExist(err)
	if exists {
		fmt.Fprintln(Stderr, "[Confirmation disabled because ~/.dotnc exists]")
	}
	return exists
}
