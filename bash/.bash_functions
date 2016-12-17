#!/usr/bin/env bash

mkd() {
  mkdir -p "$@" && cd "$_"
}

tre() {
  tree -aC -I '.git|node_modules|bower_components' --dirsfirst "$@" | less -FRNX
}

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

t() {
  if [[ $# -ne 0 ]]
  then tmux "$@"
  elif [[ -n "$TMUX" ]]
  then tmux new-session -d
  else tmux attach || tmux new-session
  fi
}

reload() {
  # bash --login
  source ~/.bashrc
}

tgit() {
  git --git-dir="$TERMUX_GIT" --work-tree="$TERMUX_CFG" "$@"
}

tstow() {
  stow -d "$TERMUX_CFG" -t "$HOME" "$@"
}

ttpm() {
  termux-fix-shebang ~/.tmux/plugged/**/*
}
