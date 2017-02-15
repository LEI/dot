package cmd

import (
	"fmt"
	"github.com/LEI/dot/git"
	"github.com/spf13/cobra"
	"os"
)

func init() {
	RootCmd.AddCommand(cloneCmd)
}

var cloneCmd = &cobra.Command{
	// Hidden: true,
	Use: "clone",
	Short: "",
	Long: ``,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) != 1 {
			fmt.Fprintln(os.Stderr, "clone requires 1 argument")
			// cmd.Usage()
			os.Exit(1)
		}
		err := cloneRepo(args[0])
		if err != nil {
			er(err)
		}
	},
}

func cloneRepo(spec string) error {
	repo, err := git.NewRepository(spec)
	if err != nil {
		return err
	}
	err = repo.CloneOrPull()
	if err != nil {
		return err
	}
	return nil
}
