package role

import (
	"fmt"
	// "path/filepath"
	// "strings"
	"github.com/LEI/dot/cmd"
	"github.com/LEI/dot/git"
	"github.com/spf13/viper"
	"os"
	"strings"
)

type Package struct {
	Name   string
	Path   string
	Origin string
	Repo   *git.Repository
	Config *viper.Viper
	Os     []string

	Dir string
	Dirs []string

	Link interface{}
	Links []interface{}

	Lines map[string]string
}

func NewPackage(spec string) (*Package, error) {
	pkg := &Package{
		Origin: spec,
	}
	err := pkg.InitRepo()
	if err != nil {
		return nil, err
	}
	// repo, err := git.NewRepository(spec)
	// p := &Package{
	// 	Name: repo.Name,
	// 	Path: repo.Path,
	// 	Repo: repo,
	// }
	return pkg, nil
}

func (pkg *Package) String() string {
	return fmt.Sprintf("%+v", *pkg)
}

// func (pkg *Package) Set(value string?) error {
// }

func (pkg *Package) InitRepo() error {
	if pkg.Origin != "" && pkg.Repo == nil {
		spec := pkg.Origin
		if pkg.Name != "" && !strings.Contains(pkg.Origin, "=") {
			spec = pkg.Name + "=" + pkg.Origin
		}
		repo, err := git.NewRepository(spec)
		if err != nil {
			return err
		}
		pkg.Repo = repo
		pkg.Path = repo.Path
		// if pkg.Name == "" {
		// 	pkg.Name = repo.Name
		// }
	}
	return nil
}

func (pkg *Package) Check(currentOs string) bool {
	for _, osType := range pkg.Os {
		switch osType {
		case currentOs:
			return true
		default:
			// if Debug {
			// 	fmt.Fprintf(os.Stderr,
			// 		"[%s] %s: unsupported platform, only for %+v\n",
			// 		pkg.Name, OS, pkg.Os)
			// }
			return false
		}
	}
	return true
}

func (pkg *Package) Sync(source string, target string) error {
	fmt.Println(source, "->", target)
	err := cmd.InitPackage(pkg)
	if err != nil {
		return err
	}
	if pkg.Config != nil {
		fmt.Println("-> DIRS", pkg.Config.GetStringSlice("dirs"))
		fmt.Println("-> LINKS", pkg.Config.Get("links"))
		fmt.Println("-> LINES", pkg.Config.GetStringMapString("lines"))
	} else {
		fmt.Println("-> NIL")
	}
	return nil
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
	return nil
}
// type PackageMap map[string]Package
// (*pkgMap)[p.Name] = *pkg
