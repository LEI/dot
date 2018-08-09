package system

import (
	"fmt"
	// "os"
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
)
