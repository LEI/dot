package dotfile

import (
	"crypto/md5"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

var (
	cacheDir string

	cachePathSep = "=+"

	osPathSep = string(os.PathSeparator)
)

func init() {
	cacheDir = os.ExpandEnv("$HOME/.cache/dot")
	_, err := CreateDir(cacheDir)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to create directory: %s", cacheDir)
		os.Exit(1)
	}
}

// Cache ...
type Cache struct {
	Map map[string]string
}

// Get ...
func (c *Cache) Get(k string) string {
	if c.Map == nil {
		return ""
	}
	// v, ok := c.Map[k]
	return c.Map[k]
}

// Put ...
func (c *Cache) Put(k, v string) error {
	if c.Map == nil {
		c.Map = make(map[string]string, 0)
	}
	// Replace path delimiters
	key := strings.Replace(k, osPathSep, cachePathSep, -1)
	// Create hash from contents
	h := md5.Sum([]byte(v))
	val := fmt.Sprintf("%x", string(h[:16]))
	c.Map[key] = val
	return nil
}

// WriteKey ...
func (c *Cache) WriteKey(k, v string) error {
	if err := c.Put(k, v); err != nil {
		return err
	}
	return c.Write()
}

// Write ...
func (c *Cache) Write() error {
	for k, v := range c.Map {
		f := filepath.Join(cacheDir, k)
		fmt.Println("Write cache", f, v)
		if err := ioutil.WriteFile(f, []byte(v), FileMode); err != nil {
			return err
		}
	}
	return nil
}

// Read ...
func (c *Cache) Read() error {
	return nil
}
