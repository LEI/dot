package system

import (
	"fmt"
	// "os"

	"github.com/LEI/dot/pkg/cache"
)

var (
	// DryRun disables task execution
	DryRun bool

	// ErrIsNotExist ...
	ErrIsNotExist = fmt.Errorf("file or directory does not exists") // no such file or directory

	// ErrFileExist ...
	ErrFileExist = fmt.Errorf("file already exists")

	// ErrLinkExist ...
	ErrLinkExist = fmt.Errorf("symlink already exists")

	// ErrDirExist ...
	ErrDirExist = fmt.Errorf("directory already exists")

	cacheDir = "dot"
	cacheName = "managed"
	store *cache.Store
)

// Init system
func Init() (err error) {
	store, err = cache.New(cacheDir)
	if err := store.Open(cacheName); err != nil {
		return err
	}
	return nil
}
