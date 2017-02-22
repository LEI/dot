package prompt

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

// var interactive bool

func Confirm(format string, v ...interface{}) bool {
	msg := fmt.Sprintf(format, v...)
	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Printf("%s [y/n]: ", msg)
		// TODO Force / AssumeYes
		res, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		res = strings.ToLower(strings.TrimSpace(res))
		if res != "" {
			fmt.Println()
		}
		switch res {
		case "y", "yes":
			return true
		case "n", "no":
			return false
		}
		// TODO limit retries
	}
}
