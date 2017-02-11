package main

import (
	"bufio"
	// "errors"
	"encoding/json"
	flag "github.com/ogier/pflag"
	"fmt"
	"io/ioutil"
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
	debug         bool
	sourceFlag    string
	targetFlag    string
	defaultSource = os.Getenv("PWD")
	defaultTarget = os.Getenv("HOME")
	config        = Configuration{}
	configFile    string
	configName    = ".dotrc"
	dotDir        = ".dot"
	packages      PackageFlag
	PathSeparator = string(os.PathSeparator)
	InfoSymbol    = "›"
	OkSymbol      = "✓" // ✓ ✔
	ErrSymbol     = "✘" // × ✕ ✖ ✗ ✘
	WarnSymbol    = "!" // ⚠ !
)

type Configuration struct {
	// Target   string
	Packages map[string]Package
	*Package
}

type Package struct {
	Name   string
	Origin string
	Path   string
	Source string
	Target string
	Dir    string
	Dirs   []string
	Link   interface{}
	Links  []interface{}
	Lines  map[string]string
	PreInstall string
	PostInstall string
	OsType string `json:"os_type"`
}

type PackageFlag []string

func (pkg *PackageFlag) String() string {
	return fmt.Sprintf("%s", *pkg)
}

func (pkg *PackageFlag) Set(origin string) error {
	*pkg = append(*pkg, origin)
	// for _, o := range strings.Split(",", origin) 
	return nil
}

// type Link struct{Type string `json:type`, Path string `json:path`}

var f = flag.NewFlagSet("flag", flag.ExitOnError)

func init() {
	// log.SetFlags(log.LstdFlags | log.Lshortfile)

	f.StringVarP(&configFile, "config", "c", "", "Configuration file")
	f.BoolVarP(&debug, "debug", "d", false, "Print more")
	f.StringVarP(&sourceFlag, "source", "s", defaultSource, "Source directory")
	f.StringVarP(&targetFlag, "target", "t", defaultTarget, "Destination directory")
	f.VarP(&packages, "package", "p", "List of packages")

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
	var source, target string

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

	// fmt.Println(configFile)
	// fmt.Println(source)
	// if exists(configFile)

	if sourceFlag != "" {
		source = sourceFlag
	}

	if configFile != "" && source == defaultSource {
		source = filepath.Dir(configFile)
	} else if configFile == "" {
		configFile = filepath.Join(source, configName)
	} else {
		handleError(fmt.Errorf("--config and --source conflict"))
	}

	// if configFile == "" {
	// 	configFile = filepath.Join(source, configName)
	// } else if source == "" {
	// 	source = filepath.Dir(configFile)
	// }

	// if source == "" {
	// 	handleError(fmt.Errorf("Empty source path"))
	// }

	err = readConfig(configFile, &config)
	if err != nil && os.IsExist(err) {
		handleError(err)
		// log.Printf("%s %s\n", configFile, "not found")
	} else if err != nil && os.IsNotExist(err) {
		config.Packages = map[string]Package{}
	}

	// if config.Packages == nil {
	// }

	if len(packages) > 0 {
		for _, p := range packages {
			// fmt.Println(reflect.TypeOf(p))
			pkg := &Package{}
			pkg.Origin = p
			config.Packages[p] = *pkg
		}
	}

	if len(config.Packages) == 0 {
		// fmt.Printf("%s: %+v\n", "No packages found", config)
		pkg := &Package{}
		err = readConfig(configFile, &pkg)
		if err != nil {
			if os.IsNotExist(err) {
				fmt.Printf("Config file not found: %s\n", err.Error())
				os.Exit(1)
			}
			handleError(err)
		}
		config.Packages[filepath.Base(source)] = *pkg
		// config.Packages = map[string]Package{filepath.Base(source): *pkg}
	}

	for name, pkg := range config.Packages {

		if pkg.Name == "" {
			pkg.Name = name
		}

		if pkg.Source != "" {
			pkg.Source = expand(pkg.Source)
		} else {
			pkg.Source = source
		}

		if pkg.Path != "" {
			pkg.Path = expand(pkg.Path)
		} else {
			pkg.Path = pkg.Source
		}

		if pkg.Target != "" {
			pkg.Target = expand(pkg.Target)
		} else {
			pkg.Target = target
		}

		// TODO check target != source

		if pkg.Origin != "" {
			repo := "https://github.com/" + pkg.Origin + ".git"
			// repo := "git@github.com:" + pkg.Origin + ".git"
			if pkg.Path == pkg.Source || pkg.Path == "" {
				pkg.Path = filepath.Join(pkg.Target, dotDir, pkg.Name)
			}
			err := cloneOrPull(pkg.Origin, repo, pkg.Path)

			pkgConfigFile := filepath.Join(pkg.Path, configName)
			err = readConfig(pkgConfigFile, &pkg)
			if err != nil && os.IsExist(err) {
				handleError(err)
			}
		}

		config.Packages[name] = pkg
	}

	for name, pkg := range config.Packages {
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
			"--git-dir", dir + "/.git",
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

func readConfig(path string, v interface{}) error {
	path += ".json"
	_, err := os.Stat(path)
	if err != nil {
		return err
	}
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
