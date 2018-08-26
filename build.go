// +build ignore

package main

// https://github.com/restic/restic/blob/master/build.go

import (
	"errors"
	"flag"
	"fmt"
	"go/ast"
	"go/doc"
	"go/parser"
	"go/token"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"reflect"
	"runtime"
	"sort"
	"strconv"
	"strings"
	"text/tabwriter"
	"time"

	"github.com/LEI/dot/internal/cli"
	"github.com/LEI/dot/internal/shell"
)

// Target ...
type Target struct {
	Name string
	Func TargetFunc
	Doc  string
}

// NameSorter ...
type NameSorter []Target

func (a NameSorter) Len() int           { return len(a) }
func (a NameSorter) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a NameSorter) Less(i, j int) bool { return a[i].Name < a[j].Name }

// TargetFunc ...
type TargetFunc func() error

var (
	name        = "dot"                        // name of the program executable and directories
	namespace   = "github.com/LEI/dot"         // subdir of GOPATH
	mainPackage = "github.com/LEI/dot/cmd/dot" // package name for the main package

	constants = map[string]string{
		// "main.packageName": mainPackage,
		"main.version": version(),
		"main.commit":  gitCommit(),
		"main.date":    time.Now().Format("2006-01-02T15:04:05Z0700"),
	}

	defaultBuildTags = []string{}

	listFlag bool
	// testFlag    bool
	vendorOnlyFlag bool
	verboseFlag    bool
	versionFlag    bool

	// docMap map[string]string
	funcMap = map[string]TargetFunc{
		"vendor":        Vendor,
		"dep":           dep,
		"check":         Check,
		"test":          Test,
		"testrace":      TestRace,
		"coverage":      Coverage,
		"vet":           Vet,
		"lint":          Lint,
		"fmt":           Fmt,
		"install":       Install,
		"build":         Build,
		"build:darwin":  Build_Darwin,
		"build:linux":   Build_Linux,
		"build:windows": Build_Windows,
		"clean":         Clean,
		"docs":          Docs,
		"docker":        Docker,
		"dockeros":      DockerOS,
		"goreleaser":    goreleaser,
		"release":       Release,
		"snapshot":      Snapshot,
	}

	targetList = []Target{}

	usageFormat   = "Usage: %s [OPTIONS] TARGET...]\n"
	versionFormat = "dot version %s build script\n"
)

func init() {
	flag.Usage = usage

	// buildFlag = flag.Bool("build", true, "build main binary")
	// testFlag = flag.String("test", "./...", "test packages")
	flag.BoolVar(&listFlag, "l", listFlag, "list targets")
	// flag.BoolVar(&testFlag, "t", testFlag, "only test packages")
	flag.BoolVar(&vendorOnlyFlag, "only", vendorOnlyFlag, "dep ensure -vendor-only")
	flag.BoolVar(&verboseFlag, "v", verboseFlag, "verbose mode")
	flag.BoolVar(&versionFlag, "version", versionFlag, "print version")
}

func usage() {
	showUsage(os.Stdout)
	fmt.Printf("\nOptions:\n")
	printDefaults()
	// os.Exit(0)
}

// showUsage prints a description of the flags
func showUsage(output io.Writer) {
	_, binary := filepath.Split(os.Args[0])
	// output := flag.CommandLine.Output()
	fmt.Fprintf(output, usageFormat, binary)
}

// showTargets prints a list of the build actions
func showTargets(output io.Writer) {
	const padding = 1
	w := tabwriter.NewWriter(output, 0, 0, 1, ' ', 0)
	for _, t := range targetList {
		name := t.Name
		desc := t.Doc
		if desc != "" { // strings.SplitN(desc, " ", 2)
			parts := strings.Fields(desc)
			if strings.ToLower(parts[0]) == name {
				desc = strings.TrimPrefix(desc, parts[0])
			}
		}
		name = strings.Replace(name, "_", ":", 1)
		fmt.Fprintf(w, "  %s\t%s\n", name, desc)
	}
	w.Flush()
}

// PrintDefaults prints, to standard error unless configured otherwise, the
// default values of all defined command-line flags in the set.
func printDefaults() {
	output := os.Stderr // flag.CommandLine.Output()
	w := tabwriter.NewWriter(output, 0, 0, 1, ' ', 0)
	flag.VisitAll(func(f *flag.Flag) {
		s := fmt.Sprintf("  -%s", f.Name) // Two spaces before -; see next two comments.
		name, usage := flag.UnquoteUsage(f)
		if len(name) > 0 {
			s += " " + name
		}
		// // Boolean flags of one ASCII letter are so common we
		// // treat them specially, putting their usage on the same line.
		// if len(s) <= 4 { // space, space, '-', 'x'.
		// 	s += "\t"
		// } else {
		// 	// Four spaces before the tab triggers good alignment
		// 	// for both 4- and 8-space tab stops.
		// 	s += "\n    \t"
		// }
		s += "\t "
		s += strings.Replace(usage, "\n", "\n    \t", -1)

		if !isZeroValue(f, f.DefValue) {
			// if _, ok := f.Value.(*stringValue); ok {
			// 	// put quotes on the value
			// 	s += fmt.Sprintf(" (default %q)", f.DefValue)
			// } else {
			s += fmt.Sprintf(" (default: %v)", f.DefValue)
			// }
		}
		fmt.Fprintf(w, "%s\n", s)
	})
	w.Flush()
}

// isZeroValue guesses whether the string represents the zero
// value for a flag. It is not accurate but in practice works OK.
func isZeroValue(f *flag.Flag, value string) bool {
	// Build a zero value of the flag's Value type, and see if the
	// result of calling its String method equals the value passed in.
	// This works unless the Value type is itself an interface type.
	typ := reflect.TypeOf(f.Value)
	var z reflect.Value
	if typ.Kind() == reflect.Ptr {
		z = reflect.New(typ.Elem())
	} else {
		z = reflect.Zero(typ)
	}
	if value == z.Interface().(flag.Value).String() {
		return true
	}

	switch value {
	case "false", "", "0":
		return true
	}
	return false
}

// Execute build command
func execute() error {
	// Default target when no args are provided
	if len(os.Args[1:]) == 0 {
		usage() // showUsage(os.Stderr)
		// fmt.Printf("\nTargets:\n")
		// showTargets(os.Stdout)
		os.Exit(2) // return nil
	}
	// Parse targets
	ts, err := parse()
	if err != nil {
		return err
	}
	// if listFlag && versionFlag {
	// 	return errors.New("-l and -version cannot be specified at the same time")
	// }
	switch {
	case versionFlag:
		// Print program version and exit
		fmt.Printf(versionFormat, version())
		return nil
	case listFlag:
		fmt.Printf("Targets:\n")
		showTargets(os.Stdout)
		return nil
		// case testFlag:
		// 	return testV()
	}
	return serial(ts...)
}

// serial targets
func serial(ts ...Target) error {
	for _, t := range ts {
		if err := execTarget(t); err != nil {
			return err
		}
	}
	return nil
}

// // serialFuncExit
// func serialFuncExit(tf ...TargetFunc) {
// 	if err := serialFunc(tf...); err != nil {
// 		fmt.Fprintf(os.Stderr, "%s\n", err)
// 		os.Exit(1)
// 	}
// }

// serialFunc targets
func serialFunc(fs ...TargetFunc) error {
	for _, f := range fs {
		if err := f(); err != nil {
			return err
		}
	}
	return nil
}

// Async target funcs
func funk(fs ...TargetFunc) error {
	// done := make(chan bool, 1)
	length := len(fs)
	errs := make(chan error, length)
	for _, f := range fs {
		go func(f TargetFunc) {
			if err := f(); err != nil {
				errs <- err
				return
			}
			errs <- nil
			// close(errs)?
		}(f)
	}
	for i := 0; i < length; i++ {
		if err := <-errs; err != nil {
			fmt.Fprintf(os.Stderr, "failed: %d/%d\n", i+1, length)
			return err
		}
	}
	return nil
}

// Execute target
func execTarget(t Target) error {
	if verboseFlag {
		fmt.Printf("Running target: %s\n", t.Name)
	}
	if err := t.Func(); err != nil {
		return err
	}
	return nil
}

// Parse arguments (targets) and command flags
func parse() ([]Target, error) {
	ts := []Target{}
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
		var t Target
		if j >= 0 {
			t = targetList[j]
		} else {
			f, ok := funcMap[a]
			if !ok {
				// unable to find target
				showUsage(os.Stderr) // usage()
				return ts, fmt.Errorf("target not found: %s", a)
				// return ts, fmt.Errorf(
				// 	"%s: invalid arguments",
				// 	strings.Join(args, " "),
				// )
			}
			t = Target{
				Name: a,
				Func: f,
			}
			// return ts, fmt.Errorf(
			// 	"%s: invalid target in args '%s'",
			// 	a,
			// 	strings.Join(args[1:], " "),
			// )
		}
		// Remove target from arguments once registered
		os.Args = append(os.Args[:i-diff], os.Args[i+1-diff:]...)
		// Append target to queue
		ts = append(ts, t)
	}
	// Parse flag once targets are removed
	flag.Parse()
	return ts, nil
}

// Vendor run dep ensure to install dependencies specified in Gopkg.toml
func Vendor() error {
	if err := dep(); err != nil {
		return err
	}
	env := map[string]string{}
	if runtime.GOOS == "android" {
		env["DEPNOLOCK"] = "1"
	}
	args := []string{"ensure"}
	if vendorOnlyFlag {
		args = append(args, "-vendor-only")
	}
	return runWith(env, "dep", "ensure")
}

// Dep install go dep
func dep() error {
	if executable("dep") {
		// Nothing to be done for 'dep'
		return nil
	}
	if runtime.GOOS == "darwin" {
		return run("brew", "install", "dep")
	}
	// curl -sSL https://raw.githubusercontent.com/golang/dep/master/install.sh | sh
	return run("go", "get", "-u", "github.com/golang/dep/cmd/dep")
}

// Check run tests and linters
func Check() error {
	if strings.Contains(runtime.Version(), "1.8") {
		// Go 1.8 doesn't play along with go test ./... and /vendor.
		// We could fix that, but that would take time.
		fmt.Printf("Skip Check on %s\n", runtime.Version())
		return nil
	}
	// return serialFunc(Test, Vet, Lint, Fmt)
	return funk(Test, Vet, Lint, Fmt)
}

// Test run go tests
func Test() error {
	args := []string{"test"}
	if verboseFlag {
		args = append(args, "-v")
	}
	args = append(args, "./...")
	// return run("go", args...)
	out, err := runOutput("go", args...)
	if out != "" && (verboseFlag || err != nil) {
		fmt.Println(out)
	}
	if err != nil {
		return err
	}
	return nil
}

// Run verbose go tests
func testV() error {
	return run("go", "test", "-v", "./...")
}

// TestRace run go tests with race detector
func TestRace() error {
	return run("go", "test", "-v", "-race", "./...")
}

// Coverage run test coverage
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

// Vet run go vet
func Vet() error {
	args := []string{"vet"}
	if verboseFlag {
		args = append(args, "-v")
	}
	args = append(args, "./...")
	return run("go", args...)
}

// Lint run golint
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
		// so we don't pass -set_exit_status, but we still print out any failures
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

// Fmt run gofmt as a linter
func Fmt() error {
	if !isGoLatest() {
		return nil
	}
	// gofmt -l -s . | grep -v ^vendor/
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
			// should fail this target
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

// List packages
func findPackages() ([]string, error) {
	// if err := dep(); err != nil {
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

// Install package with go install
func Install() error {
	args := []string{
		"install",
		"-ldflags", ldflags(constants),
		"-tags", buildTags(),
	}
	if verboseFlag {
		args = append(args, "-v")
	}
	args = append(args, mainPackage)
	return run("go", args...)
}

// Build binaries for all platforms
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
	// 	"-ldflags", ldflags(constants),
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

// Build_Darwin build binary for macOS
func Build_Darwin() error {
	return buildPlatform("darwin", "amd64")
}

// Build_Linux build binary for Linux
func Build_Linux() error {
	return buildPlatform("linux", "amd64")
}

// Build_Windows build binary for Windows
func Build_Windows() error {
	return buildPlatform("windows", "amd64")
}

// Run go build for a given platform
func buildPlatform(goos, goarch string) error {
	output := "dist/" + goos + "_" + goarch + "/dot"
	args := []string{
		"build",
		"-ldflags", ldflags(constants),
		"-tags", buildTags(),
		"-o", output,
		mainPackage,
	}
	env := map[string]string{
		"CGO_ENABLED": "0",
		"GOOS":        goos,
		"GOARCH":      goarch,
		"GOARM":       "",
	}
	return runWith(env, "go", args...)
}

// Build LDFlags with provided constants
func ldflags(cs map[string]string) string {
	l := make([]string, 0, len(cs))
	for k, v := range cs {
		l = append(l, fmt.Sprintf(`-X "%s=%s"`, k, v))
	}
	return "-s -w " + strings.Join(l, " ")
}

// gitCommit returns a version string that identifies the currently
// checked out git commit
func gitCommit() string {
	// runOutput("git", "rev-parse", "--short", "HEAD")
	// cmd := exec.Command("git", "describe",
	// 	"--long", "--tags", "--dirty", "--always")
	cmd := exec.Command("git", "rev-parse", "--short", "HEAD")
	cmd.Stderr = os.Stderr
	out, err := cmd.Output()
	if err != nil {
		if verboseFlag {
			fmt.Fprintf(os.Stderr, "git returned error: %v\n", err)
		}
		return ""
	}
	return strings.TrimSpace(string(out))
}

// version returns the version string from the file VERSION
// in the current directory
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

// Parse build tags from a given environment variable
func buildTags() string {
	bd := defaultBuildTags
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

// Clean remove dist directory
func Clean() error {
	if _, err := os.Stat("dist"); err != nil && os.IsNotExist(err) {
		return err
	}
	if verboseFlag {
		fmt.Println("Removing dist...")
	}
	return os.RemoveAll("dist")
}

// Docs generates markdown documentation
func Docs() error {
	defaultBuildTags = []string{"doc"}
	// os.Setenv("DOT_BUILD_TAGS", "doc")
	serialFunc(Vendor, Install)
	path := "./docs"
	if err := os.Mkdir(path, 0755); err != nil {
		return err
	}
	return run("dot", "doc", "--md", path)
}

// Run an external command
func run(name string, args ...string) error {
	cmd := exec.Command(name, args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin
	if verboseFlag {
		fmt.Printf("exec: %s %s\n", name, cli.FormatArgs(args))
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
		fmt.Printf("exec: %s %s\n", name, cli.FormatArgs(args))
	}
	return cmd.Run()
}

func runOutput(name string, args ...string) (string, error) {
	// buf := bytes.Buffer{}
	cmd := exec.Command(name, args...)
	// cmd.Stdout = &buf
	// cmd.Stderr = os.Stderr
	if verboseFlag {
		fmt.Printf("exec: %s %s\n", name, cli.FormatArgs(args))
	}
	// if err := cmd.Run(); err != nil {
	// 	return "", err
	// }
	// return buf.String(), nil
	buf, err := cmd.Output()
	s := strings.TrimSuffix(string(buf), "\n")
	return s, err
}

// Check if a command is available
func executable(name string) bool {
	// cmd := exec.Command("command", "-v", name)
	// err := cmd.Run()
	// return err == nil
	out, err := exec.LookPath(name)
	return err == nil && len(out) > 0
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
	_, file, _, ok := runtime.Caller(1) // file = "build.go"
	if !ok || file == "" {
		return nil, fmt.Errorf("invalid program file name %s", file)
	}
	filename := filepath.Base(file)
	filter := func(f os.FileInfo) bool {
		return f.Name() == filename // fm[f.Name()]
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

// Docker build container based on OS env var (default: debian)
func Docker() error {
	// serial(Vendor, Check)
	envOS, ok := os.LookupEnv("OS")
	if !ok {
		// Build from golang if OS is undefined
		return testDockerCompose("base", "test")
		// return errors.New("OS is undefined")
	}
	if envOS == "" {
		return errors.New("OS is empty")
	}
	return testDockerOS(envOS)
	// if err := testDockerCompose("test_os", "test_os"); err != nil {
	// 	return err
	// }
	// return nil
}

// DockerOS build all OS containers
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
	goarch := "amd64"
	if err := buildPlatform("linux", goarch); err != nil {
		return err
	}
	// for _, platform := range list {
	// 	if platform != "windows" {
	// 		continue
	// 	}
	// 	if err := buildPlatform("windows", goarch); err != nil {
	// 		return err
	// 	}
	// 	break
	// }
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
func testDockerCompose(build, test string) error {
	if err := run("docker-compose", "build", build); err != nil {
		return err
	}
	if err := run("docker-compose", "run", test); err != nil {
		return err
	}
	return nil
}

// Release releases with goreleaser
func Release() error {
	if err := goreleaser(); err != nil {
		return err
	}
	return run("goreleaser", "--rm-dist")
}

func goreleaser() error {
	if executable("goreleaser") {
		return nil
	}
	if runtime.GOOS == "darwin" {
		return run("brew", "install", "goreleaser/tap/goreleaser")
	}
	if err := dep(); err != nil {
		return err
	}
	repo := "github.com/goreleaser/goreleaser"
	// curl -sSL https://raw.githubusercontent.com/golang/dep/master/install.sh | sh
	if err := run("go", "get", "-d", repo); err != nil {
		return err
	}
	// Installation command
	c := "dep ensure -vendor-only && make setup build"
	if err := run(shell.Get(), "-c", "cd $GOPATH/src/"+repo+"; "+c); err != nil {
		return err
	}
	return run("go", "install", repo)
}

// Snapshot creates a snapshot release
func Snapshot() error {
	if err := goreleaser(); err != nil {
		return err
	}
	args := []string{"--rm-dist", "--snapshot"}
	if debug := os.Getenv("DEBUG"); debug == "1" {
		args = append(args, "--debug")
	}
	return run("goreleaser", args...)
}

// https://github.com/hashicorp/go-version
func isGoLatest() bool {
	// return strings.Contains(runtime.Version(), "1.10")
	ver := runtime.Version()
	ver = strings.TrimPrefix(ver, "go")
	parts := strings.SplitN(ver, ".", 3)
	if len(parts) < 2 {
		panic("invalid go version " + ver)
	}
	major, err := strconv.Atoi(parts[0])
	if err != nil {
		panic(err)
	}
	minor, err := strconv.Atoi(parts[1])
	if err != nil {
		panic(err)
	}
	// patch := strconv.Atoi(parts[2])
	return major >= 1 && minor >= 10
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
	sort.Sort(NameSorter(targetList))
}

func main() {
	if err := execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
