package dot

import (
	"fmt"
)

// RegisterLink ...
func RegisterLink(s, t string) error {
	CacheStore.Link[s] = t

	return nil
}

// Link ...
func Link(s, t string) error {
	source := s
	target := t

	fmt.Println("ln -s", source, target)

	return nil
}

// Unlink ...
func Unlink(s, t string) error {
	// source := s
	target := t

	fmt.Println("rm", target)

	return nil
}
