package types

import (
	"fmt"
	"github.com/LEI/dot/pkg/ostype"
)

// HasOS ...
type HasOS struct {
	OS Slice
}

// CheckOS ...
func (h *HasOS) CheckOS() bool {
	if len(h.OS) == 0 {
		return true
	}
	hasOS := false
	for _, o := range h.OS {
		if ostype.Has(o) {
			hasOS = true
			break
		}
	}
	return hasOS
	// return ostype.Has(h.OS...)
}

// HasIf ...
type HasIf struct {
	If Slice
}

// CheckIf ...
func (h *HasIf) CheckIf() bool {
	if len(h.If) == 0 {
		return true
	}
	// https://golang.org/pkg/text/template/#hdr-Functions
	fmt.Println("TODO check if", h.If)
	return false
}
