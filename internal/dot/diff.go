package dot

// "github.com/sergi/go-diff/diffmatchpatch"
// "github.com/sourcegraph/go-diff"

import (
	"io/ioutil"

	"github.com/pmezard/go-difflib/difflib"
)

func getDiff(src, dst, content string) (string, error) {
	b, err := ioutil.ReadFile(dst)
	if err != nil {
		return "", err
	}
	original := string(b)
	diff := difflib.UnifiedDiff{
		A:        difflib.SplitLines(original),
		B:        difflib.SplitLines(content),
		FromFile: tildify(src), // "Original",
		ToFile:   tildify(dst), // "Current",
		Context:  diffContextLines,
	}
	return difflib.GetUnifiedDiffString(diff)
}

func printDiff(s, content string) error {
	// // stdout, stderr, status := ExecCommand("")
	// diffCmd := exec.Command("diff", s, "-")
	// // --side-by-side --suppress-common-lines
	// stdin, err := diffCmd.StdinPipe()
	// if err != nil {
	// 	return err
	// }
	// defer stdin.Close()
	// diffCmd.Stdout = os.Stdout
	// diffCmd.Stderr = os.Stderr
	// fmt.Println("START DIFF", s)
	// if err := diffCmd.Start(); err != nil {
	// 	return err
	// }
	// io.WriteString(stdin, a)
	// // fmt.Println("WAIT")
	// stdin.Close()
	// diffCmd.Wait()
	// fmt.Println("END DIFF", s)

	/*
		b, err := ioutil.ReadFile(s)
		if err != nil {
			return err
		}
		newContent := string(b)
		diffStr := diffPatchMatch(newContent, content)
		// fmt.Printf("--- %[1]s\n+++ %[1]s\n%s\n", tildify(s), diffStr)
		fmt.Printf("--- START DIFF %[1]s\n%s\n--- END DIFF %[1]s\n", tildify(s), diffStr)
	*/

	return nil
}

/*
func diffPatchMatch(text1, text2 string) string {
	dmp := diffmatchpatch.New()
	// checkLines := false
	// diffs := dmp.DiffMain(text1, text2, checkLines)
	// diffs := dmp.DiffBisect(text1, text2, time.Date(0001, time.January, 01, 00, 00, 00, 00, time.UTC))

	wSrc, wDst, lineArray := dmp.DiffLinesToRunes(text1, text2)
	diffs := []diffmatchpatch.Diff{}
	diffs = dmp.DiffMainRunes(wSrc, wDst, false)
	diffs = dmp.DiffCharsToLines(diffs, lineArray)
	// return dmp.DiffPrettyText(diffs)
	// return dmp.DiffToDelta(diffs)
	patches := dmp.PatchMake(diffs)
	str := dmp.PatchToText(patches)
	// str = strings.Replace(str, "%0A", "", -1)
	// str, _ = url.QueryUnescape(str)
	return str
}
*/
