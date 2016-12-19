#!/usr/bin/env bash

LIB="$(dirname "$BASH_SOURCE")"
source $LIB/utils.bash

directory() {
  local state="$1"
  shift
  local dir
  for dir in "$@"
  do
    case "$state" in
      present) [[ -d "$dir" ]] || dry_run mkdir -p "$dir" ;;
      absent) [[ -d "$dir" ]] && dry_run rmdir "$dir" ;;
    esac
  done
}
