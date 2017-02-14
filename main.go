package main

import (
	"flag"
	"fmt"
	"os"
	"runtime"
	"strings"
	"github.com/spf13/viper"
	"github.com/spf13/pflag"
)

const (
	OS = runtime.GOOS
	version = "master"
)

var (
	HomeDir     = os.Getenv("HOME")
	CurrentDir  = os.Getenv("PWD")
	Sync, Remove bool
	Source, Target string
	Verbose     = 0
	Debug       = false
	ForceYes    = false
	ConfigName  = ".dotrc"
	Packages    PackageMap //= make(PackageSlice, 0)
	IgnoreFiles = []string{".git", ".*\\.md"}
)

type Package struct {
	Name string
	Origin string
	Path string
}

type PackageMap map[string]Package

func (pkg *PackageMap) String() string {
	return fmt.Sprintf("%+v", *pkg)
}

func (pkg *PackageMap) Type() string {
	return fmt.Sprintf("%T", *pkg)
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

func init() {
	pflag.BoolVarP(&Sync, "sync", "S", Sync, "Synchronize packages")
	pflag.BoolVarP(&Remove, "remove", "R", Remove, "Remove packages")

	pflag.StringVarP(&Source, "source", "s", Source, "Source `directory`")
	pflag.StringVarP(&Target, "target", "t", Target, "Destination `directory`")
	// flag.StringVarP(&ConfigFile, "config", "c", "", "Configuration `file`")

	pflag.IntVarP(&Verbose, "verbose", "v", Verbose, "Print more")
	pflag.BoolVarP(&Debug, "debug", "d", Debug, "Check mode")
	pflag.BoolVarP(&ForceYes, "force", "f", ForceYes, "Force yes")

	// flag.VarP(&Packages, "add", "a", "List of packages: `[path=]user/repo`")
	pflag.VarP(&Packages, "packages", "p", "List of packages: `[path=]user/repo`")
}

func main() {
	pflag.CommandLine.AddGoFlagSet(flag.CommandLine)
	flag.Parse()

	viper.SetConfigName(ConfigName)
	viper.AddConfigPath("/etc")
	viper.AddConfigPath(".") // CurrentDir
	// viper.AddConfigPath(Target)

	err := viper.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("Fatal error reading configuration: %s\n", err))
	}

	// viper.SetDefault("Source", CurrentDir)
	// viper.SetDefault("Target", HomeDir)

	fmt.Println(Source, "->", Target)
	fmt.Println(Verbose, Debug, ForceYes)
	fmt.Printf("%s: %+v\n", "Packages", Packages)
}

func exists(path string) bool {
	if _, err := os.Stat(path); err != nil {
		if os.IsNotExist(err) {
			return false
		}
	}
	return true
}
