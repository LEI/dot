package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func confirm(str string) bool {
	reader := bufio.NewReader(os.Stdin)

	for {
		fmt.Printf("%s [y/n]: ", str)

		if ForceYes {
			fmt.Printf("%s", "Forced")
			return true
		}

		res, err := reader.ReadString('\n')
		if err != nil {
			handleError(err)
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
