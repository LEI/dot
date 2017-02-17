package cli

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

// var interactive bool

func Confirm(str string) bool {
	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Printf("%s [y/n]: ", str)
		// if ForceYes {
		// 	fmt.Printf("%s", "Forced")
		// 	return true
		// }
		res, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		res = strings.ToLower(strings.TrimSpace(res))
		switch res {
		case "y", "yes":
			return true
		case "n", "no":
			return false
		}
		// TODO limit retries
	}
}
