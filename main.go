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

const OS = runtime.GOOS

var (
    f = flag.NewFlagSet("flag", flag.ExitOnError)
    // Skip = fmt.Errorf("Skip this path")
    debug bool
    source, target string
    defaultSource = os.Getenv("PWD")
    defaultTarget = os.Getenv("HOME")
    configPath string
    configName = "config.json"
    pkgConfigName = ".json"
)

type Configuration struct {
    Target string
    Packages map[string]Package
}

type Package struct {
    Source string
    Target string
    Dir string
    Dirs []string
    Link interface{}
    Links []interface{}
    Lines map[string]string
}

// type Link struct {
//     Type string `json:type`
//     Path string `json:path`
// }

func init() {
    f.StringVar(&configPath, "c", "", "configuration file")
    f.BoolVar(&debug, "d", false, "enable check-mode")
    // f.BoolVar(&sync, "sync", true, "install")
    f.StringVar(&source, "s", defaultSource, "target directory")
    f.StringVar(&target, "t", defaultTarget, "source directory")

    // flag.ErrHelp = errors.New("flag: help requested")
    f.Usage = func() {
        fmt.Fprintf(os.Stderr, "usage: %s [args]\n", os.Args[0])
        f.PrintDefaults()
    }

    err := f.Parse(os.Args[1:])
    if err != nil {
        log.Fatal(err)
    }
}

func main() {
    err := os.Setenv("OS", OS)
    if err != nil {
        log.Fatal(err)
    }
    // logFlag := func(a *flag.Flag) {
    //     fmt.Println(">", a.Name, "value=", a.Value)
    // }
    // f.Visit(logFlag)
    // fmt.Println(OS, f.Args())

    // fmt.Println(configPath)
    // fmt.Println(source)
    if source == "" {
        log.Fatal("Empty source path")
    }
    if configPath == "" {
        configPath = filepath.Join(source, configName)
    }
    // if exists(configPath)
    config := Configuration{}
    err = readConfig(configPath, &config)
    if err != nil || len(config.Packages) == 0 {
        pkg := Package{}
        err = readConfig(filepath.Join(source, pkgConfigName), &pkg)
        if err != nil {
            log.Fatal(err)
        }
        config.Packages = map[string]Package{filepath.Base(source): pkg}
    }

    // fmt.Println("Source:", source, "/", "Target:", target)

    for name, pkg := range config.Packages {
        err = handlePackage(name, pkg)
        if err != nil {
            log.Fatal(err)
        }
    }

    // scanner := bufio.NewScanner(os.Stdin)
    // for scanner.Scan() {
    //     fmt.Println(scanner.Text()) // Println will add back the final '\n'
    // }
    // if err := scanner.Err(); err != nil {
    //     fmt.Fprintln(os.Stderr, "reading standard input:", err)
    // }

    // err := sync(source)
    // if err != nil {
    //     log.Fatal(err)
    // }

    fmt.Println("[Done]")
}

func handlePackage(name string, pkg Package) error {
    fmt.Printf("%+v\n", name)

    src := source
    if pkg.Source != "" {
        src = expand(pkg.Source)
    }
    dst := target
    if pkg.Target != "" {
        dst = expand(pkg.Target)
    }

    if pkg.Dir != "" {
        pkg.Dirs = append(pkg.Dirs, pkg.Dir)
    }
    fmt.Printf("[%d] Dirs\n", len(pkg.Dirs))
    err := makeDirs(src, dst, pkg.Dirs)
    if err != nil {
        return err
    }

    if pkg.Link != nil && pkg.Link != "" {
        pkg.Links = append(pkg.Links, pkg.Link)
    }
    fmt.Printf("[%d] Links\n", len(pkg.Links))
    err = linkFiles(src, dst, pkg.Links)
    if err != nil {
        return err
    }

    fmt.Printf("[%d] Lines\n", len(pkg.Lines))
    err = linesInFiles(src, dst, pkg.Lines)
    if err != nil {
        return err
    }
    return nil
}

func expand(str string) string {
    str = os.ExpandEnv(str)
    str = strings.Replace(str, "$OS", OS, -1)
    return str
}

// func stat(path string) (os.FileInfo, error) {
//     fi, err := os.Stat(path)
//     if err != nil {
//         msg := strings.Replace(err.Error(), "stat ", os.Args[0]+": ", 1)
//         // fmt.Fprintf(os.Stderr, "%s\n", msg)
//         return fi, fmt.Errorf(msg)
//     }
//     return fi, nil
// }

func makeDirs(src string, dst string, paths []string) error {
    for _, dir := range paths {
        dir = filepath.Join(dst, expand(dir))
        err := os.MkdirAll(dir, 0755)
        if err != nil {
            return err
        }
        fmt.Printf("Created %s\n", dir)
    }
    return nil
}

func linkFiles(source string, target string, globs []interface{}) (error) {
    if _, err := os.Stat(source); err != nil {
        return err
    }
    // var filePaths []string
    for _, glob := range globs {
        switch l := glob.(type) {
            case string:
                src := filepath.Join(source, expand(l))
                paths, _ := filepath.Glob(src)
                // filePaths = append(filePaths, paths...)
                for _, p := range paths {
                    name := strings.Replace(p, source+"/", "", 1)
                    dst := strings.Replace(p, source, target, 1)
                    _, err := os.Stat(dst)
                    if err != nil && os.IsNotExist(err) == false {
                        return err
                    }
                    if os.IsNotExist(err) == false {
                        link, err := filepath.EvalSymlinks(dst)
                        if err != nil {
                            return err
                        }
                        if link == p {
                            fmt.Printf("%s == %s\n", name, dst)
                            continue
                        } else {
                            // fmt.Println(dst, "is a symlink to", link, "not", name)
                            msg := "Do you want to remove "+dst+", which is a symlink to "+link+"?"
                            if ok := confirm(msg); ok {
                                err := os.Remove(dst)
                                if err != nil {
                                    return err
                                }
                            } else {
                                continue
                            }
                        }
                    }

                    err = os.Symlink(p, dst)
                    if err != nil {
                        return err
                    }
                    fmt.Printf("%s -> %s\n", name, dst)
                }
            default:
                fmt.Println("Unhandled type for", l)
        }
    }
    return nil
}

func linesInFiles(src string, target string, lines map[string]string) error {
    for file, line := range lines {
        dst := filepath.Join(target, file)

        contains, err := fileContainsString(dst, line+"\n")
        if err != nil || contains == true {
            continue
        }

        err = appendStringToFile(dst, line+"\n")
        if err != nil {
            return err
        }

        fmt.Printf("Line '%s' -> %s\n", line, dst)
    }
    return nil
}

func fileContainsString(path string, text string) (bool, error) {
    if _, err := os.Stat(path); err != nil && os.IsNotExist(err) == false {
        return false, err
    }
    b, err := ioutil.ReadFile(path)
    if err != nil {
        return false, err
    }
    content := string(b)
    fmt.Println("content of", path, "is", content)
    if content != "" {
        for _, str := range strings.Split(content, "\n") {
            if strings.Contains(str, text) {
                fmt.Printf("Line '%s' already in %s\n", text, path)
                return true, err
            }
        }
    }
    return false, err
}

func appendStringToFile(path string, text string) error {
    // fi, err := os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0611)
    fi, err := os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, os.ModeAppend)
    if err != nil {
        return err
    }
    defer fi.Close()

    // fmt.Fprintf(fi, line+"\n")
    _, err = fi.WriteString(text)
    if err != nil {
        return err
    }
    return nil
}

func readConfig(path string, v interface{}) error {
    // file, _ := os.Open(path)
    // decoder := json.NewDecoder(file)
    // err := decoder.Decode(&config)
    // if _, err := os.Stat(configPath); err != nil {
    //     return err
    // }
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
            log.Fatal(err)
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
//     f, err := os.Open(dirname)
//     if err != nil {
//         return nil, err
//     }
//     defer f.Close()
//     paths, err := f.Readdir(-1) // names
//     if err != nil {
//         return nil, err
//     }
//     // sort.Strings(paths)
//     return paths, nil
// }

// func usage(code int, msg ...string) {
//     if len(msg) > 0 {
//         fmt.Fprintf(os.Stderr, "%s: ", HOME)
//     }
//     for _, m := range msg {
//         fmt.Fprintf(os.Stderr, "%s\n", m)
//     }
//     flag.Usage()
//     os.Exit(code)
// }
