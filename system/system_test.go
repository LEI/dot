package system

import (
	"io/ioutil"
	"log"
	"os"
	"testing"
)

var (
	testDir       string
	testTmpDir    = "" // "/tmp"
	testTmpPrefix = "test"
)

func TestMain(m *testing.M) {
	var err error
	dir, err := ioutil.TempDir(testTmpDir, testTmpPrefix)
	if err != nil {
		log.Fatal(err)
	}
	testDir = dir
	// Override global system cache
	cacheDir = "dot-testing"
	if err := Init(); err != nil {
		log.Fatal(err)
	}
	// store, err = cache.New(filepath.Join(testDir, ".cache"))
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// fmt.Printf("Using test dir: %s\n", dir)
	exitCode := m.Run()
	// defer os.RemoveAll(dir)
	if err := os.RemoveAll(dir); err != nil {
		log.Fatal(err)
	}
	os.Exit(exitCode)
}
