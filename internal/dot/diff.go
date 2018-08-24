package dot

import (
	"bytes"
	"io"
	"os"
	"os/exec"
)

// "github.com/sourcegraph/go-diff"

func getDiff(dst, content string) (string, error) {
	// stdout, stderr, status := ExecCommand("")
	diffCmd := exec.Command("diff", dst, "-")
	// --side-by-side --suppress-common-lines
	stdin, err := diffCmd.StdinPipe()
	if err != nil {
		return "", err
	}
	defer stdin.Close()
	var buf bytes.Buffer
	diffCmd.Stdout = &buf
	diffCmd.Stderr = os.Stderr
	if err := diffCmd.Start(); err != nil {
		return buf.String(), err
	}
	io.WriteString(stdin, content)
	// fmt.Println("WAIT")
	stdin.Close()
	diffCmd.Wait()
	return buf.String(), nil
}

/* // github.com/pmezard/go-difflib/difflib
func getDiff(src, dst, content string) (string, error) {
	b, err := ioutil.ReadFile(dst)
	if err != nil {
		return "", err
	}
	original := string(b)
	// Number of context lines in difflib output
	diffContextLines := 3
	diff := difflib.UnifiedDiff{
		A:        difflib.SplitLines(original),
		B:        difflib.SplitLines(content),
		FromFile: tildify(src), // "Original",
		ToFile:   tildify(dst), // "Current",
		Context:  diffContextLines,
	}
	return difflib.GetUnifiedDiffString(diff)
} */

/* // github.com/sergi/go-diff/diffmatchpatch
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
	// str, _ = url.PathUnescape(str) // url.QueryUnescape(str)
	return str
} */
