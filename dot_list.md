## dot list

List managed files

### Synopsis

The "list" command lists roles and their tasks.

```
dot list [flags]
```

### Options

```
  -c, --config-file DOT_FILE   main configuration DOT_FILE (default ".dotrc.yml")
  -d, --dry-run                do not execute tasks
  -F, --force                  force execution
      --format string          Pretty-print roles using a Go template
  -h, --help                   help for list
      --https                  use HTTPS to clone repositories
  -l, --long                   Output role tasks
  -P, --packages               manage system packages
  -r, --role-filter strings    filter roles by name
  -s, --source DOT_SOURCE      DOT_SOURCE directory (default "$HOME")
  -t, --target DOT_TARGET      DOT_TARGET directory (default "$HOME")
```

### Options inherited from parent commands

```
  -q, --quiet       do not output
  -v, --verbose n   be verbose (specify --verbose multiple times or level n)
```

### SEE ALSO

* [dot](dot.md)	 - Manage files

