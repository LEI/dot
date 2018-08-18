package dot

import "fmt"

var (
	// ErrAlreadyExist ...
	ErrAlreadyExist = fmt.Errorf("already exists")

	// ErrSkip ...
	ErrSkip = fmt.Errorf("skip")
)
