package cli

import (
	"fmt"
	"strings"
)

// FormatArgs formats command line aruments to a string.
func FormatArgs(args []string) string {
	for i, a := range args {
		if strings.Contains(a, " ") {
			args[i] = fmt.Sprintf("%q", a)
			// windows? args[i] = syscall.EscapeArg(a)
		}
		// switch v := a.(type) {
		// case string:
		// }
	}
	return strings.Join(args, " ")
}
