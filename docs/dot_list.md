## dot list

List managed files

### Synopsis

The "list" command lists roles and their tasks.

```
dot list [flags]
```

### Options

```
  -a, --all             Show all roles (default hides incompatible platforms)
      --format string   Pretty-print roles using a Go template (default "{{.Name}} {{if .Ok}}✓{{end}}")
  -h, --help            help for list
```

### Options inherited from parent commands

```
  -q, --quiet       do not output
  -v, --verbose n   be verbose (specify --verbose multiple times or level n)
```

### SEE ALSO

* [dot](dot.md)	 - Manage files

