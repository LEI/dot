#!/usr/bin/env bash

LIB="$(dirname "$BASH_SOURCE")"
source $LIB/utils.bash

directory() {
  local state="$1"
  shift
  local dir cmd=
  for dir in "$@"
  do
    case "$state" in
      present)
        if [[ ! -d "$dir" ]]
        then [[ -n "${RUN:-}" ]] && [[ "$RUN" -ne 0 ]] \
          && mkdir -p "$dir"
          log "directory: create $dir"
        fi
        ;;
      absent)
        if [[ -d "$dir" ]]
        then [[ -n "${RUN:-}" ]] && [[ "$RUN" -ne 0 ]] \
          && rmdir "$dir"
          log "directory: remove $dir"
        fi
        ;;
    esac
  done
}
