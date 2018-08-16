package system

import (
	"fmt"
	"io"
	"os"
)

// CheckCopy ... (verify/validate)
func CheckCopy(src, dst string) error {
	// fmt.Println("CheckCopy", src, dst)
	// if src == "" || dst == "" {
	// 	return fmt.Errorf("missing symlink arg: [src:%s dst:%s]", src, dst)
	// }
	if !Exists(src) {
		// return ErrIsNotExist
		return fmt.Errorf("%s: no such file to copy to %s", src, dst)
	}
	if !Exists(dst) {
		// Stop here if the target does not exist
		return nil
	}
	ok, err := store.CompareFile(dst) // (src, dst)
	if err != nil {
		return err
	}
	if !ok {
		return ErrFileExist
	}
	return ErrFileAlreadyExist
}

// Copy ...
func Copy(src, dst string) error {
	// if src == "" || dst == "" {
	// 	return fmt.Errorf("missing symlink arg! [src:%s dst:%s]", src, dst)
	// }
	if DryRun {
		return nil
	}
	in, err := os.Open(src)
	if err != nil {
		return err
	}
	defer in.Close()
	out, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer out.Close()
	// defer func() {
	// 	cerr := out.Close()
	// 	if err == nil {
	// 		err = cerr
	// 	}
	// }()
	if _, err = io.Copy(out, in); err != nil {
		return err
	}
	if err := out.Sync(); err != nil {
		return err
	}
	return store.PutFile(dst)
}
