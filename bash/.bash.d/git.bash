#!/usr/bin/env bash

g() {
  if [[ $# -gt 0 ]]
  then git "$@"
  else git status
  fi

}

# [[ -f /usr/local/etc/bash_completion.d/git-completion.bash ]]
if hash _git 2>/dev/null
then complete -o default -o nospace -F _git g
fi
