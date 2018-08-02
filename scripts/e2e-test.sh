#!/bin/bash

set -e

DIR="${BASH_SOURCE%/*}"

source "$DIR/functions.sh"

if [[ -n "$DOT" ]] && [[ ! -f "$HOME/.dot.yml" ]]; then
  ln -sf "$DOT/.dot.yml" "$HOME/.dot.yml"
fi

# run dot sync -u "https://github.com/LEI/dot-git" -s ~/.dot/git
# run dot install link -u "https://github.com/LEI/dot-git" -s ~/.dot/git ".gitconfig" ".gitignore"
# CREDENTIAL_HELPER=cache run dot install template -u "https://github.com/LEI/dot-git" -s ~/.dot/git ".gitconfig.local.tpl"
# run dot install line -u "https://github.com/LEI/dot-git" -s ~/.dot/git
# run dot -- git

tail_bashrc="$(tail -n1 ~/.bashrc)"
# Travis CI fix: sudo chmod +x /usr/local/bin/pacapt
run dot list

dot_args=()
if [[ "$(uname -s)" != "Darwin" ]] && [[ "$(id -u)" -ne 0 ]]; then
  dot_args+=("--sudo")
fi
yes | run dot install --packages "${dot_args[@]}" # --verbose

run tmux -2 -u new-session -n test "vim -E -s -u $HOME/.vimrc +Install +qall; exit"
for f in "$HOME"/.gitconfig; do run test -f "$f"; done
for d in "$HOME"/{.tmux/plugins/tpm,.vim/pack/config}; do run test -d "$d"; done
[[ "$(tail -n1 ~/.bashrc)" != "$tail_bashrc" ]] || exit 1
yes | run dot remove # "${dot_args[@]}"
# yes | run dot rm --https --empty;
[[ "$(tail -n1 ~/.bashrc)" == "$tail_bashrc" ]] || exit 1
# touch ~/{.bashrc,.vim/init.vim}
# yes | run dot install -s -f bash,vim -c $DOT/.dotrc.yml
# # for d in $HOME/.dot/*; do yes | run dot "${d##*/}"; done'
