package dot

// "github.com/sergi/go-diff/diffmatchpatch"

// func printDiff(s, content string) error {
// 	// stdout, stderr, status := ExecCommand("")
// 	diffCmd := exec.Command("diff", s, "-")
// 	// --side-by-side --suppress-common-lines
// 	stdin, err := diffCmd.StdinPipe()
// 	if err != nil {
// 		return err
// 	}
// 	defer stdin.Close()
// 	diffCmd.Stdout = os.Stdout
// 	diffCmd.Stderr = os.Stderr
// 	fmt.Println("START DIFF", s)
// 	if err := diffCmd.Start(); err != nil {
// 		return err
// 	}
// 	io.WriteString(stdin, content)
// 	// fmt.Println("WAIT")
// 	stdin.Close()
// 	diffCmd.Wait()
// 	fmt.Println("END DIFF", s)
// 	return nil
// }
