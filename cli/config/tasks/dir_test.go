package tasks

import (
	"testing"
)

func TestDirTaskCheck(t *testing.T) {
	task := &Dir{
		Path: "/tmp/dir",
	}
	if err := task.Check(); err != nil {
		t.Fatal(err)
	}
	// if err := task.Install(); err != nil {
	// 	t.Fatal(err)
	// }
}
