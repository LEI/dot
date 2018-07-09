package cmd

import (
	// "fmt"
	// "reflect"

	"github.com/jessevdk/go-flags"
)

// BaseCmd ...
type BaseCmd struct {
	// FIXME *flags.Command
	flags.Commander

	// Role map[string]string `short:"r" long:"role" description:""`
}


/*
// CommandHandler ...
func (cmd *BaseCmd) CommandHandler(command flags.Commander, args []string) error {
	fmt.Println("EXEC", command, cmd, args)

	return command.Execute(args)
}
*/

/*
// IsEmpty ...
func (cmd *BaseCmd) IsEmpty() bool {
	// FIXME cmd -> nil?
	// TODO: Dynamic type instead of BaseCmd
	return reflect.DeepEqual(cmd, BaseCmd{})
}
*/
