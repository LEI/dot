#!/usr/bin/env bash

load() {
  local f
  for f in "$@"
  do f="$ROOT/lib/$f.bash"
    [[ -f "$f" ]] && source "$f"
  done
}

log() {
  # printf "%s" "${0##*/}: "
  printf "%s\n" "$@"
}

err() {
  local ret=$?
  >&2 log "$@"
  return $ret
}

has() {
  if ! hash "$1" 2>/dev/null
  then err "$1: command not found"
    return 127
  fi
}
