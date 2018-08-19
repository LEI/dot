// +build mage

package main

// https://github.com/gohugoio/hugo/blob/master/magefile.go
// https://github.com/restic/restic/blob/master/build.go
// https://github.com/oxequa/realize

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	"github.com/magefile/mage/mg"
	"github.com/magefile/mage/sh"
)

const (
	name        = "dot"                        // name of the program executable and directories
	namespace   = "github.com/LEI/dot"         // subdir of GOPATH
	mainPackage = "github.com/LEI/dot/cmd/dot" // package name for the main package
)

var (
	// Default target to run when none is specified
	// If not set, running mage will list available targets
	Default = All

	defaultLDFlags = "-s -w"

	defaultBuildTags = []string{} // {"selfupdate"}

	goexe = "go"

	dockerCompose = RunVCmd("docker-compose")
	goTest        = sh.RunCmd(goexe, "test")
	goTestV       = RunVCmd(goexe, "test")
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

// Install go dep and sync vendored dependencies
func Vendor() error {
	mg.Deps(getDep)
	return sh.Run("dep", "ensure")
}

func getDep() error {
	if executable("dep") {
		return nil
	}
	if runtime.GOOS == "darwin" {
		return sh.RunV("brew", "install", "dep")
	}
	// curl -sSL https://raw.githubusercontent.com/golang/dep/master/install.sh | sh
	return sh.Run(goexe, "get", "-u", "github.com/golang/dep/cmd/dep")
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

// func test(args ...string) error {
// 	args = append([]string{"-tags", buildTags()}, args...)
// 	return goTest(args...)
// }

// func testV(args ...string) error {
// 	args = append([]string{"-tags", buildTags()}, args...)
// 	if mg.Verbose() {
// 		args = append([]string{"-v"}, args...)
// 	}
// 	return goTestV(args...)
// }

// Run go tests
func Test() error {
	// verbose := ""
	// if mg.Verbose() {
	// 	verbose = "-v"
	// }
	// return goTest(verbose, "./...")
	return goTest("-v", "./...")
}

// Run go tests with race detector
func TestRace() error {
	return goTest("-v", "-race", "./...")
}

// Run test coverage
func Coverage() error {
	profile := os.Getenv("COVERPROFILE")
	if profile == "" {
		profile = "coverage.txt"
	}
	mode := os.Getenv("COVERMODE")
	if mode == "" {
		mode = "atomic"
	}
	// mg.Deps(Vendor)
	verbose := ""
	if mg.Verbose() {
		verbose = "-v"
	}
	return goTestV(verbose, "-race", "-coverprofile="+profile, "-covermode="+mode, "./...")
}

// Run go vet
func Vet() error {
	// verbose := ""
	// if mg.Verbose() {
	// 	verbose = "-v"
	// }
	// return sh.RunV(goexe, "vet", verbose, "./...")
	return sh.RunV(goexe, "vet", "./...")
}

// Run golint
func Lint() error {
	if !executable("golint") {
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
	// if !executable("goimports") {
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

var pkgPrefixLen = len(namespace)

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

// type Build mg.Namespace

// Build binary for macOS
func Darwin() error {
	return buildDist("darwin", "amd64")
}

// Build binary for Linux
func Linux() error {
	return buildDist("linux", "amd64")
}

// Build binary for Windows
func Windows() error {
	return buildDist("windows", "amd64")
}

// Build binary for a specific platform
func buildDist(platform, arch string) error {
	env := map[string]string{
		"CGO_ENABLED": "0",
		"GOOS":        platform,
		"GOARCH":      arch,
		// "GOARM":       "",
	}
	for k, v := range flagEnv() {
		env[k] = v
	}
	// output := "dist/${GOOS}_${GOARCH}/dot"
	output := "dist/" + platform + "_" + arch + "/dot"
	return build(env, "-o", output, mainPackage)
}

// func buildEnv(args ...string) error {
// 	return build(flagEnv(), args...)
// }

func build(env map[string]string, args ...string) error {
	mg.Deps(Vendor)
	a := []string{"build"}
	ldflags := defaultLDFlags + " " + constantsLDFlags()
	a = append(a, []string{
		"-ldflags", ldflags,
		"-tags", buildTags(),
	}...)
	a = append(a, args...)
	return sh.RunWith(env, goexe, a...)
}

func Clean() error {
	if _, err := os.Stat("dist"); err == nil || os.IsExist(err) {
		fmt.Println("Removing dist...")
	}
	return sh.Rm("dist") // os.Remove("dot")
}

// Run go install
func Install() error {
	// mg.Deps(Vendor)
	ldflags := defaultLDFlags + " " + constantsLDFlags()
	args := []string{
		"install",
		"-ldflags", ldflags,
		"-tags", buildTags(),
		mainPackage,
	}
	return sh.RunWith(flagEnv(), goexe, args...)
}

func constantsLDFlags() string {
	cs := map[string]string{
		"main.version":   getVersionFromFile(),
		"main.commit":    getVersionFromGit(),
		"main.timestamp": time.Now().Format("2006-01-02T15:04:05Z0700"),
	}
	l := make([]string, 0, len(cs))
	for k, v := range cs {
		l = append(l, fmt.Sprintf(`-X "%s=%s"`, k, v))
	}
	return strings.Join(l, " ")
}

func buildTags() string {
	bd := defaultBuildTags
	if envTags := os.Getenv("DOT_BUILD_TAGS"); envTags != "" {
		for _, et := range strings.Split(envTags, " ") {
			bd = append(bd, et)
		}
	}
	if len(bd) == 0 {
		// bd = append(bd, "release")
		return "none"
	}
	for i := range bd {
		bd[i] = strings.TrimSpace(bd[i])
	}
	return strings.Join(bd, " ")
	// return getEnv("DOT_BUILD_TAGS", "none")
}

func flagEnv() map[string]string {
	// hash, _ := parseRev()
	return map[string]string{
		"PACKAGE": mainPackage,
		"VERSION": getVersionFromFile(),
		"COMMIT":  getVersionFromGit(),
		"DATE":    time.Now().Format("2006-01-02T15:04:05Z0700"),
	}
}

/*
// Build container with docker
func Docker() error {
	if err := docker("build", "-t", "hugo", "."); err != nil {
		return err
	}
	// yes ignore errors here
	docker("rm", "-f", "hugo-build")
	if err := docker("run", "--name", "hugo-build", "hugo ls /go/bin"); err != nil {
		return err
	}
	if err := docker("cp", "hugo-build:/go/bin/hugo", "."); err != nil {
		return err
	}
	return docker("rm", "hugo-build")
}
*/

// Build container for a given OS
func Docker() error {
	// mg.SerialDeps(Vendor, Check)
	envOS, ok := os.LookupEnv("OS")
	if !ok {
		// Build from golang if OS is undefined
		return testDockerCompose("base", "test")
		// return fmt.Errorf("OS is undefined")
	}
	if envOS == "" {
		return fmt.Errorf("OS is empty")
	}
	return testDockerOS(envOS)
	// if err := testDockerCompose("test_os", "test_os"); err != nil {
	// 	return err
	// }
	// return nil
}

// Build all OS containers
func DockerOS() error {
	return testDockerOS()
}

var platforms = []string{
	"alpine",
	"archlinux",
	"centos",
	"debian",
}

// Docker compose OS
func testDockerOS(list ...string) error {
	if len(list) == 0 {
		list = platforms
	}
	envOS, _ := os.LookupEnv("OS")
	// mg.Deps(Linux) // Snapshot
	if err := buildDist("linux", "amd64"); err != nil {
		return err
	}
	defer os.Setenv("OS", envOS)
	for _, platform := range list {
		os.Setenv("OS", platform)
		if err := testDockerCompose("test_os", "test_os"); err != nil {
			return err
		}
	}
	return nil
}

// var docker = sh.RunCmd("docker")
func testDockerCompose(build, run string) error {
	if err := dockerCompose("build", build); err != nil {
		return err
	}
	if err := dockerCompose("run", run); err != nil {
		return err
	}
	return nil
}

// Create release
func Release() error {
	mg.Deps(getGoreleaser)
	return sh.RunV("goreleaser", "--rm-dist")
}

func getGoreleaser() error {
	if executable("goreleaser") {
		return nil
	}
	if runtime.GOOS == "darwin" {
		return sh.RunV("brew", "install", "goreleaser/tap/goreleaser")
	}
	mg.Deps(getDep)
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

// Create snapshot release
func Snapshot() error {
	mg.Deps(getGoreleaser)
	args := []string{"--rm-dist", "--snapshot"}
	if debug := os.Getenv("DEBUG"); debug == "1" {
		args = append(args, "--debug")
	}
	return sh.RunV("goreleaser", args...)
}

// func Clean() error {
// 	return sh.Rm("dist")
// }

// RunVCmd uses Exec underneath
func RunVCmd(cmd string, args ...string) func(args ...string) error {
	return func(args2 ...string) error {
		return sh.RunV(cmd, append(args, args2...)...)
	}
}

// getVersionFromGit returns a version string that identifies the currently
// checked out git commit.
func getVersionFromGit() string {
	// sh.Output("git", "rev-parse", "--short", "HEAD")
	cmd := exec.Command("git", "rev-parse", "--short", "HEAD")
	// cmd := exec.Command("git", "describe",
	// 	"--long", "--tags", "--dirty", "--always")
	out, err := cmd.Output()
	if err != nil {
		if mg.Verbose() {
			fmt.Fprintf(os.Stderr, "git returned error: %v\n", err)
		}
		return ""
	}

	version := strings.TrimSpace(string(out))
	if mg.Verbose() {
		fmt.Printf("git version is %s\n", version)
	}
	return version
}

// getVersion returns the version string from the file VERSION in the current
// directory.
func getVersionFromFile() string {
	buf, err := ioutil.ReadFile("VERSION")
	if err != nil {
		if mg.Verbose() {
			fmt.Fprintf(os.Stderr, "error reading file VERSION: %v\n", err)
		}
		return ""
	}

	return strings.TrimSpace(string(buf))
}

func executable(bin string) bool {
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
