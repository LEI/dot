package dotfile

import (
	"crypto/md5"
	"fmt"
	"io/ioutil"
	"net/url"
	"os"
	"path/filepath"
	// "strings"

	"github.com/LEI/dot/utils"
)

var (
	// CacheDir ...
	CacheDir string

	// ClearCache  ...
	ClearCache bool // = true

	// ErrCacheKeyNotFound ...
	ErrCacheKeyNotFound = fmt.Errorf("cache entry not found")

	dotCache *Cache
)

func init() {
	CacheDir = os.ExpandEnv("$HOME/.cache/dot")

	dotCache = NewCache() // &Cache{Map: map[string]string{}}
}

// Cache ...
type Cache struct {
	// FIXME: allow different cache types
	// (copies, templates...) TODO init
	// Dir string
	Map map[string]string
}

// NewCache ...
func NewCache() *Cache {
	c := &Cache{
		Map: map[string]string{},
	}
	return c.New() // .Init()
}

// InitCache ...
func InitCache() {
	dotCache.Init()
}

// New ...
func (c *Cache) New() *Cache {
	return &Cache{Map: map[string]string{}}
}

// Init ...
func (c *Cache) Init() *Cache {
	if c.Map == nil {
		// *c = *c.New()
		c.Map = make(map[string]string, 0)
	}
	_, err := CreateDir(CacheDir)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to create directory: %s", CacheDir)
		os.Exit(1)
	}
	if ClearCache {
		if err := c.Clear(); err != nil {
			fmt.Fprintf(os.Stderr, "Unable to clear cache: %s", CacheDir)
			os.Exit(1)
		}
	} else if _, err := c.Read(); err != nil {
		fmt.Fprintf(os.Stderr, "Unable to read cache: %s", CacheDir)
		os.Exit(1)
	}
	return c
}

// Get ...
func (c *Cache) Get(k string) (string, error) {
	// key := cacheSerialize(k)
	v, ok := c.Map[k]
	if !ok {
		return v, ErrCacheKeyNotFound
	}
	return v, nil
}

// Validate ...
func (c *Cache) Validate(k, v string) (bool, error) {
	cached, err := c.Get(k)
	if err != nil && err != ErrCacheKeyNotFound {
		return false, err
	}
	// if cached == "" {
	// 	fmt.Println("Validated empty cache", k, cacheHashValue(v))
	// 	return true, nil
	// }
	if err == ErrCacheKeyNotFound {
		// fmt.Println("Cache key not found", k)
		return false, nil
	}
	return cached == cacheHashValue(v), nil
}

// Put ...
func (c *Cache) Put(k, v string) error {
	key := cacheSerialize(k)
	val := cacheHashValue(v)
	c.Map[key] = val // (*c)
	file := filepath.Join(CacheDir, key)
	// fmt.Println("Write cache", file, val)
	return ioutil.WriteFile(file, []byte(val), FileMode)
}

// Del ...
func (c *Cache) Del(k string) error {
	m := make(map[string]string, 0)
	for key, val := range c.Map {
		if key != k {
			m[key] = val
		}
	}
	c.Map = m
	key := cacheSerialize(k)
	file := filepath.Join(CacheDir, key)
	if !utils.Exist(file) {
		return nil
	}
	return os.Remove(file)
}

// Read ...
func (c *Cache) Read() (map[string]string, error) {
	list, err := c.List()
	if err != nil {
		return c.Map, err
	}
	for _, f := range list {
		b, err := ioutil.ReadFile(f)
		if err != nil {
			return c.Map, err
		}
		_, n := filepath.Split(f)
		k, err := cacheUnserialize(n)
		if err != nil {
			return c.Map, err
		}
		v := string(b)
		c.Map[k] = v // (*c)
	}
	// fmt.Println("Read cache", c)
	return c.Map, nil
}

// Write ...
func (c *Cache) Write() error {
	for k, v := range c.Map {
		err := c.Put(k, v)
		if err != nil {
			return err
		}
	}
	return nil
}

// List ...
func (c *Cache) List() ([]string, error) {
	p := filepath.Join(CacheDir, "*")
	return filepath.Glob(p)
}

// Clear ...
func (c *Cache) Clear() error {
	c.Map = make(map[string]string, 0)
	list, err := c.List()
	if err != nil {
		return err
	}
	for _, f := range list {
		if err := os.Remove(f); err != nil {
			return err
		}
	}
	return nil
}

func cacheSerialize(s string) string {
	return url.QueryEscape(s)
}

func cacheUnserialize(s string) (string, error) {
	s, err := url.QueryUnescape(s)
	if err != nil {
		return s, err
	}
	return s, err
}

// Create hash from file content
func cacheHashValue(s string) string {
	h := md5.Sum([]byte(s))
	return fmt.Sprintf("%x", string(h[:16]))
}
