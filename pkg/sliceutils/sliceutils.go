package sliceutils

import (
	"fmt"
	"os"
	"regexp"
)

// Matches ...
func Matches(in []string, list []string) bool {
	for _, pattern := range in {
		negated := pattern[0] == '!'
		if negated {
			pattern = pattern[1:]
		}
		for _, str := range list {
			matched, err := regexp.MatchString(pattern, str)
			if err != nil {
				fmt.Fprintf(os.Stderr, "%s: %s\n", pattern, err)
				os.Exit(1)
			}
			if negated && matched {
				return false
			}
			if matched {
				return true
			}
		}
		if negated {
			return true
		}
	}
	return false
}

// Contains check if a slice contains a given string
func Contains(in []string, s string) bool {
	for _, a := range in {
		if a == s {
			return true
		}
	}
	return false
}

// Intersects check if a slice contains at least one element from the other
func Intersects(in []string, list []string) bool {
	for _, a := range in {
		if Contains(list, a) {
			return true
		}
	}
	return false
}
