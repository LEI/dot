package role

import (
	"fmt"
	// "path/filepath"
	// "strings"
	"github.com/LEI/dot/git"
)

type Package struct {
	Name   string
	Path   string
	Origin string
	Repo   *git.Repository
	// Os OsType
}

// func (*pkg Package) Set(value string)

func NewPackage(s string) (*Package, error) {
	repo, err := git.NewRepository(s)
	if err != nil {
		return nil, err
	}
	p := &Package{
		Name: repo.Name,
		Path: repo.Path,
		Repo: repo,
	}
	return p, nil
}

type PackageSlice []Package

func (slice *PackageSlice) String() string {
	return fmt.Sprintf("%+v", *slice)
}

func (slice *PackageSlice) Type() string {
	return fmt.Sprintf("%T", *slice)
}

func (slice *PackageSlice) Set(value string) error {
	pkg, err := NewPackage(value)
	if err != nil {
		return err
	}
	*slice = append(*slice, *pkg)
	// (*pkgMap)[p.Name] = *pkg
	return nil
}
