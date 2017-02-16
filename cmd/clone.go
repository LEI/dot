package cmd

import (
	// "fmt"
	"github.com/LEI/dot/git"
	"github.com/spf13/cobra"
	// "os"
)

var (
	Worktree string
	Pull bool
)

func init() {
	RootCmd.AddCommand(CloneCmd)

	flags := CloneCmd.Flags()

	flags.StringVarP(&Worktree, "work-tree", "w", "", "Overwrite clone destination")
	flags.BoolVarP(&Pull, "pull", "p", false, "Pull from branch master")
}

var CloneCmd = &cobra.Command{
	// Hidden: true,
	Use: "clone [repo]",
	Short: "",
	Long: ``,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) != 1 {
			er("clone requires at least 1 argument")
		}
		err := clone(args[0], Pull)
		if err != nil {
			er(err)
		}
	},
}

func clone(spec string, pull bool) error {
	repo, err := git.NewRepository(spec)
	if err != nil {
		return err
	}
	if pull {
		err = repo.CloneOrPull()
	} else {
		err = repo.CloneOrPull()
	}
	if err != nil {
		return err
	}
	return nil
}
