package main

import (
	// "errors"
	// "encoding/json"
	"fmt"
	"github.com/LEI/dot/git"
	// "github.com/LEI/dot/role"
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
	// OSTYPE = os.Getenv("OSTYPE")
	HomeDir     = os.Getenv("HOME")
	CurrentDir  = os.Getenv("PWD")
	Sync        = true
	Remove      = false
	Verbose     = 0
	Debug       = false
	ForceYes    = false
	ConfigFile  = ""
	ConfigDir   = ".dot"
	ConfigName  = ".dotrc"
	ConfigExts  = []string{"", ".json", ".yml", ".yaml", ".toml"}
	IgnoreFiles = []string{".git", ".*\\.md"}
)

var (
	InfoSymbol    = "›"
	SuccessSymbol = "✓" // ✓ ✔
	ErrorSymbol   = "✘" // × ✕ ✖ ✗ ✘
	WarnSymbol    = "!" // ⚠ !
	logInfo       = log.New(os.Stdout, InfoSymbol+" ", 0)
	logSuccess    = log.New(os.Stdout, SuccessSymbol+" ", 0)
	logWarn       = log.New(os.Stderr, WarnSymbol+" ", log.Lshortfile)
	logError      = log.New(os.Stderr, ErrorSymbol+" ", log.Llongfile)
	// Skip = fmt.Errorf("Skip this path")
)

var Config = Configuration{
	Source:   CurrentDir,
	Target:   HomeDir,
	Packages: make(PackageMap, 0),
}

// type Configuration struct {
// 	// Target   string
// 	Packages map[string]Package
// 	*Package
// }

type Configuration struct {
	Source string
	Target string
	// Debug, ForceYes
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
	PreRemove   string `json:"pre_remove"`
	PostRemove  string `json:"post_remove"`
	Os          OsType // `json:"os_type"`
	repo        *git.Repo
}

type PackageMap map[string]Package

func (pkg *PackageMap) String() string {
	return fmt.Sprintf("%+v", *pkg)
}

func (pkg *PackageMap) Set(origin string) error {
	p := &Package{}
	if strings.Contains(origin, "=") {
		s := strings.Split(origin, "=")
		p.Name = s[0]
		p.Origin = s[1]
	} else {
		p.Name = origin
		p.Origin = origin
	}
	// *pkg = append(*pkg, *p)
	(*pkg)[p.Name] = *p
	return nil
}

type Link struct {
	Type string `json:"type"`
	Path string `json:"path"`
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
		logWarn.Printf("could not set value of type %T: %+v\n", val, val)
	}
	return nil
}

// type IncrementFlag struct {
// 	Bool  bool
// 	Level int
// }

// func (i *IncrementFlag) String() string {
// 	return fmt.Sprintf("%d", i.Level)
// 	// return fmt.Sprintf("%t", *i.Level > 0)
// }

// func (i *IncrementFlag) Set(value string) error {
// 	if value != "" {
// 		i.Level++
// 	} else {
// 		i.Level = 0
// 	}
// 	i.Bool = i.Level > 0
// 	return nil
// }

func init() {
	flagSet := flag.NewFlagSet(os.Args[0], flag.ExitOnError)
	// log.SetFlags(log.LstdFlags | log.Lshortfile)

	flagSet.BoolVarP(&Sync, "sync", "S", Sync, "Synchronize packages")
	flagSet.BoolVarP(&Remove, "remove", "R", Remove, "Remove packages")

	flagSet.StringVarP(&Config.Source, "source", "s", Config.Source, "Source `directory`")
	flagSet.StringVarP(&Config.Target, "target", "t", Config.Target, "Destination `directory`")
	flagSet.StringVarP(&ConfigFile, "config", "c", "", "Configuration `file`")
	flagSet.IntVarP(&Verbose, "verbose", "v", Verbose, "Print more")
	flagSet.BoolVarP(&Debug, "debug", "d", Debug, "Check mode")
	flagSet.BoolVarP(&ForceYes, "force", "f", ForceYes, "Force yes")
	flagSet.VarP(&Config.Packages, "add", "a", "List of packages: `[path=]user/repo`")

	// flag.ErrHelp = errors.New("flag: help requested")
	// flagSet.Usage = func() {
	// 	fmt.Fprintf(os.Stderr, "usage: %s [args]\n", os.Args[0])
	// 	flagSet.PrintDefaults()
	// }

	err := flagSet.Parse(os.Args[1:])
	if err != nil {
		handleError(err)
	}

	parseFlag := func(f *flag.Flag) {
		switch f.Name {
		case "remove":
			Sync = false
		}
	}
	flagSet.Visit(parseFlag)

	// processed := flagSet.NFlag()
	remaining := flagSet.NArg()
	if remaining > 0 {
		logWarn.Printf("%d arguments remaining in '%v'\n", remaining, os.Args)
	}
}

func main() {
	fmt.Println()
	logInfo.Printf("[%s]\n", OS)
	defer logInfo.Printf("[%s]\n\n", "Done")

	err := os.Setenv("OS", OS)
	if err != nil {
		handleError(err)
	}

	// Validate arguments
	if Sync && Remove {
		handleError(fmt.Errorf("--sync and --remove cannot be used together"))
	}
	if ConfigFile != "" && Config.Source != CurrentDir {
		handleError(fmt.Errorf("config (-c %s) and source (-s %s) cannot be used together", ConfigFile, Config.Source))
	}

	// Cleanup source path
	if !filepath.IsAbs(Config.Source) {
		str, err := filepath.Abs(Config.Source)
		if err != nil {
			handleError(err)
		}
		Config.Source = str
	}
	Config.Source = filepath.Clean(Config.Source)

	// Look for config file if no package manually added
	if len(Config.Packages) == 0 || confirm("Ignore global configuration files?") != true {
		err = handleConfig(&Config)
		if err != nil {
			handleError(err)
		}
	}

	if len(Config.Packages) == 0 {
		handleError(fmt.Errorf("empty package list"))
	}

	for name, pkg := range Config.Packages {
		err = handlePackage(name, pkg)
		if err != nil {
			handleError(err)
		}
	}
}

func handleError(err error) {
	if err != nil {
		// pc, fn, line, _ := runtime.Caller(1)
		// fmt.Fprintf(os.Stderr, "Error: %s[%s:%d] %v\n", runtime.FuncForPC(pc).Name(), fn, line, err)
		// log.Printf("[error] %s:%d %v", fn, line, err)
		logError.Printf("%v", err)
		os.Exit(1)
	}
}

func handleConfig(Config *Configuration) error {
	// /etc/dotrc

	if ConfigFile != "" {
		Config.Source = filepath.Dir(ConfigFile)
	} else if ConfigFile == "" && Config.Source != CurrentDir {
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
				logError.Printf("Config file not found: %s\n", err.Error())
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
	for _, osType := range pkg.Os {
		switch osType {
		case OS: //, OSTYPE, OS + "-" + OSTYPE:
			break
		default:
			logInfo.Printf("[%s] %s: %s\n", name, osType, "skip (only for "+OS+")")
			return nil
		}
	}

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
		if pkg.Path == pkg.Source || pkg.Path == "" {
			pkg.Path = filepath.Join(pkg.Target, ConfigDir, pkg.Name)
		}
		pkg.repo = git.New(pkg.Name, pkg.Path).AddRemote("origin", pkg.Origin)
		pkg.repo.Branch = "master"
	}
	pkg.repo.CloneOrPull()

	pkgConfigFile := filepath.Join(pkg.Path, ConfigName)
	err := readConfig(pkgConfigFile, &pkg)
	if err != nil && os.IsExist(err) {
		return err
	}

	logInfo.Printf("Package: %+v\n", name)
	if Verbose > 2 {
		logInfo.Printf("%+v\n", pkg)
	}

	if pkg.Dir != "" {
		pkg.Dirs = append(pkg.Dirs, pkg.Dir)
		pkg.Dir = ""
	}

	if pkg.Link != nil && pkg.Link != "" {
		pkg.Links = append(pkg.Links, pkg.Link)
		pkg.Link = nil
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

// func handle(name string, pkg Package, fn func(), args ...interface{}) error {
// 	count := len(args)
// 	logInfo.Printf("[%s] %d items\n", name, count)
// 	fmt.Println("ARGS", count, args[count - 1])
// 	return nil
// }

func syncPackage(name string, pkg Package) error {
	if pkg.PreInstall != "" {
		// parts := string.Fields(pkg.PreInstall)
		// exe = parts[0]; args = [1:len(parts)]
		preInstall := exec.Command("sh", "-c", pkg.PreInstall)
		out, err := preInstall.CombinedOutput()
		if len(out) > 0 {
			logInfo.Printf("%s: %s\n", "Pre-install", out)
		}
		if err != nil {
			return err
		}
	}

	nbDirs := len(pkg.Dirs)
	if nbDirs > 0 {
		logInfo.Printf("[%s] %d directories\n", name, nbDirs)
		err := makeDirs(pkg.Target, pkg.Dirs)
		if err != nil {
			return err
		}
	}

	nbLinks := len(pkg.Links)
	if nbLinks > 0 {
		logInfo.Printf("[%s] %d symlink patterns\n", name, nbLinks)
		err := linkFiles(pkg.Path, pkg.Target, pkg.Links)
		if err != nil {
			return err
		}
	}

	nbLines := len(pkg.Lines)
	if nbLines > 0 {
		logInfo.Printf("[%s] %d lines in files\n", name, nbLines)
		err := linesInFiles(pkg.Path, pkg.Target, pkg.Lines)
		if err != nil {
			return err
		}
	}

	if pkg.PostInstall != "" {
		postInstall := exec.Command("sh", "-c", pkg.PostInstall)
		out, err := postInstall.CombinedOutput()
		if len(out) > 0 {
			logInfo.Printf("%s: %s\n", "Post-install", out)
		}
		if err != nil {
			return err
		}
	}

	return nil
}

func removePackage(name string, pkg Package) error {
	if pkg.PreRemove != "" {
		// parts := string.Fields(pkg.PreRemove)
		// exe = parts[0]; args = [1:len(parts)]
		preRemove := exec.Command("sh", "-c", pkg.PreRemove)
		out, err := preRemove.CombinedOutput()
		if len(out) > 0 {
			logInfo.Printf("%s: %s\n", "Pre-remove", out)
		}
		if err != nil {
			return err
		}
	}

	nbDirs := len(pkg.Dirs)
	if nbDirs > 0 {
		logInfo.Printf("[%s] %d directories\n", name, nbDirs)
		err := removeDirs(pkg.Target, pkg.Dirs)
		if err != nil {
			return err
		}
	}

	nbLinks := len(pkg.Links)
	if nbLinks > 0 {
		logInfo.Printf("[%s] %d symlink patterns\n", name, nbLinks)
		err := unlinkFiles(pkg.Path, pkg.Target, pkg.Links)
		if err != nil {
			return err
		}
	}

	nbLines := len(pkg.Lines)
	if nbLines > 0 {
		logInfo.Printf("[%s] %d lines in files\n", name, nbLines)
		err := linesOutFiles(pkg.Path, pkg.Target, pkg.Lines)
		if err != nil {
			return err
		}
	}

	if pkg.PostRemove != "" {
		postRemove := exec.Command("sh", "-c", pkg.PostRemove)
		out, err := postRemove.CombinedOutput()
		if len(out) > 0 {
			logInfo.Printf("%s: %s\n", "Post-remove", out)
		}
		if err != nil {
			return err
		}
	}
	return nil
}

func expand(str string) string {
	str = os.ExpandEnv(str)
	// for _, v := range []string{"OS", "OSTYPE"} ...
	str = strings.Replace(str, "$OS", OS, -1)
	return str
}

func readConfig(path string, v interface{}) error {
	var e error
	var paths []string
	for _, ext := range ConfigExts {
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
	if len(paths) > 0 {
		configor.Load(&v, paths...)
		if Verbose > 0 {
			logInfo.Printf("configor.Load %+v: %+v\n", paths, v)
		}
	}
	// _, err := os.Stat(path)
	// if err != nil {
	// 	return err
	// }
	// fi, err := ioutil.ReadFile(path)
	// if err != nil {
	// 	return err
	// }
	// err = json.Unmarshal([]byte(string(fi)), &v)
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
