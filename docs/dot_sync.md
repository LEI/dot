## dot sync

Synchronize roles

### Synopsis

The "sync" command clone or pull a role repository.

```
dot sync [flags]
```

### Options

```
  -c, --config-file DOT_FILE   main configuration DOT_FILE (default ".dotrc.yml")
  -d, --dry-run                do not execute tasks
  -F, --force                  force execution
  -h, --help                   help for sync
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

