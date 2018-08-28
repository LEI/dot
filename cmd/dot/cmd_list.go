package main

import (
	"bytes"
	"fmt"
	"strconv"
	"text/tabwriter"
	"text/template"

	"github.com/LEI/dot/internal/dot"
	"github.com/spf13/cobra"
)

// Options for the list command.
type listOptions struct {
	// all    bool
	// filter []string
	format string
	long   bool
	quiet  bool
	// noTab  bool

	// listLong bool
	// host      string
	// tags      restic.TagLists
	// paths     []string
	// recursive bool
}

var listOpts listOptions

var (
	// List templates
	defaultListTpl = "{{.Name}}\t{{if .Ok}}✓{{else}}×{{end}}\t[{{.Path}}]({{.URL}})"
	quietListTpl   = "{{.Name}}\t{{if .Ok}}✓{{end}}"
	longListTpl    = "{{.}}"
)

var cmdList = &cobra.Command{
	Use:     "list [flags]", //  [snapshotID] [dir...]
	Aliases: []string{"ls"},
	Short:   "List managed files",
	Long:    `The "list" command lists roles and their tasks.`,
	// Example: ``,
	Args:    cobra.NoArgs,
	PreRunE: preRunList,
	RunE:    runList,
	// DisableAutoGenTag: true,
}

func init() {
	cmdRoot.AddCommand(cmdList)

	flags := cmdList.Flags()
	// flags.BoolVarP(&listOpts.all, "all", "a", false, "Show all roles (default hides incompatible platforms)")
	// flags.StringSliceVarP(&listOpts.filter, "filter", "f", []string{}, "Filter task list")
	flags.StringVarP(&listOpts.format, "format", "", "", "Pretty-print roles using a Go template")
	flags.BoolVarP(&listOpts.long, "long", "l", false, "Output role tasks")
	flags.BoolVarP(&listOpts.quiet, "quiet", "q", false, "Only show role names and status")
	// flags.BoolVarP(&listOpts.noTab, "no-tab", "n", false, "Disable tabwriter")

	// flags.BoolVarP(&listOpts.listLong, "long", "l", false, "use a long listing format showing size and mode")
	// flags.StringVarP(&listOpts.host, "host", "H", "", "only consider snapshots for this `host`, when no snapshot ID is given")
	// flags.Var(&listOpts.tags, "tag", "only consider snapshots which include this `taglist`, when no snapshot ID is given")
	// flags.StringArrayVar(&listOpts.paths, "path", nil, "only consider snapshots which include this (absolute) `path`, when no snapshot ID is given")
	// flags.BoolVar(&listOpts.recursive, "recursive", false, "include files in subfolders of the listed directories")

	addActionFlags(cmdList)
}

func preRunList(cmd *cobra.Command, args []string) error {
	if err := setActionEnv(cmd); err != nil {
		return err
	}
	if err := initList(&listOpts); err != nil {
		return err
	}
	// if listOpts.format == "" { // && dotOpts.Verbose > 0 {
	// 	listOpts.format = "{{.}}"
	// }
	// if len(listOpts.filter) > 0 {
	// 	fmt.Fprintf(dotOpts.stderr, "--filter not implemented\n")
	// }
	if err := preRunAction(cmd, args); err != nil {
		if _, ok := err.(*dot.DiffError); !ok {
			return err
		}
	}
	return nil
}

func initList(opts *listOptions) error {
	switch {
	case opts.long:
		if opts.format != "" {
			return fmt.Errorf("--format and --long cannot be specified at the same time")
		}
		if opts.quiet {
			return fmt.Errorf("--long and --quiet cannot be specified at the same time")
		}
		opts.format = longListTpl
	case opts.quiet:
		if opts.format != "" {
			return fmt.Errorf("--format and --quiet cannot be specified at the same time")
		}
		if opts.long {
			return fmt.Errorf("--long and --quiet cannot be specified at the same time")
		}
		opts.format = quietListTpl
	default:
		if opts.format == "" {
			opts.format = defaultListTpl
		}
	}
	return nil
}

func runList(cmd *cobra.Command, args []string) error {
	// if len(args) == 0 && opts.Host == "" && len(opts.tags) == 0 && len(opts.paths) == 0 {
	// 	return errors.Fatal("Invalid arguments, either give one or more snapshot IDs or set filters.")
	// }
	// if !listOpts.all {
	// 	dotConfig.Roles.FilterOS()
	// }
	str, err := rolesTable(dotConfig.Roles, 0, 8, 1, listOpts.format)
	if err != nil {
		return err
	}
	fmt.Fprint(dotOpts.stdout, str)

	// // extract any specific directories to walk
	// var dirs []string
	// if len(args) > 1 {
	// 	dirs = args[1:]
	// 	for _, dir := range dirs {
	// 		if !strings.HasPrefix(dir, "/") {
	// 			return errors.Fatal("All path filters must be absolute, starting with a forward slash '/'")
	// 		}
	// 	}
	// }

	// withinDir := func(nodepath string) bool {
	// 	if len(dirs) == 0 {
	// 		return true
	// 	}

	// 	for _, dir := range dirs {
	// 		// we're within one of the selected dirs, example:
	// 		//   nodepath: "/test/foo"
	// 		//   dir:      "/test"
	// 		if fs.HasPathPrefix(dir, nodepath) {
	// 			return true
	// 		}
	// 	}
	// 	return false
	// }

	// approachingMatchingTree := func(nodepath string) bool {
	// 	if len(dirs) == 0 {
	// 		return true
	// 	}

	// 	for _, dir := range dirs {
	// 		// the current node path is a prefix for one of the
	// 		// directories, so we're interested in something deeper in the
	// 		// tree. Example:
	// 		//   nodepath: "/test"
	// 		//   dir:      "/test/foo"
	// 		if fs.HasPathPrefix(nodepath, dir) {
	// 			return true
	// 		}
	// 	}
	// 	return false
	// }

	// repo, err := OpenRepository(gopts)
	// if err != nil {
	// 	return err
	// }

	// if err = repo.LoadIndex(gopts.ctx); err != nil {
	// 	return err
	// }

	// ctx, cancel := context.WithCancel(gopts.ctx)
	// defer cancel()
	// for sn := range FindFilteredSnapshots(ctx, repo, opts.Host, opts.tags, opts.paths, args[:1]) {
	// 	Verbosef("snapshot %s of %v filtered by %v at %s):\n", sn.ID().Str(), sn.paths, dirs, sn.Time)

	// 	err := walker.Walk(ctx, repo, *sn.Tree, nil, func(nodepath string, node *restic.Node, err error) (bool, error) {
	// 		if err != nil {
	// 			return false, err
	// 		}
	// 		if node == nil {
	// 			return false, nil
	// 		}

	// 		if withinDir(nodepath) {
	// 			// if we're within a dir, print the node
	// 			Printf("%s\n", formatNode(nodepath, node, lsOptions.ListLong))

	// 			// if recursive listing is requested, signal the walker that it
	// 			// should continue walking recursively
	// 			if opts.recursive {
	// 				return false, nil
	// 			}
	// 		}

	// 		// if there's an upcoming match deeper in the tree (but we're not
	// 		// there yet), signal the walker to descend into any subdirs
	// 		if approachingMatchingTree(nodepath) {
	// 			return false, nil
	// 		}

	// 		// otherwise, signal the walker to not walk recursively into any
	// 		// subdirs
	// 		if node.Type == "dir" {
	// 			return false, walker.SkipNode
	// 		}
	// 		return false, nil
	// 	})

	// 	if err != nil {
	// 		return err
	// 	}
	// }
	return nil
}

func templateString(name, format string, data interface{}) (string, error) {
	t, err := template.New(name).Parse(format)
	if err != nil {
		return "", err
	}
	var tpl bytes.Buffer
	if err := t.Execute(&tpl, data); err != nil {
		return "", err
	}
	return tpl.String(), nil
}

func rolesTable(roles []*dot.Role, minWidth, tabWidth, padding int, format string) (string, error) {
	var b bytes.Buffer
	w := tabwriter.NewWriter(&b, minWidth, tabWidth, padding, ' ', 0)
	for i, v := range roles {
		name := "index " + strconv.Itoa(i)
		str, err := templateString(name, format, v)
		if err != nil {
			return "", err
		}
		fmt.Fprintf(w, "%s\n", str)
	}
	if err := w.Flush(); err != nil {
		return b.String(), err
	}
	return b.String(), nil
}
