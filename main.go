package main

import (
	"bufio"
	// "errors"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"strings"
)

const (
	OS = runtime.GOOS
)

var (
	f = flag.NewFlagSet("flag", flag.ExitOnError)
	// Skip = fmt.Errorf("Skip this path")
	debug         bool
	source, dest  string
	defaultSource = os.Getenv("PWD")
	defaultDest   = os.Getenv("HOME")
	configPath    string
	configName    = "config.json"
	pkgConfigName = ".json"
	PathSeparator = string(os.PathSeparator)
	InfoSymbol    = "›"
	OkSymbol      = "✓" // ✓ ✔
	ErrSymbol     = "✘" // × ✕ ✖ ✗ ✘
	WarnSymbol    = "!" // ⚠ !
)

type Configuration struct {
	Target   string
	Packages map[string]Package
}

type Package struct {
	Origin string
	Source string
	Target string
	Dir    string
	Dirs   []string
	Link   interface{}
	Links  []interface{}
	Lines  map[string]string
	OsType string `json:"os_type"`
}

// type Link struct{Type string `json:type`, Path string `json:path`}

func init() {
	// log.SetFlags(log.LstdFlags | log.Lshortfile)

	f.StringVar(&configPath, "c", "", "configuration file")
	f.BoolVar(&debug, "d", false, "enable check-mode")
	// f.BoolVar(&sync, "sync", true, "install")
	f.StringVar(&source, "s", defaultSource, "source directory")
	f.StringVar(&dest, "t", defaultDest, "destination directory")

	// flag.ErrHelp = errors.New("flag: help requested")
	f.Usage = func() {
		fmt.Fprintf(os.Stderr, "usage: %s [args]\n", os.Args[0])
		f.PrintDefaults()
	}

	err := f.Parse(os.Args[1:])
	if err != nil {
		handleError(err)
	}
}

func main() {
	err := os.Setenv("OS", OS)
	if err != nil {
		handleError(err)
	}

	// logFlag := func(a *flag.Flag) {
	//     fmt.Println(">", a.Name, "value=", a.Value)
	// }
	// f.Visit(logFlag)
	// fmt.Println(OS, f.Args())

	// fmt.Println(configPath)
	// fmt.Println(source)
	if source == "" {
		handleError(fmt.Errorf("Empty source path"))
	}
	if configPath == "" {
		configPath = filepath.Join(source, configName)
	}
	// if exists(configPath)

	config := Configuration{}

	err = readConfig(configPath, &config)
	if err != nil || len(config.Packages) == 0 {
		log.Printf("%s %s\n", configPath, "not found")
		pkg := Package{}
		pkgConfigPath := filepath.Join(source, pkgConfigName)
		err = readConfig(pkgConfigPath, &pkg)
		if err != nil {
			log.Printf("%s %s\n", pkgConfigPath, "not found")
			os.Exit(1)
		}
		config.Packages = map[string]Package{filepath.Base(source): pkg}
	}

	packages := map[string]Package{}
	for name, pkg := range config.Packages {

		if pkg.Source != "" {
			pkg.Source = expand(pkg.Source)
		} else {
			pkg.Source = source
		}

		if pkg.Target != "" {
			pkg.Target = expand(pkg.Target)
		} else {
			pkg.Target = dest
		}

		if pkg.Origin == "" {
			// fmt.Printf("PKG\n%s %+v\n", name, pkg)
			packages[name] = pkg
		} else {
			// fmt.Println(name, "comes from", pkg.Origin)
			subPkg := Package{}
			subCfgPath := filepath.Join(pkg.Source, pkg.Origin, pkgConfigName)
			// if subPkg == Package{} {
			//     handleError(name+": empty sub-package for origin "+pkg.Origin)
			// }
			if _, err = os.Stat(subCfgPath); err != nil && os.IsExist(err) {
				handleError(err)
			} else {
				err = readConfig(subCfgPath, &subPkg)
				if err != nil {
					handleError(err)
				}
			}
			// fmt.Printf("SUBPKG\n%s %+v\n", name, subPkg)
			subPkg.Source = filepath.Join(pkg.Source, pkg.Origin)
			subPkg.Target = pkg.Target
			// subPkg.Origin = pkg.Origin
			// subPkg.OsType = pkg.OsType
			packages[name] = subPkg
		}
	}

	for name, pkg := range packages {
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
	}
}

func handlePackage(name string, pkg Package) error {
	fmt.Printf("[%+v]\n", name)

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
		err := linkFiles(pkg.Source, pkg.Target, pkg.Links)
		if err != nil {
			return err
		}
	}

	nbLines := len(pkg.Lines)
	if nbLines > 0 {
		fmt.Printf("[%s] %d lines in files\n", name, nbLines)
		err := linesInFiles(pkg.Source, pkg.Target, pkg.Lines)
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

func readConfig(path string, v interface{}) error {
	file, err := ioutil.ReadFile(path)
	if err != nil {
		return err
	}
	err = json.Unmarshal([]byte(string(file)), &v)
	if err != nil {
		return err
	}
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
