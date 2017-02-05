#!/usr/bin/env bash

log() {
  printf "%s\n" "$*"
}

err() {
  local ret=$?
  >&2 log "$*"
  return $ret
}

usage_exit() {
  local ret="${1:-$?}"
  shift
  >&2 log "$*"
  exit $ret
}

confirm() {
  printf "%s" "$1 [y/N] "
  read -n 1 REPLY
  printf "\n"
  [[ "$REPLY" =~ ^[Yy]$ ]]
}

run() {
  [[ "$VERBOSE" -gt 1 ]] && log "$*"
  if [[ "$DRY_RUN" -eq 0 ]]
  then "$@"
  fi
}
