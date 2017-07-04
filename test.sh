#!/bin/bash

set -e

sep() { printf %${COLUMNS:-100}s |tr " " "${1:-=}"; printf "\n"; }
log() { sep "-"; printf "\n\t%s\n\n" "$@"; sep "-"; }
run() { log "\$ $*"; "$@" || exit $?; }

run dot -u "https://github.com/LEI/dot-git" -d "$HOME/.dot/git"
run dot -u "https://github.com/LEI/dot-tmux" -d "$HOME/.dot/tmux"
run dot -u "https://github.com/LEI/dot-vim" -d "$HOME/.dot/vim"

run tmux -2 -u new-session -n test "vim -E -s -u $HOME/.vim/vimrc +PlugInstall +qall; exit";

ls -la ~

# tail_bashrc="$(tail -n1 ~/.bashrc)"
# yes | run dot -s $DOT --https
# run tmux -2 -u new-session -n test "vim -E -s -u $HOME/.vimrc +PlugInstall +qall; exit"
# for f in "$HOME"/.gitconfig; do run test -f "$f"; done
# for d in "$HOME"/{.tmux/plugins/tpm,.vim/plugged}; do run test -d "$d"; done
# [[ "$(tail -n1 ~/.bashrc)" != "$tail_bashrc" ]] || exit 1
# yes | run dot remove -c $DOT/.dotrc.yml
# yes | run dot rm --https --empty;
# [[ "$(tail -n1 ~/.bashrc)" == "$tail_bashrc" ]] || exit 1
# touch ~/{.bashrc,.vim/init.vim}
# yes | run dot install -d -f bash,vim -c $DOT/.dotrc.yml
# # for d in $HOME/.dot/*; do yes | run dot "${d##*/}"; done'
