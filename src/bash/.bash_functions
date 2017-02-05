#!/usr/bin/env bash

e() {
  if [[ -z "$EDITOR" ]]
  then >&2 printf "%s" "EDITOR is undefined"
    return 1
  fi
  if [[ $# -ne 0 ]]
  then $EDITOR "$@"
  else $EDITOR .
  fi
}

mkd() {
  [[ $# -ne 1 ]] && return 1
  mkdir -p "$1"
}

mcd() {
  [[ $# -ne 1 ]] && return 1
  mkd "$1" && cd "$_"
}

# Append or prepend to PATH
if ! hash pathmunge 2>/dev/null; then
  pathmunge() {
    local p="$1"
    local pos="$2"
    if [[ ! -d "$p" ]]; then
      # Not a directory
      return 1
    fi
    if [[ $PATH =~ (^|:)$p($|:) ]]; then
      # Already in PATH
      return
    fi
    if [[ "$pos" == "after" ]]; then
      # pathmunge /path/to/dir after
      PATH=$PATH:$p
    else
      # pathmunge /path/to/dir
      PATH=$p:$PATH
    fi
  }
fi

to() {
  case "$1" in
    lower) tr "[:upper:]" "[:lower:]" ;;
    upper) tr "[:lower:]" "[:upper:]" ;;
    *) >&2 printf "%s\n" "to: $1: illegal option"; return 1 ;;
  esac
}

tre() { # hash tree 2>/dev/null
  tree -aC -I '.git|node_modules|bower_components' --dirsfirst "$@" | less -FRNX
}

for method in GET HEAD POST PUT DELETE TRACE OPTIONS; do
  alias "$method"="lwp-request -m '$method'"
done
unset method
