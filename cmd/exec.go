// Copyright © 2017 NAME HERE <EMAIL ADDRESS>
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cmd

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"os/exec"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// type execFunc func(name string, args ...string) error

var (
	// Shell to use for commands
	Shell = "bash"
)

// execCmd represents the exec command
var execCmd = &cobra.Command{
	Use:   "exec",
	Short: "Execute",
	Long:  ``,
	RunE: func(cmd *cobra.Command, args []string) error {
		// r := &role{
		// 	Dir: Source,
		// 	URL: URL,
		// }
		// if err := r.Init(); err != nil {
		// 	return err
		// }
		if len(args) == 0 {
			if err := viper.UnmarshalKey("exec", &args); err != nil {
				return err
			}
		}
		return ExecCommand(args)
	},
}

func init() {
	// DotCmd.AddCommand(execCmd)

	// execCmd.Flags().StringVarP(&URL, "url", "u", URL, "Repository URL")
	execCmd.Flags().StringVarP(&Shell, "shell", "", Shell, "Shell")
}

// ExecCommand ...
func ExecCommand(in []string) error {
	// if len(args) == 0 {
	// 	args = append(args, viper.GetStringSlice("exec")...)
	// }
	// args = append([]string{"-c"}, str)
	for _, str := range in {
		fmt.Printf("%s\n", str)
		// fmt.Printf("%s %s\n", name, strings.Join(args, " "))
		err := executeCmd(Shell, []string{"-c", str}...)
		if err != nil {
			return err
		}
	}
	return nil
}

// Safe guard execution in test mode
func executeCmd(name string, args ...string) error {
	if DryRun {
		return nil
	}
	return execute(name, args...)
}

func execute(name string, args ...string) error {
	c := exec.Command(name, args...)
	c.Stdout = os.Stdout
	c.Stderr = os.Stderr
	return c.Run()
	// out, err := c.CombinedOutput()
	// fmt.Println(out)
	// if err != nil {
	// 	return err
	// }
	// return nil
}

func execStdout(name string, args ...string) (string, error) {
	c := exec.Command(name, args...)
	stdout, err := c.StdoutPipe()
	if err != nil {
		return toString(stdout), err
	}
	if err := c.Start(); err != nil {
		return toString(stdout), err
	}
	// var person struct {
	// 	Name string
	// 	Age  int
	// }
	// if err := json.NewDecoder(stdout).Decode(&person); err != nil {
	// 	return toString(stdout), err
	// }
	out := toString(stdout)
	if err := c.Wait(); err != nil {
		return out, err
	}
	// fmt.Printf("%s is %d years old\n", person.Name, person.Age)
	return out, nil
}

func toString(in io.ReadCloser) string {
	buf := new(bytes.Buffer)
	buf.ReadFrom(in)
	return buf.String()
}

func captureStdout(cb func() error) (string, error) {
	old := os.Stdout // keep backup of the real stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	// fmt.Println("output")
	err := cb()
	if err != nil {
		return "", err
	}

	outC := make(chan string)
	// copy the output in a separate goroutine so printing can't block indefinitely
	go func() {
		var buf bytes.Buffer
		io.Copy(&buf, r)
		outC <- buf.String()
	}()

	// back to normal state
	w.Close()
	os.Stdout = old // restoring the real stdout
	out := <-outC

	// reading our temp stdout
	fmt.Println("previous output:")
	fmt.Print(out)

	return out, nil
}
