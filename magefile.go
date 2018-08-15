// +build mage
// https://github.com/gohugoio/hugo/blob/master/magefile.go

package main

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	"github.com/magefile/mage/mg" // mg contains helpful utility functions, like Deps
	"github.com/magefile/mage/sh"
)

const (
	packageName = "github.com/LEI/dot"
)

var (
	// Default target to run when none is specified
	// If not set, running mage will list available targets
	Default = All

	goexe = "go"
)

func init() {
	if exe := os.Getenv("GOEXE"); exe != "" {
		goexe = exe
	}
}

// Default target
func All() {
	mg.SerialDeps(Vendor, Check, Install)
	// cmd := exec.Command(goexe, "build", "-o", "bin/dot", ".")
}

func getDep() error {
	if has("dep") {
		return nil
	}
	// curl -sSL https://raw.githubusercontent.com/golang/dep/master/install.sh | sh
	return sh.Run(goexe, "get", "-u", "github.com/golang/dep/cmd/dep")
}

// Install go dep and sync vendored dependencies
func Vendor() error {
	mg.Deps(getDep)
	return sh.Run("dep", "ensure")
}

// Run tests and linters
func Check() {
	if strings.Contains(runtime.Version(), "1.8") {
		// Go 1.8 doesn't play along with go test ./... and /vendor.
		// We could fix that, but that would take time.
		fmt.Fprintf(os.Stderr, "Skip Check on %s\n", runtime.Version())
		return
	}
	mg.Deps(Test, Vet, Lint, Fmt)
	// mg.Deps(TestRace)
}

// Run go tests
func Test() error {
	// v := ""
	// if mg.Verbose() {
	// 	v = "-v"
	// }
	// return sh.RunV(goexe, "test", v, "./...") // -tags none
	return sh.Run(goexe, "test", "./...")
}

// Run go tests with race detector
func TestRace() error {
	v := ""
	if mg.Verbose() {
		v = "-v"
	}
	return sh.RunV(goexe, "test", "-race", v, "./...")
}

// Run go vet
func Vet() error {
	// v := ""
	// if mg.Verbose() {
	// 	v = "-v"
	// }
	// return sh.RunV(goexe, "vet", v, "./...")
	return sh.RunV(goexe, "vet", "./...")
}

// Run golint
func Lint() error {
	if !has("golint") {
		if err := sh.Run(goexe, "get", "golang.org/x/lint/golint"); err != nil {
			return err
		}
	}
	pkgs, err := findPackages()
	if err != nil {
		return err
	}
	failed := false
	for _, pkg := range pkgs {
		// We don't actually want to fail this target if we find golint errors,
		// so we don't pass -set_exit_status, but we still print out any failures.
		if _, err := sh.Exec(nil, os.Stderr, nil, "golint", pkg); err != nil {
			// fmt.Fprintf(os.Stderr, "ERROR: running go lint on %q: %v\n", pkg, err)
			fmt.Fprintf(os.Stderr, "%s\n", err)
			failed = true
		}
	}
	if failed {
		return errors.New("errors running golint")
	}
	// -min_confidence=$GOLINT_MIN_CONFIDENCE
	// return sh.RunV("golint", "-set_exit_status", verbose("-v"), "$("+goexe+" list ./...)")
	return nil
}

// Run gofmt linter
// gofmt -l -s . | grep -v ^vendor/
func Fmt() error {
	// if !has("goimports") {
	// 	if err := sh.Run(goexe, "get", "golang.org/x/tools/cmd/goimports"); err != nil {
	// 		return err
	// 	}
	// }
	pkgs, err := findPackages()
	if err != nil {
		return err
	}
	failed := false
	first := true
	for _, pkg := range pkgs {
		files, err := filepath.Glob(filepath.Join(pkg, "*.go"))
		if err != nil {
			return nil
		}
		for _, f := range files {
			// gofmt doesn't exit with non-zero when it finds unformatted code
			// so we have to explicitly look for output, and if we find any, we
			// should fail this target.
			s, err := sh.Output("gofmt", "-l", f)
			if err != nil {
				fmt.Fprintf(os.Stderr, "ERROR: running gofmt on %q: %v\n", f, err)
				failed = true
			}
			if s != "" {
				if first {
					fmt.Fprintln(os.Stderr, "The following files are not gofmt'ed:")
					first = false
				}
				failed = true
				fmt.Fprintln(os.Stderr, s)
			}
		}
	}
	if failed {
		return errors.New("improperly formatted go files")
	}
	return nil
}

func has(bin string) bool {
	err := sh.Run("command", "-v", bin)
	return err == nil
}

// func verbose(s string) string {
// 	// if val, err := strconv.ParseBool(os.Getenv("MAGE_VERBOSE")); err == nil && val {
// 	if !mg.Verbose() {
// 		return ""
// 	}
// 	return s
// }

var pkgPrefixLen = len(packageName)

func findPackages() ([]string, error) {
	mg.Deps(getDep)
	s, err := sh.Output(goexe, "list", "./...")
	if err != nil {
		return nil, err
	}
	pkgs := strings.Split(s, "\n")
	for i := range pkgs {
		pkgs[i] = "." + pkgs[i][pkgPrefixLen:]
	}
	return pkgs, nil
}

func parseRev() (string, error) {
	return sh.Output("git", "rev-parse", "--short", "HEAD")
}

func flagEnv() map[string]string {
	hash, _ := parseRev()
	return map[string]string{
		"PACKAGE":     packageName,
		"COMMIT_HASH": hash,
		"BUILD_DATE":  time.Now().Format("2006-01-02T15:04:05Z0700"),
	}
}

// Run go install
func Install() error {
	// mg.Deps(Vendor)
	// return sh.RunWith(flagEnv(), goexe, "install", packageName)
	return sh.Run(goexe, "install", packageName)
}

// var docker = sh.RunCmd("docker")

// Build container with docker compose
func Docker() error {
	// if err := docker("build", "-t", "hugo", "."); err != nil {
	// 	return err
	// }
	// // yes ignore errors here
	// docker("rm", "-f", "hugo-build")
	// if err := docker("run", "--name", "hugo-build", "hugo ls /go/bin"); err != nil {
	// 	return err
	// }
	// if err := docker("cp", "hugo-build:/go/bin/hugo", "."); err != nil {
	// 	return err
	// }
	// return docker("rm", "hugo-build")
	if err := sh.RunV("docker-compose", "build", "base"); err != nil {
		return err
	}
	if err := sh.RunV("docker-compose", "run", "test"); err != nil {
		return err
	}
	return nil
}

// Build container for each OS
func DockerOS() error {
	// mg.SerialDeps(Vendor, Check)
	mg.Deps(Snapshot)
	envOS := os.Getenv("OS")
	defer os.Setenv("OS", envOS)
	for _, platform := range []string{
		"alpine",
		"archlinux",
		"centos",
		"debian",
	} {
		os.Setenv("OS", platform)
		if err := sh.RunV("docker-compose", "build", "test_os"); err != nil {
			return err
		}
		if err := sh.RunV("docker-compose", "run", "test_os"); err != nil {
			return err
		}
	}
	return nil

}

func getGoreleaser() error {
	if has("goreleaser") {
		return nil
	}
	repo := "github.com/goreleaser/goreleaser"
	installCmd := "dep ensure -vendor-only && make setup build"
	// curl -sSL https://raw.githubusercontent.com/golang/dep/master/install.sh | sh
	if err := sh.Run(goexe, "get", "-d", repo); err != nil {
		return err
	}
	if err := sh.Run("sh", "-c", "cd $GOPATH/src/"+repo+"; "+installCmd); err != nil {
		return err
	}
	return sh.Run(goexe, "install", repo)
}

// Create release
func Release() error {
	mg.Deps(getGoreleaser)
	return sh.RunV("goreleaser", "--rm-dist")
}

// Create snapshot release
func Snapshot() error {
	mg.Deps(getGoreleaser)
	return sh.RunV("goreleaser", "--rm-dist", "--snapshot")
}

// func Clean() error {
// 	return sh.Rm("dist")
// }
