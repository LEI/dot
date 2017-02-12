package main

import (
	"bufio"
	// "errors"
	// "encoding/json"
	"fmt"
	"github.com/jinzhu/configor"
	flag "github.com/ogier/pflag"
	// "io/ioutil"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	// "reflect"
	"runtime"
	"strings"
)

const (
	OS = runtime.GOOS
)

var (
	// Skip = fmt.Errorf("Skip this path")
	debug         = Config.Debug
	defaultSource = os.Getenv("PWD")
	defaultTarget = os.Getenv("HOME")
	Config        = Configuration{}
	ConfigFile    string
	ConfigName    = ".dotrc"
	dotDir        = ".dot"
	packages      PackageFlag
	PathSeparator = string(os.PathSeparator)
	InfoSymbol    = "›"
	OkSymbol      = "✓" // ✓ ✔
	ErrSymbol     = "✘" // × ✕ ✖ ✗ ✘
	WarnSymbol    = "!" // ⚠ !
)

// type Configuration struct {
// 	// Target   string
// 	Packages map[string]Package
// 	*Package
// }

type Configuration struct {
	Source string
	Target string
	Debug  bool
	// Name string `default:"?"`
	Packages PackageMap
	// Source string `required:"true"`
	// Target string
}

type Package struct {
	Name        string
	Origin      string
	Path        string
	Source      string
	Target      string
	Dir         string
	Dirs        []string
	Link        interface{}
	Links       []interface{}
	Lines       map[string]string
	PreInstall  string `json:"pre_install"`
	PostInstall string `json:"post_install"`
	OsType      string `json:"os_type"`
}

type Link struct {
	Type string `json:type`
	Path string `json:path`
}

type PackageMap map[string]Package

type PackageFlag []Package

func (pkg *PackageFlag) String() string {
	return fmt.Sprintf("%s", *pkg)
}

func (pkg *PackageFlag) Set(origin string) error {
	p := &Package{}
	if strings.Contains(origin, "=") {
		s := strings.Split(origin, "=")
		p.Name = s[0]
		p.Origin = s[1]
	} else {
		p.Name = origin
		p.Origin = origin
	}
	// if *pkg == nil {
	// 	*pkg = &PackageFlag{}
	// }
	*pkg = append(*pkg, *p)
	// (*pkg)[p.Name] = *p
	// for _, o := range strings.Split(",", origin)
	return nil
}

var f = flag.NewFlagSet("flag", flag.ExitOnError)

func init() {
	// log.SetFlags(log.LstdFlags | log.Lshortfile)

	f.StringVarP(&ConfigFile, "config", "c", "", "Configuration file")
	f.BoolVarP(&Config.Debug, "debug", "d", false, "Print more")
	f.StringVarP(&Config.Source, "source", "s", defaultSource, "Source directory")
	f.StringVarP(&Config.Target, "target", "t", defaultTarget, "Destination directory")
	f.VarP(&packages, "package", "p", "List of packages")

	// flag.ErrHelp = errors.New("flag: help requested")
	// f.Usage = func() {
	// 	fmt.Fprintf(os.Stderr, "usage: %s [args]\n", os.Args[0])
	// 	f.PrintDefaults()
	// }

	err := f.Parse(os.Args[1:])
	if err != nil {
		handleError(err)
	}

	if Config.Packages == nil {
		Config.Packages = map[string]Package{}
	}

	if len(packages) > 0 {
		for _, pkg := range packages {
			// pkg := &Package{}
			// pkg.Origin = p
			// fmt.Println("PKG ===", pkg)
			Config.Packages[pkg.Name] = pkg
		}
	}
}

func main() {
	err := os.Setenv("OS", OS)
	fmt.Printf("[OS: %s]\n", OS)
	if err != nil {
		handleError(err)
	}

	// logFlag := func(a *flag.Flag) {
	//     fmt.Println(">", a.Name, "value=", a.Value)
	// }
	// f.Visit(logFlag)
	// fmt.Println(OS, f.Args())

	// fmt.Println(ConfigFile)
	// fmt.Println(source)
	// if exists(ConfigFile)

	// if Config.Source != "" {
	// 	source = Config.Source
	// }

	if ConfigFile != "" && Config.Source == defaultSource {
		Config.Source = filepath.Dir(ConfigFile)
	} else if ConfigFile == "" {
		ConfigFile = filepath.Join(Config.Source, ConfigName)
	} else {
		handleError(fmt.Errorf("--config and --source conflict"))
	}

	// if ConfigFile == "" {
	// 	ConfigFile = filepath.Join(source, ConfigName)
	// } else if source == "" {
	// 	source = filepath.Dir(ConfigFile)
	// }

	// if source == "" {
	// 	handleError(fmt.Errorf("Empty source path"))
	// }

	err = readConfig(ConfigFile, &Config)
	if err != nil && os.IsExist(err) {
		handleError(err)
		// log.Printf("%s %s\n", ConfigFile, "not found")
		// } else if err != nil && os.IsNotExist(err) {
		// 	Config.Packages = map[string]Package{}
	}

	if len(Config.Packages) == 0 {
		// fmt.Printf("%s: %+v\n", "No packages found", Config)
		pkg := &Package{}
		err = readConfig(ConfigFile, &pkg)
		if err != nil {
			if os.IsNotExist(err) {
				fmt.Printf("Config file not found: %s\n", err.Error())
				os.Exit(1)
			}
			handleError(err)
		}
		Config.Packages[filepath.Base(Config.Source)] = *pkg
		// Config.Packages = map[string]Package{filepath.Base(source): *pkg}
	}

	for name, pkg := range Config.Packages {

		if pkg.Name == "" {
			pkg.Name = name
		}

		if pkg.Source != "" {
			pkg.Source = expand(pkg.Source)
		} else {
			pkg.Source = Config.Source
		}

		if pkg.Path != "" {
			pkg.Path = expand(pkg.Path)
		} else {
			pkg.Path = pkg.Source
		}

		if pkg.Target != "" {
			pkg.Target = expand(pkg.Target)
		} else {
			pkg.Target = Config.Target
		}

		// TODO check target != source

		if pkg.Origin != "" {
			repo := "https://github.com/" + pkg.Origin + ".git"
			// repo := "git@github.com:" + pkg.Origin + ".git"
			if pkg.Path == pkg.Source || pkg.Path == "" {
				pkg.Path = filepath.Join(pkg.Target, dotDir, pkg.Name)
			}
			err := cloneOrPull(pkg.Origin, repo, pkg.Path)

			pkgConfigFile := filepath.Join(pkg.Path, ConfigName)
			err = readConfig(pkgConfigFile, &pkg)
			if err != nil && os.IsExist(err) {
				handleError(err)
			}
		}

		Config.Packages[name] = pkg
	}

	for name, pkg := range Config.Packages {
		if debug {
			fmt.Printf("%+v\n", pkg)
		}
		err = handlePackage(name, pkg)
		if err != nil {
			handleError(err)
		}
	}

	fmt.Println("[Done]")
}

func handleError(err error) {
	if err != nil {
		pc, fn, line, _ := runtime.Caller(1)
		// log.Printf("[error] %s:%d %v", fn, line, err)
		log.Printf("[error] in %s[%s:%d] %v", runtime.FuncForPC(pc).Name(), fn, line, err)
		log.Fatal(err)
	}
}

func handlePackage(name string, pkg Package) error {
	fmt.Printf("Package: %+v\n", name)

	if pkg.OsType == OS {
		fmt.Printf("[%s] %s: %s", name, "Skipping, only for", pkg.OsType)
		return nil
	}

	if pkg.PostInstall != "" {
		fmt.Println(pkg.PreInstall)
	}

	if pkg.Dir != "" {
		pkg.Dirs = append(pkg.Dirs, pkg.Dir)
	}
	nbDirs := len(pkg.Dirs)
	if nbDirs > 0 {
		fmt.Printf("[%s] %d directories\n", name, nbDirs)
		err := makeDirs(pkg.Target, pkg.Dirs)
		if err != nil {
			return err
		}
	}

	if pkg.Link != nil && pkg.Link != "" {
		pkg.Links = append(pkg.Links, pkg.Link)
	}
	nbLinks := len(pkg.Links)
	if nbLinks > 0 {
		fmt.Printf("[%s] %d symlink patterns\n", name, nbLinks)
		err := linkFiles(pkg.Path, pkg.Target, pkg.Links)
		if err != nil {
			return err
		}
	}

	nbLines := len(pkg.Lines)
	if nbLines > 0 {
		fmt.Printf("[%s] %d lines in files\n", name, nbLines)
		err := linesInFiles(pkg.Path, pkg.Target, pkg.Lines)
		if err != nil {
			return err
		}
	}

	if pkg.PostInstall != "" {
		fmt.Println(pkg.PostInstall)
	}

	return nil
}

// func handleMultiple() error {
// }

func cloneOrPull(name string, repo string, dir string) error {
	_, err := os.Stat(dir)
	if err != nil && os.IsNotExist(err) {
		err := os.MkdirAll(dir, 0755)
		if err != nil {
			return err
		}
		gitClone := exec.Command("git", "clone", repo, dir)
		out, err := gitClone.CombinedOutput()
		if len(out) > 0 {
			fmt.Printf("%s: %s", name, out)
		}
		if err != nil {
			return err
		}
	} else {
		gitPull := exec.Command("git",
			"--git-dir", dir+"/.git",
			"--work-tree", dir,
			"pull")
		out, err := gitPull.CombinedOutput()
		if len(out) > 0 {
			fmt.Printf("%s: %s", name, out)
		}
		if err != nil {
			return err
		}
	}
	return nil
}

func expand(str string) string {
	str = os.ExpandEnv(str)
	str = strings.Replace(str, "$OS", OS, -1)
	return str
}

var ConfigExt = []string{"", ".json", ".yml", ".yaml"}

func readConfig(path string, v interface{}) error {
	var e error
	var paths []string
	for _, ext := range ConfigExt {
		p := path + ext
		_, err := os.Stat(p)
		if err != nil {
			e = err
		} else {
			paths = append(paths, path+ext)
			break
		}
	}

	if len(paths) == 0 && e != nil {
		return e
	}

	configor.Load(&v, paths...)
	// fmt.Printf("config: %#v", v)
	// _, err := os.Stat(path)
	// if err != nil {
	// 	return err
	// }
	// file, err := ioutil.ReadFile(path)
	// if err != nil {
	// 	return err
	// }
	// err = json.Unmarshal([]byte(string(file)), &v)
	// if err != nil {
	// 	return err
	// }
	return nil
}

func confirm(str string) bool {
	reader := bufio.NewReader(os.Stdin)

	for {
		fmt.Printf("%s [y/n]: ", str)

		res, err := reader.ReadString('\n')
		if err != nil {
			handleError(err)
		}

		res = strings.ToLower(strings.TrimSpace(res))

		switch res {
		case "y", "yes":
			return true
		case "n", "no":
			return false
		}
		// TODO limit retries
	}
}

// func readDir(dirname string) ([]os.FileInfo, error) {
// 	f, err := os.Open(dirname)
// 	if err != nil {
// 		return nil, err
// 	}
// 	defer f.Close()
// 	paths, err := f.Readdir(-1) // names
// 	if err != nil {
// 		return nil, err
// 	}
// 	// sort.Strings(paths)
// 	return paths, nil
// }

// func usage(code int, msg ...string) {
// 	if len(msg) > 0 {
// 		fmt.Fprintf(os.Stderr, "%s: ", HOME)
// 	}
// 	for _, m := range msg {
// 		fmt.Fprintf(os.Stderr, "%s\n", m)
// 	}
// 	flag.Usage()
// 	os.Exit(code)
// }
