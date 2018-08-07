package tasks

import (
	"fmt"
	"strings"
)

func parseDest(s string) (src, dst string, err error) {
	src = s
	if strings.Contains(s, ":") {
		parts := strings.Split(s, ":")
		if len(parts) != 2 {
			return src, dst, fmt.Errorf("unable to parse dest spec: %s", s)
		}
		src = parts[0]
		dst = parts[1]
	}
	// if dst == "" && isDir(src) {
	// 	dst = PathHead(src)
	// }
	return src, dst, nil
}
