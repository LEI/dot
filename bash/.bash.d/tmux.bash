#!/usr/bin/env bash

t() {
  if [[ $# -ne 0 ]]
  then tmux "$@"
  elif [[ -n "$TMUX" ]]
  then tmux new-session -d
  else tmux attach || tmux new-session
  fi
}

if hash _tmux 2>/dev/null
then complete -o default -o nospace -F _tmux t
fi
