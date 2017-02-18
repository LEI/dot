package role

import (
	"fmt"
	"github.com/LEI/dot/git"
	"github.com/spf13/viper"
	"os"
	"strings"
)

type Role struct {
	// Origin string
	Source string
	Target string
	Packages []*Package
}

func (role *Role) String() string {
	return fmt.Sprintf("%+v", *role)
}

type Package struct {
	Name   string
	Path   string
	Origin string
	Os     []string
	Repo   *git.Repository
	Config *viper.Viper

	Dir string
	Dirs []string
	Link  interface{} // *PackageLink
	Links []interface{}
	Lines map[string]string

}

// type PackageLink interface{}

func (pkg *Package) String() string {
	return fmt.Sprintf("%+v", *pkg)
}

// func (pkg *Package) Set(value string?) error {
// }

func (pkg *Package) CheckOsType(types []string) bool {
	if len(pkg.Os) == 0 {
		return true
	}
	for _, osType := range pkg.Os {
		for _, t := range types {
			if t == osType {
				return true
			}
		}
	}
	// fmt.Fprintf(os.Stderr, "[%s] %s: unsupported platform, only for %+v\n", pkg.Name, OS, pkg.Os)
	return false
}

func (pkg *Package) InitRepo(useHttps bool) error {
	git.Https = useHttps
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
	// PkgConfig[pkg.Name] = pkg.Config
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


// func NewPackage(spec string) (*Package, error) {
// 	pkg := &Package{
// 		Origin: spec,
// 	}
// 	err := pkg.InitRepo(git.Https)
// 	if err != nil {
// 		return nil, err
// 	}
// 	// repo, err := git.NewRepository(spec)
// 	// p := &Package{
// 	// 	Name: repo.Name,
// 	// 	Path: repo.Path,
// 	// 	Repo: repo,
// 	// }
// 	return pkg, nil
// }

// func (pkg *Package) SyncLink(link *File, target string) error {
// 	paths, err := filepath.Glob(path)
// 	// fmt.Printf("Find: %s -> %+v\n", link.Path, paths)
// 	if err != nil {
// 		return err
// 	}
// GLOB:
// 	for _, src := range paths {

// 		// file := NewFile(src)
// 		// file.Link(dst)
// 		// err = fileutil.Link(src, dst)
// 		// if err != nil {
// 		// 	return err
// 		// }
// 		// filepath.Rel(pkg.Path, dst)
// 	}
// 	return nil
// }

func (pkg *Package) UnmarshalKey(k string, v interface{}) interface{} {
	err := pkg.Config.UnmarshalKey(k, v)
	if err != nil {
		fmt.Fprintf(os.Stderr, "could not unmarshal %s: %s\n", k, err)
	}
	return v
}

// type PackageSlice []Package

// func (slice *PackageSlice) String() string {
// 	return fmt.Sprintf("%+v", *slice)
// }

// func (slice *PackageSlice) Type() string {
// 	return fmt.Sprintf("%T", *slice)
// }

// func (slice *PackageSlice) Set(value string) error {
// 	pkg, err := NewPackage(value)
// 	if err != nil {
// 		return err
// 	}
// 	*slice = append(*slice, *pkg)
// 	return nil
// }

// type PackageMap map[string]Package
// (*pkgMap)[p.Name] = *pkg
