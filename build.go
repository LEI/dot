package main

import (
	"errors"
	"flag"
	"fmt"
	"go/ast"
	"go/doc"
	"go/parser"
	"go/token"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"text/tabwriter"
	"time"
)

// Target ...
type Target struct {
	Name string
	Func targetFunc
	Doc  string
}

type targetFunc func() error

var (
	name        = "dot"                        // name of the program executable and directories
	namespace   = "github.com/LEI/dot"         // subdir of GOPATH
	mainPackage = "github.com/LEI/dot/cmd/dot" // package name for the main package

	listFlag    bool
	testFlag    bool
	verboseFlag bool
	versionFlag bool

	// docMap map[string]string
	funcMap = map[string]targetFunc{
		"vendor":   Vendor,
		"dep":      getDep,
		"check":    Check,
		"test":     Test,
		"testrace": TestRace,
		"coverage": Coverage,
		"vet":      Vet,
		"lint":     Lint,
		"fmt":      Fmt,
		"install":  Install,
		"build":    Build,
	}

	targetList = []Target{}

	versionFormat = "Dot version %s\n"
)

var usageFormat = `Usage: %s [flags] [target...]
`

func init() {
	flag.Usage = usage

	// buildFlag = flag.Bool("build", true, "build main binary")
	// testFlag = flag.String("test", "./...", "test packages")
	flag.BoolVar(&listFlag, "l", listFlag, "list targets")
	flag.BoolVar(&testFlag, "t", testFlag, "only test packages")
	flag.BoolVar(&verboseFlag, "v", verboseFlag, "verbose mode")
	flag.BoolVar(&versionFlag, "V", versionFlag, "print version")
}

// Usage of the flags.
func usage() {
	_, binary := filepath.Split(os.Args[0])
	fmt.Fprintf(flag.CommandLine.Output(), usageFormat, binary)
	flag.PrintDefaults()
	// os.Exit(0)
}

// Execute build command.
func execute() error {
	if len(os.Args) == 1 {
		usage()
		return nil
	}
	// Parse targets
	tl, err := parse()
	if err != nil {
		return err
	}
	switch {
	case listFlag:
		printTargets()
		return nil
	case testFlag:
		return testV() // run("go", "test", "./...")
	case versionFlag:
		fmt.Printf(versionFormat, version())
		return nil
	}
	for _, t := range tl {
		if verboseFlag {
			fmt.Printf("Running target: %s", t.Name)
		}
		if err := t.Func(); err != nil {
			return err
		}
	}
	return nil
}

// Parse arguments (targets) and command flags.
func parse() ([]Target, error) {
	tl := []Target{}
	// args := os.Args[1:]
	args := make([]string, len(os.Args)) // os.Args[1:]
	copy(args, os.Args)
	for i, a := range args {
		if i == 0 {
			continue
		}
		diff := len(args) - len(os.Args)
		// if diff < 0 {
		// 	diff = 0
		// }
		// a := args[i]
		if len(a) > 1 && strings.HasPrefix(a, "-") {
			flag.Parse()
			continue
		}
		j := -1
		for k, t := range targetList {
			if t.Name == a {
				j = k
				break
			}
		}
		if j < 0 {
			// return fmt.Errorf("unable to find target %s", a)
			usage()
			return tl, fmt.Errorf(
				"%s: invalid arguments",
				strings.Join(args, " "),
			)
		}
		t := targetList[j]
		// t, ok := funcMap[a]
		// if !ok {
		// 	// fmt.Fprintf(os.Stderr, "Target not found: %s\n", a)
		// 	usage()
		// 	return tl, fmt.Errorf(
		// 		"%s: invalid arguments",
		// 		strings.Join(args, " "),
		// 	)
		// }
		// Remove target from arguments once registered
		os.Args = append(os.Args[:i-diff], os.Args[i+1-diff:]...)
		// Append target to queue
		tl = append(tl, t)
	}
	flag.Parse()
	return tl, nil
}

// Vendor install dependencies specified in Gopkg.toml.
func Vendor() error {
	if err := getDep(); err != nil {
		return err
	}
	return run("dep", "ensure")
}

// Install go dep.
func getDep() error {
	if executable("dep") {
		return nil
	}
	if runtime.GOOS == "darwin" {
		return run("brew", "install", "dep")
	}
	// curl -sSL https://raw.githubusercontent.com/golang/dep/master/install.sh | sh
	return run("go", "get", "-u", "github.com/golang/dep/cmd/dep")
}

// Check run tests and linters.
func Check() error {
	if err := Test(); err != nil {
		return err
	}
	if err := Vet(); err != nil {
		return err
	}
	if err := Lint(); err != nil {
		return err
	}
	if err := Fmt(); err != nil {
		return err
	}
	return nil
}

// Run go tests.
func Test() error {
	args := []string{"test", "./..."}
	if verboseFlag {
		args = append(args, "-v")
	}
	return run("go", args...)
}

// Run verbose go tests.
func testV() error {
	return run("go", "test", "-v", "./...")
}

// Run go tests with race detector.
func TestRace() error {
	return run("go", "test", "-v", "-race", "./...")
}

// Run test coverage.
func Coverage() error {
	profile := os.Getenv("COVERPROFILE")
	if profile == "" {
		profile = "coverage.txt"
	}
	mode := os.Getenv("COVERMODE")
	if mode == "" {
		mode = "atomic"
	}
	return run("go", "test", "-v", "-race", "-coverprofile="+profile, "-covermode="+mode, "./...")
}

// Run go vet.
func Vet() error {
	args := []string{"vet", "./..."}
	if verboseFlag {
		args = append(args, "-v")
	}
	return run("go", args...)
}

// Run golint.
func Lint() error {
	if !executable("golint") {
		if err := run("go", "get", "golang.org/x/lint/golint"); err != nil {
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
		if verboseFlag {
			fmt.Printf("exec: golint %s\n", pkg)
		}
		cmd := exec.Command("golint", pkg)
		cmd.Stdout = os.Stderr
		if err := cmd.Run(); err != nil {
			fmt.Fprintf(os.Stderr, "ERROR: running go lint on %q: %v\n", pkg, err)
			// fmt.Fprintf(os.Stderr, "%s\n", err)
			failed = true
		}
	}
	if failed {
		return errors.New("errors running golint")
	}
	// -min_confidence=$GOLINT_MIN_CONFIDENCE
	// return runV("golint", "-set_exit_status", verbose("-v"), "$(go list ./...)")
	return nil
}

// // Run gofmt as a linter.
// // gofmt -l -s . | grep -v ^vendor/
func Fmt() error {
	// if !executable("goimports") {
	// 	if err := run("go", "get", "golang.org/x/tools/cmd/goimports"); err != nil {
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
			// so we have to explicitly look for runOutput, and if we find any, we
			// should fail this target.
			s, err := runOutput("gofmt", "-l", f)
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

// List packages.
func findPackages() ([]string, error) {
	// if err := getDep(); err != nil {
	// 	return []string{}, err
	// }
	s, err := runOutput("go", "list", "./...")
	if err != nil {
		return nil, err
	}
	pkgs := strings.Split(s, "\n")
	for i := range pkgs {
		if len(pkgs[i]) == 0 {
			continue
		}
		if len(pkgs[i]) < pkgPrefixLen {
			return pkgs, fmt.Errorf("%s: invalid pkg", pkgs[i])
		}
		pkgs[i] = "." + pkgs[i][pkgPrefixLen:]
	}
	return pkgs, nil
}

// Install dot package.
func Install() error {
	args := []string{
		"install",
		"-ldflags", ldflags(),
		"-tags", buildTags(),
		mainPackage,
	}
	return run("go", args...)
}

// Build binaries for all platforms.
func Build() error {
	platforms := []struct {
		os   string
		arch string
		// arm  string
	}{
		{"darwin", "amd64"},
		{"linux", "amd64"},
		{"windows", "amd64"},
	}
	// args = append([]string{
	// 	"build",
	// 	"-ldflags", ldflags(),
	// 	"-tags", buildTags(),
	// }, args...)
	// return run("go", args...)
	for _, p := range platforms {
		if err := buildPlatform(p.os, p.arch); err != nil {
			return err
		}
	}
	return nil
}

// Run go build for a given platform.
func buildPlatform(goos, goarch string) error {
	output := "dist/" + goos + "_" + goarch + "/dot"
	args := []string{
		"build",
		"-ldflags", ldflags(),
		"-tags", buildTags(),
		"-o", output,
	}
	env := map[string]string{
		"CGO_ENABLED": "0",
		"GOOS":        goos,
		"GOARCH":      goarch,
		"GOARM":       "",
	}
	return runWith(env, "go", args...)
}

// Build LDFlags.
func ldflags() string {
	cs := map[string]string{
		// "main.packageName": mainPackage,
		"main.version":   version(),
		"main.commit":    gitCommit(),
		"main.timestamp": time.Now().Format("2006-01-02T15:04:05Z0700"),
	}
	l := make([]string, 0, len(cs))
	for k, v := range cs {
		l = append(l, fmt.Sprintf(`-X "%s=%s"`, k, v))
	}
	return "-s -w " + strings.Join(l, " ")
}

// gitCommit returns a version string that identifies the currently
// checked out git commit.
func gitCommit() string {
	// runOutput("git", "rev-parse", "--short", "HEAD")
	// cmd := exec.Command("git", "describe",
	// 	"--long", "--tags", "--dirty", "--always")
	out, err := runOutput("git", "rev-parse", "--short", "HEAD")
	if err != nil {
		if verboseFlag {
			fmt.Fprintf(os.Stderr, "git returned error: %v\n", err)
		}
		return ""
	}
	version := strings.TrimSpace(string(out))
	if verboseFlag {
		fmt.Printf("git version is %s\n", version)
	}
	return version
}

// version returns the version string from the file VERSION
// in the current directory.
func version() string {
	buf, err := ioutil.ReadFile("VERSION")
	if err != nil {
		if verboseFlag {
			fmt.Fprintf(os.Stderr, "error reading file VERSION: %v\n", err)
		}
		return ""
	}
	return strings.TrimSpace(string(buf))
}

func buildTags() string {
	bd := []string{} // defaultBuildTags
	if envTags := os.Getenv("DOT_BUILD_TAGS"); envTags != "" {
		for _, et := range strings.Fields(envTags) {
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
}

func clean() error {
	if _, err := os.Stat("dist"); err != nil && os.IsNotExist(err) {
		return err
	}
	if verboseFlag {
		fmt.Println("Removing dist...")
	}
	return os.RemoveAll("dist")
}

// Run an external command.
func run(name string, args ...string) error {
	cmd := exec.Command(name, args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin
	if verboseFlag {
		fmt.Printf("exec: %s %s\n", name, strings.Join(args, " "))
	}
	return cmd.Run()
}

func runWith(env map[string]string, name string, args ...string) error {
	cmd := exec.Command(name, args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin
	cmd.Env = os.Environ()
	for k, v := range env {
		cmd.Env = append(cmd.Env, k+"="+v)
	}
	if verboseFlag {
		fmt.Printf("exec: %s %s\n", name, strings.Join(args, " "))
	}
	return cmd.Run()
}

func runOutput(name string, args ...string) (string, error) {
	// buf := bytes.Buffer{}
	cmd := exec.Command(name, args...)
	// cmd.Stdout = &buf
	// cmd.Stderr = os.Stderr
	if verboseFlag {
		fmt.Printf("exec: %s %s\n", name, strings.Join(args, " "))
	}
	// if err := cmd.Run(); err != nil {
	// 	return "", err
	// }
	// return buf.String(), nil
	buf, err := cmd.Output()
	s := strings.TrimSuffix(string(buf), "\n")
	return s, err
}

// Check if a command is available.
func executable(name string) bool {
	cmd := exec.Command("command", "-v", name)
	err := cmd.Run()
	return err == nil
}

func printTargets() {
	const padding = 1
	fmt.Println("Targets:")
	w := &tabwriter.Writer{}
	w.Init(os.Stdout, 0, 8, 3, '\t', 0)
	for _, t := range targetList {
		fmt.Fprintf(w, "  %s\t%s\n", t.Name, t.Doc)
	}
	w.Flush()
}

// func getFunctionName(i interface{}) string {
// 	return runtime.FuncForPC(reflect.ValueOf(i).Pointer()).Name()
// }

// getPackage
// https://github.com/magefile/mage/blob/master/parse/parse.go
func getPackage(path string, files []string) (*ast.Package, error) {
	fset := token.NewFileSet()
	// fm := make(map[string]bool, len(files))
	// for _, f := range files {
	// 	fm[f] = true
	// }
	filter := func(f os.FileInfo) bool {
		return f.Name() == "build.go" // fm[f.Name()]
	}
	pkgs, err := parser.ParseDir(fset, path, filter, parser.ParseComments)
	if err != nil {
		return nil, fmt.Errorf("failed to parse directory: %v", err)
	}
	for name, pkg := range pkgs {
		if !strings.HasSuffix(name, "_test") {
			return pkg, nil
		}
	}
	return nil, fmt.Errorf("no non-test packages found in %s", path)
}

func parseDoc() (map[string]string, error) {
	m := map[string]string{}
	pkg, err := getPackage(".", []string{"build.go"})
	if err != nil {
		return m, err
	}
	p := doc.New(pkg, "./", 0)
	for _, f := range p.Funcs {
		if f.Recv != "" {
			// skip methods
			continue
		}
		if !ast.IsExported(f.Name) {
			// skip non-exported functions
			continue
		}
		name := strings.ToLower(f.Name)
		docStr := strings.TrimSpace(f.Doc)
		if docStr != "" {
			docStr = strings.Split(docStr, "\n")[0]
		}
		m[name] = docStr
	}
	return m, nil
}

func init() {
	dm, err := parseDoc()
	if err != nil {
		// error while parsing doc
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	for name, docStr := range dm {
		t := Target{
			Name: name,
			Func: funcMap[name],
			Doc:  docStr,
		}
		targetList = append(targetList, t)
	}
}

func main() {
	if err := execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
