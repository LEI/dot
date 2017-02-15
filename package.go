package main

import (
	"fmt"
	// "path/filepath"
	// "strings"
	"github.com/LEI/dot/git"
)

type Package struct {
	Name   string
	Origin string
	Path   string
	// Os OsType
	Repo *git.Repository
}

type PackageSlice []Package

func (list *PackageSlice) String() string {
	return fmt.Sprintf("%+v", *list)
}

func (list *PackageSlice) Type() string {
	return fmt.Sprintf("%T", *list)
}

func (list *PackageSlice) Set(value string) error {
	p := &Package{}
	repo, err := git.NewRepository(value)
	if err != nil {
		return err
	}
	p.Repo = repo
	p.Name = repo.Name
	*list = append(*list, *p)
	// (*pkgMap)[p.Name] = *p
	return nil
}
