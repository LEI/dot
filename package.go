package main

import (
	"fmt"
	"strings"
)

type Package struct {
	Name string
	Origin string
	Path string
	// Os OsType
}

type PackageSlice []Package

func (list *PackageSlice) String() string {
	return fmt.Sprintf("%+v", *list)
}

func (list *PackageSlice) Type() string {
	return fmt.Sprintf("%T", *list)
}

func (list *PackageSlice) Set(origin string) error {
	p := &Package{}
	if strings.Contains(origin, "=") {
		s := strings.Split(origin, "=")
		p.Name = s[0]
		p.Origin = s[1]
	} else {
		p.Name = origin
		p.Origin = origin
	}
	*list = append(*list, *p)
	// (*pkgMap)[p.Name] = *p
	return nil
}
