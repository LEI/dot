package dot

import (
	"fmt"
	"os"
	"time"
)

var (
	backupExtFmt = "2006-01-02@15:04:05~"
)

func backup(src string) error {
	// Ensure the source file exists
	if _, err := os.Stat(src); err != nil && os.IsNotExist(err) {
		return err
	}
	// Retrieve backup extentsion
	ext := fmt.Sprintf("%s", time.Now().Format(backupExtFmt))
	if ext == "" {
		return fmt.Errorf("empty backup extension for %s with format: %s", src, backupExtFmt)
	}
	// Build destination file path
	dst := fmt.Sprintf("%s.%s", src, ext)
	if _, err := os.Stat(src); os.IsExist(err) {
		return fmt.Errorf("could not make backup of %s to %s: destination already exists", src, dst)
	}
	// Copy over file contents
	err := copyFile(src, dst)
	if err != nil {
		return fmt.Errorf("could not make backup of %s to: %s: %s", src, dst, err)
	}
	fmt.Printf("Backuped %s to %s\n", src, dst)
	return nil
}
