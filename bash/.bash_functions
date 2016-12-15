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
  git --git-dir="$HOME/termux-config.git" --work-tree="$TERMUX_CFG" "$@"
}

tstow() {
  stow -d "$TERMUX_CFV" -t "$HOME" "$@"
}

ttpm() {
  termux-fix-shebang ~/.tmux/plugged/**/*
}
