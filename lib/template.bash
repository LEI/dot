#!/usr/bin/env bash

LIB="$(dirname "$BASH_SOURCE")"
source $LIB/utils.bash

template() {
  local src="$1"
  shift
  local dst="$1" # ${1:-$HOME/$(basename "${1%.tpl}")}
  shift
  local vars=("$@")

  if [[ ! -f "$src" ]] && [[ ! -L "$src" ]]
  then err "$src: not found"
    return 1
  fi

  if [[ -f "$dst" ]] || [[ -L "$dst" ]]
  then err "${dst/$HOME/~}: already exists"
    return
  fi

  local opts=()
  while [[ $# -ne 0 ]]
  do
    local var="${1%%:*}"
    local val="${!var}"
    if [[ -n "$val" ]]
    then opts+=("-e" "'s/$var/$val/g'")
    else err "$var: undefined variable"
    fi
    shift
  done

  if [[ -n "${#opts[*]}" ]]
  then run "sed "${opts[*]}" "$src" > "$dst""
  else err "template: no variables provided"; return 1
  fi
}
