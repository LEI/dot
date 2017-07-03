#!/bin/bash

set -e

sep() { printf %${COLUMNS:-100}s |tr " " "${1:-=}"; printf "\n"; }
log() { sep "-"; printf "\n\t%s\n\n" "$@"; sep "-"; }
run() { log "\$ $*"; "$@" || exit $?; }

#touch ~/.dot/git;
#touch ~/.gitconfig;
#touch ~/.gitconfig.local;

dot clone -u "https://github.com/LEI/dot-git" -d "$HOME/.dot/git"
# dot link -d "$HOME/.dot/git" \
# 	".gitconfig" \
# 	".gitignore"
# dot template -d "$HOME/.dot/git" \
# 	".gitconfig.local.tpl"
dot link -d "$HOME/.dot/git" - <<< '{"link": [".gitconfig", ".gitignore"]}'
dot template -d "$HOME/.dot/git" - <<< '{"template": [".gitconfig.local.tpl"]}'

dot clone -u "https://github.com/LEI/dot-tmux" -d "$HOME/.dot/tmux"
dot link -d "$HOME/.dot/tmux" \
	"tmux.conf:$HOME/.tmux" \
	"tmux-*.conf:$HOME/.tmux" \
	"tmux.\$OS.conf:$HOME/.tmux" \
	"colors:$HOME/.tmux"
if ! test -f "$HOME/.tmux.conf" || ! grep -Fxq 'source-file $HOME/.tmux/tmux.conf' "$HOME/.tmux.conf"; then
	echo 'source-file $HOME/.tmux/tmux.conf' >> "$HOME/.tmux.conf"
fi

dot clone -u "https://github.com/LEI/dot-vim" -d "$HOME/.dot/vim"
dot link -d "$HOME/.dot/vim" \
	"autoload/*:$HOME/.vim/autoload" \
	"*.vim:$HOME/.vim" \
	"config:$HOME/.vim" \
	"ftdetect:$HOME/.vim" \
	"ftplugin:$HOME/.vim" \
	"plugin:$HOME/.vim"
	#"[^.]*[^\.md]?:$HOME/.vim"
if ! test -f "$HOME/.vim/vimrc" || ! grep -Fxq 'source ~/.vim/init.vim' "$HOME/.vim/vimrc"; then
	echo 'source ~/.vim/init.vim' >> "$HOME/.vim/vimrc"
fi

run tmux -2 -u new-session -n test "vim -E -s -u $HOME/.vim/vimrc +PlugInstall +qall; exit";

ls -la ~ ~/.tmux ~/.vim

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
