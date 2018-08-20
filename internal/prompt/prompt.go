package prompt

// https://github.com/manifoldco/promptui
// https://github.com/c-bata/go-prompt
// https://github.com/AlecAivazis/survey

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/docker/docker/pkg/homedir"
)

// AskConfirmation ...
func AskConfirmation(s string) (ret bool) {
	if noConfirm() {
		return true
	}
	reader := bufio.NewReader(os.Stdin)
	for {
		// fmt.Printf("%s [y/n]: ", s)
		fmt.Printf("%s [y/n]:\n", s)
		res, err := reader.ReadString('\n')
		if err != nil {
			fmt.Fprintf(os.Stderr, "Could not read input from stdin: %s\n", err)
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
	ncfile := filepath.Join(homedir.Get(), ".dotnc")
	_, err := os.Stat(ncfile)
	exists := err == nil || os.IsExist(err)
	if exists {
		fmt.Fprintln(os.Stderr, "(Confirmation disabled because ~/.dotnc exists)")
	}
	return exists
}
