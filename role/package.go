package role

import (
	"fmt"
	// "path/filepath"
	// "strings"
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
	Os     []string
	Config *viper.Viper

	Dir string
	Dirs []string
	Link interface{}
	Links []interface{}
	Lines map[string]string // *Lines
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

func (pkg *Package) InitRepo() error {
	if pkg.Origin == "" {
		return fmt.Errorf("Empty origin in %+v\n", pkg)
	}
	if pkg.Repo != nil {
		fmt.Println("INFO: repo is already registered, skip init")
		return nil
	}
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
	return nil
}
func (pkg *Package) InitConfig(name string) error {
	if pkg.Config == nil {
		pkg.Config = viper.New()
	}
	pkg.Config.SetConfigName(name)
	pkg.Config.AddConfigPath(pkg.Path)
	err := pkg.Config.ReadInConfig()
	if err != nil {
		return err
	}
	if pkg.Config.ConfigFileUsed() == "" {
		fmt.Println("No config file, skip package init config")
		return nil
	}

	fmt.Println(pkg.Name, "using", pkg.Config.ConfigFileUsed())
	// PkgConfig[pkg.Name] = pkg.Config
	return nil
}

func (pkg *Package) Sync(source string, target string) error {
	if pkg.Config != nil {
		for _, dir := range pkg.GetDirs() {
			fmt.Printf("Create: %s/%s\n", target, dir)
		}
		for _, link := range pkg.GetLinks() {
			fmt.Printf("Link: %s/%s into %s\n", source, link, target)
		}
		for file, line := range pkg.GetLines() {
			fmt.Printf("Line: %s in %s\n", line, file)
		}
	} else {
		fmt.Println("-> NIL")
	}
	return nil
}

func (pkg *Package) GetDirs() []string {
	pkg.Dir = pkg.Config.GetString("dir")
	if pkg.Dir != "" {
		pkg.Dirs = append(pkg.Config.GetStringSlice("dirs"), pkg.Dir)
		pkg.Config.Set("dirs", pkg.Dirs)
	}
	return pkg.Config.GetStringSlice("dirs")
}

func (pkg *Package) GetLinks() []interface{} {
	// for _, key := range []string{"link", "links"} { }
	pkg.UnmarshalKey("link", &pkg.Link)
	pkg.UnmarshalKey("links", &pkg.Links)
	// link := pkg.Config.Get("link")
	// links := pkg.Config.Get("links")
	// for i, link := range links { links[i] = link.(*Link) }
	// if links == nil {
	// 	links = make([]*Link, 0)
	// }
	if pkg.Link != nil {
		pkg.Links = append(pkg.Links, pkg.Link)
		// links = append(links.([]*Link), link.(*Link))
		// pkg.Config.Set("links", append(pkg.Config.GetStringSlice("links"), link))
		// pkg.Config.Set("links", append(pkg.Config.Get("links").(role.Link), link))
		pkg.Config.Set("links", pkg.Links)
		// pkg.Config.Set("link", nil)
	}
	// fmt.Printf("links (%T): %+v\n", links)
	return pkg.Links
}

func (pkg *Package) GetLines() map[string]string {
	pkg.UnmarshalKey("lines", &pkg.Lines)
	// if pkg.Lines == nil {
	// 	pkg.Lines = make(Lines, 0)
	// }
	// fmt.Printf("lines (%T): %+v\n", lines)
	return pkg.Lines
}

func (pkg *Package) UnmarshalKey(k string, v interface{}) interface{} {
	err := pkg.Config.UnmarshalKey(k, v)
	if err != nil {
		fmt.Fprintf(os.Stderr, "could not unmarshal %s: %s\n", k, err)
	}
	return v
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
