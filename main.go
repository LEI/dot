package main

import (
	// "errors"
	// "encoding/json"
	"fmt"
	"github.com/jinzhu/configor"
	// "github.com/LEI/dot/models"
	flag "github.com/ogier/pflag"
	// "io/ioutil"
	// "log"
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
	OSTYPE = os.Getenv("OSTYPE")
	HOME   = os.Getenv("HOME")
	PWD    = os.Getenv("PWD")
	// User          = os.User()
	DefaultSource = PWD
	DefaultTarget = HOME
	Sync          = true
	Remove        = false
	ConfigName    = ".dotrc"
	Config        = Configuration{}
	ConfigFile    string
	ConfigDir     = ".dot"
	Verbose       = 0
	Debug         = false
	ForceYes      = false
	IgnoreFiles   = []string{".git", ".*\\.md"}
	PackageList   PackageFlag
	PathSeparator = string(os.PathSeparator)
	InfoSymbol    = "›"
	OkSymbol      = "✓" // ✓ ✔
	ErrSymbol     = "✘" // × ✕ ✖ ✗ ✘
	WarnSymbol    = "!" // ⚠ !
	// Skip = fmt.Errorf("Skip this path")
)

// type Configuration struct {
// 	// Target   string
// 	Packages map[string]Package
// 	*Package
// }

type Configuration struct {
	Source string
	Target string
	// Debug, ForceYes, Verbose
	// Name string `default:"?"`
	Packages PackageMap
	// Source string `required:"true"`
	// Target string
}

// func (*cfg Configuration) Set(value interface{}) error {
// 	fmt.Print("SET CFG", value)
// 	*cfg = value
// 	return nil
// }

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
	Os          OsType `json:"os_type"`
}

type OsType []string

func (osType *OsType) String() string {
	return fmt.Sprintf("%s", *osType)
}

func (osType *OsType) Set(value interface{}) error {
	switch val := value.(type) {
	case string:
		*osType = append(*osType, val)
	case []string:
		*osType = append(*osType, val...)
	default:
		fmt.Printf("could not set value of type %T: %+v\n", val, val)
	}
	return nil
}

type Link struct {
	Type string `json:"type"`
	Path string `json:"path"`
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
	err := os.Setenv("OS", OS)
	if err != nil {
		handleError(err)
	}
	// log.SetFlags(log.LstdFlags | log.Lshortfile)

	f.BoolVarP(&Sync, "sync", "S", Sync, "Synchronize packages")
	f.BoolVarP(&Remove, "remove", "R", Remove, "Remove packages")

	f.StringVarP(&ConfigFile, "config", "c", "", "Configuration `file`")
	f.IntVarP(&Verbose, "verbose", "v", Verbose, "Verbosity `level`")
	f.BoolVarP(&Debug, "debug", "d", Debug, "Enable check-mode")
	f.BoolVarP(&ForceYes, "force", "f", ForceYes, "Force yes")
	f.VarP(&PackageList, "add", "a", "List of packages: `[path=]user/repo`")
	f.StringVarP(&Config.Source, "source", "s", DefaultSource, "Source `directory`")
	f.StringVarP(&Config.Target, "target", "t", DefaultTarget, "Destination `directory`")

	// flag.ErrHelp = errors.New("flag: help requested")
	// f.Usage = func() {
	// 	fmt.Fprintf(os.Stderr, "usage: %s [args]\n", os.Args[0])
	// 	f.PrintDefaults()
	// }
}

func main() {
	err := f.Parse(os.Args[1:])
	if err != nil {
		handleError(err)
	}

	if Config.Packages == nil {
		Config.Packages = map[string]Package{}
	}

	if len(PackageList) > 0 {
		for _, pkg := range PackageList {
			// pkg := &Package{}
			// pkg.Origin = p
			// fmt.Println("PKG ===", pkg)
			Config.Packages[pkg.Name] = pkg
		}
	}

	// logFlag := func(a *flag.Flag) {
	//     fmt.Println(">", a.Name, "value=", a.Value)
	// }
	// f.Visit(logFlag)
	// fmt.Println(OS, f.Args())

	fmt.Printf("[OS: %s]\n", OS)

	err = handleConfig(&Config)
	if err != nil {
		handleError(err)
	}

	for name, pkg := range Config.Packages {
		err = handlePackage(name, pkg)
		if err != nil {
			handleError(err)
		}
	}

	fmt.Printf("%s\n", "[Done]")
}

func handleError(err error) {
	if err != nil {
		pc, fn, line, _ := runtime.Caller(1)
		// log.Printf("[error] %s:%d %v", fn, line, err)
		fmt.Fprintf(os.Stderr, "Error: %s[%s:%d] %v\n", runtime.FuncForPC(pc).Name(), fn, line, err)
		os.Exit(1)
	}
}

func handleConfig(Config *Configuration) error {
	if ConfigFile != "" && Config.Source != DefaultSource {
		return fmt.Errorf("Can not use --config " + ConfigFile + " with --source " + Config.Source)
	}

	if ConfigFile != "" {
		Config.Source = filepath.Dir(ConfigFile)
	} else if ConfigFile == "" && Config.Source != DefaultSource {
		ConfigFile = filepath.Join(Config.Source, ConfigName)
	}
	// if ConfigFile == "" {
	// 	ConfigFile = filepath.Join(source, ConfigName)
	// } else if source == "" {
	// 	source = filepath.Dir(ConfigFile)
	// }

	// if source == "" {
	// 	return fmt.Errorf("Empty source path")
	// }

	err := readConfig(ConfigFile, &Config)
	if err != nil && os.IsExist(err) {
		return err
		// log.Printf("%s %s\n", ConfigFile, "not found")
		// } else if err != nil && os.IsNotExist(err) {
		// 	Config.Packages = map[string]Package{}
	} else if err == nil && len(Config.Packages) == 0 {
		pkg := &Package{}
		err = readConfig(ConfigFile, &pkg)
		if err != nil {
			if os.IsNotExist(err) {
				fmt.Printf("Config file not found: %s\n", err.Error())
				os.Exit(1)
			}
			return err
		}
		name := filepath.Base(Config.Source)
		Config.Packages[name] = *pkg
		// Config.Packages = map[string]Package{filepath.Base(source): *pkg}
	}

	return nil
}

func handlePackage(name string, pkg Package) error {
	if pkg.Name == "" {
		pkg.Name = name
	}

	if pkg.Source != "" {
		pkg.Source = expand(pkg.Source)
	} else {
		pkg.Source = Config.Source
	}

	if pkg.Target != "" {
		pkg.Target = expand(pkg.Target)
	} else {
		pkg.Target = Config.Target
	}

	if pkg.Path != "" {
		pkg.Path = expand(pkg.Path)
	} else {
		pkg.Path = pkg.Source
	}

	if pkg.Origin != "" {
		repo := "https://github.com/" + pkg.Origin + ".git"
		// repo := "git@github.com:" + pkg.Origin + ".git"
		if pkg.Path == pkg.Source || pkg.Path == "" {
			pkg.Path = filepath.Join(pkg.Target, ConfigDir, pkg.Name)
		}
		err := gitCloneOrPull(pkg.Origin, repo, pkg.Path)
		if err != nil {
			return err
		}
	}

	pkgConfigFile := filepath.Join(pkg.Path, ConfigName)
	err := readConfig(pkgConfigFile, &pkg)
	if err != nil && os.IsExist(err) {
		return err
	}

	fmt.Printf("Package: %+v\n", name)
	// if Verbose > 0 {
	// 	fmt.Printf("%+v\n", pkg)
	// }

	for _, osType := range pkg.Os {
		fmt.Println(osType, "vs", OS, OSTYPE)
		switch osType {
		case OS, OSTYPE, OS + "-" + OSTYPE:
			break
		default:
			fmt.Printf("[%s] %s: %s\n", name, osType, "skip")
			return nil
		}
	}

	// Config.Packages[name] = pkg
	switch {
	case Sync:
		err := syncPackage(name, pkg)
		if err != nil {
			return err
		}
	case Remove:
		err := removePackage(name, pkg)
		if err != nil {
			return err
		}
	}

	return nil
}

func syncPackage(name string, pkg Package) error {
	if pkg.PreInstall != "" {
		// parts := string.Fields(pkg.PreInstall)
		// exe = parts[0]; args = [1:len(parts)]
		preInstall := exec.Command("sh", "-c", pkg.PreInstall)
		out, err := preInstall.CombinedOutput()
		if len(out) > 0 {
			fmt.Printf("%s: %s\n", "Pre-install:", out)
		}
		if err != nil {
			return err
		}
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
		postInstall := exec.Command("sh", "-c", pkg.PostInstall)
		out, err := postInstall.CombinedOutput()
		if len(out) > 0 {
			fmt.Printf("%s: %s\n", "Post-install:", out)
		}
		if err != nil {
			return err
		}
	}

	return nil
}

func removePackage(name string, pkg Package) error {
	fmt.Println(" ---> Should remove", name, pkg)
	return nil
}

func expand(str string) string {
	str = os.ExpandEnv(str)
	// for _, v := range []string{"OS", "OSTYPE"} ...
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
	if Verbose > 1 {
		fmt.Printf("configor.Load %+v: %+v\n", paths, v)
	}

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
