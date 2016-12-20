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

run() {
  [[ "${verbose:-0}" -ne 0 ]] && log "$*"
  [[ -n "${RUN:-}" ]] && [[ "$RUN" -ne 0 ]] && "$@"
}

nchar() {
  local char="$1"
  local nb="$2"
  local max="$3"
  local i
  [[ -n "$max" ]] && [[ "$nb" -gt "$max" ]] && nb="$max"
  for i in $(seq "$nb")
  do printf "%s" "$char"
  done
}
