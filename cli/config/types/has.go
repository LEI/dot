package types

import (
	"fmt"
	"os"
	"strings"
	"text/template"

	"github.com/LEI/dot/pkg/executils"
	"github.com/LEI/dot/pkg/ostype"
	"github.com/LEI/dot/system"
)

// HasOS ...
type HasOS struct {
	OS Slice
}

// CheckOS ...
func (h *HasOS) CheckOS() bool {
	if len(h.OS) == 0 {
		return true
	}
	// hasOS := false
	// for _, o := range h.OS {
	// 	if ostype.Has(o) {
	// 		hasOS = true
	// 		break
	// 	}
	// }
	// return hasOS
	return ostype.Has(h.OS...)
}

// HasIf ...
type HasIf struct {
	If Slice
}

// CheckIf ...
func (h *HasIf) CheckIf() bool {
	if len(h.If) == 0 {
		return true
	}
	varsMap := map[string]interface{}{
		"DryRun":  system.DryRun,
		// "Verbose": tasks.Verbose,
		// "OS":      runtime.GOOS,
	}
	funcMap := template.FuncMap{
		"hasOS": ostype.Has,
	}
	// https://golang.org/pkg/text/template/#hdr-Functions
	for _, cond := range h.If {
		str, err := system.TemplateData("", cond, varsMap, funcMap)
		if err != nil {
			fmt.Fprintf(os.Stderr, "err tpl: %s\n", err)
			continue
		}
		_, stdErr, status := executils.ExecuteBuf("sh", "-c", str)
		// out := strings.TrimRight(string(stdOut), "\n")
		strErr := strings.TrimRight(string(stdErr), "\n")
		// if out != "" {
		// 	fmt.Printf("stdout: %s\n", out)
		// }
		if strErr != "" {
			fmt.Fprintf(os.Stderr, "'%s' stderr: %s\n", str, strErr)
		}
		if status == 0 {
			return true
		}
	}
	return false
}
