#!/usr/bin/env bash

log() {
  printf "%s\n" "$*"
}

err() {
  local ret=$?
  >&2 log "$*"
  return $ret
}

debug() {
  local v
  for v in "$@"
  do >&2 log "${!v:-undefined}"
  done
}

usage_exit() {
  local ret="${1:-$?}"
  shift
  >&2 log "$*"
  exit "$ret"
}

confirm() {
  printf "%s" "$1 (y/N) "
  [[ "${INTERACTIVE:-1}" -eq 0 ]] && return
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
