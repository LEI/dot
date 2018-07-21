package dotfile

import (
	"fmt"
	"os"
	"time"
)

var (
	// BackupDir  ...
	BackupDir string

	// BackupExtFmt ...
	// backups named basename-YYYY-MM-DD@HH:MM:SS~
	// ext = time.strftime("%Y-%m-%d@%H:%M:%S~", time.localtime(time.time()))
	BackupExtFmt = "2006-01-02@15:04:05~"
)

// func init() {
// }

// Backup ...
func Backup(src string) error {
	if _, err := os.Stat(src); err != nil && os.IsNotExist(err) {
		return err
	}
	ext := BackupExt()
	if ext == "" {
		return fmt.Errorf("empty backup extension for %s with format: %s", src, BackupExtFmt)
	}
	dst := fmt.Sprintf("%s.%s", src, ext)
	if _, err := os.Stat(src); os.IsExist(err) {
		return fmt.Errorf("could not make backup of %s to %s: destination already exists", src, dst)
	}
	// FIXME no such file or directory
	changed, err := Copy(src, dst)
	if err != nil {
		return fmt.Errorf("could not make backup of %s to: %s: %s", src, dst, err)
	}
	if changed {
		fmt.Printf("Backuped %s to %s\n", src, dst)
	}
	return nil
}

// BackupExt ...
func BackupExt() string {
	now := time.Now()
	return fmt.Sprintf("%s", now.Format(BackupExtFmt))
}
