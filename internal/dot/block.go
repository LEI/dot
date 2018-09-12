package dot

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
)

// Block task
type Block struct {
	Task   `mapstructure:",squash"` // Action, If, OS
	Target string                   // Target file
	Data   string                   // Block content
}

// NewBlock task
func NewBlock(s, d string) *Block {
	return &Block{Target: s, Data: d}
}

func (t *Block) String() string {
	s := fmt.Sprintf("%s:%s", t.Target, t.Data)
	switch Action {
	case "install":
		s = fmt.Sprintf("echo '%s' >> %s", t.Data, tildify(t.Target))
	case "remove":
		s = fmt.Sprintf("rmblock '%s' %s", t.Data, tildify(t.Target))
	}
	return s
}

// Init task
func (t *Block) Init() error {
	// ...
	return nil
}

// Status check task
func (t *Block) Status() error {
	exists, err := blockExists(t.Target, t.Data)
	if err != nil {
		return err
	}
	if exists {
		return ErrExist
	}
	return nil
}

// Do task
func (t *Block) Do() error {
	if err := t.Status(); err != nil {
		switch err {
		case ErrExist, ErrSkip:
			return nil
		default:
			return err
		}
	}
	// Write target file
	f, err := os.OpenFile(t.Target, os.O_CREATE|os.O_APPEND|os.O_WRONLY, defaultFileMode)
	defer f.Close()
	if err != nil {
		return err
	}
	// TODO: insert new line / add DO NOT EDIT
	n, err := f.WriteString(t.Data)
	if err != nil {
		return err
	}
	if n == 0 {
		return fmt.Errorf("%s: written %d bytes", t.Target, n)
	}
	return nil
}

// Undo task
func (t *Block) Undo() error {
	if err := t.Status(); err != nil {
		switch err {
		case ErrExist:
			// continue
		case ErrSkip:
			return nil
		default:
			return err
		}
	}
	b, err := ioutil.ReadFile(t.Target)
	if err != nil {
		return err
	}
	// TODO: remove surrounding empty lines
	b = bytes.Replace(b, []byte(t.Data), []byte{}, 1)
	// Write target file
	if err := ioutil.WriteFile(t.Target, b, defaultFileMode); err != nil {
		return err
	}
	return nil
}

// blockExists returns true if the target file contains the block.
func blockExists(target, data string) (bool, error) {
	if !exists(target) {
		return false, nil
	}
	b, err := ioutil.ReadFile(target)
	if err != nil {
		return false, err
	}
	if bytes.Contains(b, []byte(data)) {
		return true, nil
	}
	return false, nil
}

// func getBlocks(r io.Reader) ([]byte, error) {
// 	scanner := bufio.NewScanner(r)
// 	scanner.Split(bufio.ScanWords)
// 	var result []int
// 	for scanner.Scan() {
// 	    x, err := strconv.Atoi(scanner.Text())
// 	    if err != nil {
// 	        return result, err
// 	    }
// 	    result = append(result, x)
// 	}
// 	return result, scanner.Err()

// f, err := os.Open(t.Target)
// defer f.Close()
// if err != nil {
// 	return err
// }

// bytes.Contains(buf, []byte(data))
// if i := bytes.Index([]byte(t.Data), buf); i != -1 {
// 	fmt.Println(i)
// }

// block := strings.Split(t.Data, "\n")
// if len(block) == 0 {
// 	return fmt.Errorf("no block data")
// }
// scanner := bufio.NewScanner(f)
// var match int
// for scanner.Scan() {
// 	if block[match] == scanner.Text() {
// 		match++
// 	}
// 	if match == len(block)-1 {
// 		fmt.Println("GG")
// 		fmt.Printf("%q ", scanner.Text())
// 		os.Exit(1)
// 	}
// }
// if err := scanner.Err(); err != nil {
// 	fmt.Fprintln(os.Stderr, "reading input:", err)
// }
// }
