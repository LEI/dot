#!/usr/bin/env bash

alias sudo="sudo "

# alias df="df -ah" # -T --total # pydf
# alias du="du -ach | sort -h" # ncdu
alias cp="cp -v" # -i
alias ln="ln -v" # -i
alias mv="mv -v" # -i
alias rm="rm -v"
alias wget="wget -c"

alias grep="grep --color=auto"
alias fgrep="fgrep --color=auto"
alias egrep="egrep --color=auto"

alias h="history"
alias hgrep="history | grep"
alias j="jobs"
alias l="ls -lF"
alias la="ls -lAF"
alias rd="rmdir"

alias ipecho="curl ipecho.net/plain; echo"
alias map="xargs -n1"
alias urlencode="python -c 'import sys, urllib as ul; print ul.quote_plus(sys.argv[1]);'"

alias path='echo $PATH | tr -s ":" "\n"'
alias reload="source ~/.bashrc"

alias ..="cd .."
alias ...="cd ../.."
alias ....="cd ../../.."
alias .....="cd ../../../.."
alias -- -="cd -"
