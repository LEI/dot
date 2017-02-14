package main

import (
	"fmt"
	"os"
	"runtime"
	"gopkg.in/urfave/cli.v1"
	"github.com/urfave/cli/altsrc"
	// "github.com/LEI/dot/cli/command"
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
	ConfigFile  = ""
	ConfigDir   = ".dot"
	ConfigName  = ".dotrc"
	Packages    cli.StringSlice // []Package //= make(PackageSlice, 0)
	IgnoreFiles = []string{".git", ".*\\.md"}
)

type Package struct {
	Path string
}

func main() {
	app := cli.NewApp()
	app.Name = "dot"
	app.Usage = "Dotfiles manager"
	app.Version = version
	// app.Author = ""

	app.Commands = []cli.Command{
		// command.Sync,
	}

	app.Flags = []cli.Flag{
		cli.BoolTFlag{
			Name: "sync, S",
			Usage: "Synchronize packages",
			Destination: &Sync,
		},
		cli.BoolFlag{
			Name: "remove, R",
			Usage: "Remove packages",
			Destination: &Remove,
		},
		altsrc.NewStringFlag(cli.StringFlag{
			Name: "source, s",
			Value: CurrentDir,
			Usage: "Source `directory`",
			Destination: &Source,
			EnvVar: "DOT_SOURCE",
		}),
		altsrc.NewStringFlag(cli.StringFlag{
			Name: "target, t",
			Value: HomeDir,
			Usage: "Destination `directory`",
			Destination: &Target,
			EnvVar: "DOT_TARGET",
		}),
		cli.StringFlag{
			Name: "config, c",
			Usage: "Configuration `file`",
			Destination: &ConfigFile,
			EnvVar: "DOT_CONFIG",
		},
		cli.StringSliceFlag{
			Name: "packages, p",
			Usage: "Package list `name=user/repo`",
			Value: &Packages,
		},
	}

	// app.Before = altsrc.InitInputSourceWithContext(
	// 	app.Flags,
	// 	altsrc.NewYamlSourceFromFlagFunc("config"))

	app.Action = func(c *cli.Context) error {
		// fmt.Printf("%+v\n", c.App)

		// fmt.Printf("%#v\n", c.App.VisibleCategories())
		// fmt.Printf("%#v\n", c.App.VisibleCommands())
		// fmt.Printf("%#v\n", c.App.VisibleFlags())

		fmt.Println("Packages:", c.StringSlice("packages"))

		fmt.Printf("NArg: %#v\n", c.NArg())
		fmt.Printf("NumFlags: %#v\n", c.NumFlags())

		switch {
		case !exists(Source):
			return cli.NewExitError("Source directory does not exist", 1)
		case !exists(Target):
			return cli.NewExitError("Target directory does not exist", 1)
		case Sync && Remove:
			return cli.NewExitError("Sync & Remove are incompatible", 1)
		}
		return nil
	}

	app.Run(os.Args)
}

func exists(path string) bool {
	if _, err := os.Stat(path); err != nil {
		if os.IsNotExist(err) {
			return false
		}
	}
	return true
}
