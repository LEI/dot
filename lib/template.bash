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
    val="$(prompt "${1#*:} (default: $val) " "$val")"
    if [[ -n "$val" ]]
    then opts+=("-e" "'s/$var/$val/g'")
    else err "$var: empty variable"
    fi
    shift
  done

  if [[ "${#opts[@]}" -ne 0 ]]
  then
    if [[ -z "$dryrun" ]] || [[ "$dryrun" -eq 0 ]]
    then eval sed "${opts[@]}" "$src" > "$dst"
    else log "DRY-RUN: sed "${opts[@]}" "$src""
    fi
  else err "template: no options"; return 1
  fi
}
