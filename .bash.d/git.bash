#!/usr/bin/env bash

g() {
  if [[ $# -gt 0 ]]
  then git "$@"
  else git status
  fi
}

# [[ -f /usr/local/etc/bash_completion.d/git-completion.bash ]]
if hash _git 2>/dev/null
then complete -o bashdefault -o default -o nospace -F _git g # __git_wrap__git_main
fi
