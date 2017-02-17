package role

import (
	"fmt"
	"path/filepath"
	// "strings"
	"github.com/LEI/dot/fileutil"
	"github.com/LEI/dot/git"
	"github.com/spf13/viper"
	"os"
	"runtime"
	"strings"
)

const OS = runtime.GOOS

func init() {
	err := os.Setenv("OS", OS)
	if err != nil {
		fmt.Printf("Could not set env OS=%s: %s", OS, err)
	}
}

var Ignore = []string{".git", ".*\\.md"}

type Package struct {
	Name   string
	Path   string
	Origin string
	Repo   *git.Repository
	Os     []string
	Config *viper.Viper

	Dir   string
	Dirs  []string
	Link  interface{}
	Links []interface{}
	Lines map[string]string // *Lines
}

// type Role interface {
// 	Get() interface{}
// 	Check() bool
// 	Sync(string, string) error
// }

type Link struct {
	Type string
	Path string
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

func (pkg *Package) CheckOS() bool {
	for _, osType := range pkg.Os {
		switch osType {
		case OS:
			return true
		default:
			// fmt.Fprintf(os.Stderr,
			// 	"[%s] %s: unsupported platform, only for %+v\n",
			// 	pkg.Name, OS, pkg.Os)
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
	// PkgConfig[pkg.Name] = pkg.Config
	return nil
}

func (pkg *Package) Sync(source string, target string) error {
	if pkg.Config == nil {
		return fmt.Errorf("%s: no config", pkg.Name)
	}
	err := pkg.SyncDirs(target)
	if err != nil {
		return err
	}
	err = pkg.SyncLinks(source, target)
	if err != nil {
		return err
	}
	err = pkg.SyncLines(target)
	if err != nil {
		return err
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

func (pkg *Package) SyncDirs(target string) error {
	for _, dir := range pkg.GetDirs() {
		dir = os.ExpandEnv(dir)
		dir = filepath.Join(target, dir)
		err := fileutil.MakeDir(dir)
		if err != nil {
			return err
		}
		fmt.Printf("Directory: %s\n", dir)
	}
	return nil
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

func (pkg *Package) SyncLinks(source string, target string) error {
	for _, link := range pkg.GetLinks() {
		var src *Link
		switch v := link.(type) {
		case string:
			src = &Link{Type: "", Path: v}
		case map[string]interface{}:
			src = &Link{Type: v["type"].(string), Path: v["path"].(string)}
		default:
			fmt.Fprintf(os.Stderr, "Unknown type %T for %+v, skipping link\n", v, v)
			continue
		}
		src.Path = os.ExpandEnv(src.Path)
		err := pkg.SyncLink(src, target)
		if err != nil {
			return err
		}
	}
	return nil
}

func (pkg *Package) SyncLink(link *Link, target string) error {
	path := filepath.Join(pkg.Path, link.Path)
	paths, err := filepath.Glob(path)
	// fmt.Printf("Find: %s -> %+v\n", link.Path, paths)
	if err != nil {
		return err
	}
	if len(paths) == 0 {
		fmt.Fprintf(os.Stderr, "%s: No such file or directory\n", path)
	}
GLOB:
	for _, src := range paths {
		for _, pattern := range Ignore {
			matched, err := filepath.Match(pattern, filepath.Base(src))
			if err != nil {
				return err
			}
			if matched {
				fmt.Printf("Ignoring path: %s\n", src)
				continue GLOB
			}
		}
		fi, err := os.Stat(src)
		if err != nil {
			return nil
		}
		switch link.Type {
		case "directory":
			if !fi.IsDir() {
				continue
			}
		case "file":
			if fi.IsDir() {
				continue
			}
		}
		dst := strings.Replace(src, pkg.Path, target, 1)
		err = fileutil.Link(src, dst)
		if err != nil {
			return err
		}
		// filepath.Rel(pkg.Path, dst)
		fmt.Printf("Link: %s into %s\n", src, dst)
	}
	return nil
}

func (pkg *Package) GetLines() map[string]string {
	pkg.UnmarshalKey("lines", &pkg.Lines)
	// if pkg.Lines == nil {
	// 	pkg.Lines = make(Lines, 0)
	// }
	// fmt.Printf("lines (%T): %+v\n", lines)
	return pkg.Lines
}

func (pkg *Package) SyncLines(target string) error {
	for file, line := range pkg.GetLines() {
		file = os.ExpandEnv(file)
		file = filepath.Join(target, file)
		err := fileutil.LineInFile(file, line)
		if err != nil {
			return err
		}
		fmt.Printf("Line: '%s' in %s\n", line, file)
	}
	return nil
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
