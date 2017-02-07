#!/usr/bin/env bash

create_dir() {
  local opts=("--parents")
  if [[ "$VERBOSE" -ne 0 ]]
  then opts+=("--verbose")
  fi
  local d
  for d in "$@"
  do run mkdir "${opts[@]}" "$d"
  done
}

remove_dir() {
  local opts=("--ignore-fail-on-non-empty")
  if [[ "$VERBOSE" -ne 0 ]]
  then opts+=("--verbose")
  fi
  local d
  for d in "$@"
  do
    if [[ -d "$d" ]]
    then run rmdir "${opts[@]}" "$1"
    fi
  done
}

directory() {
  local state="$1"
  shift
  local d
  for d in "$@"
  do
    case "$state" in
      $state_install) [[ -d "$d" ]] || run mkdir -p "$d" ;;
      $state_remove) [[ -d "$d" ]] && run rmdir "$d" ;;
    esac
  done
}
