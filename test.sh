#!/bin/bash

dot clone -u "https://github.com/LEI/dot-git" -d "$HOME/.dot/git";
dot clone -u "https://github.com/LEI/dot-tmux" -d "$HOME/.dot/tmux";
dot clone -u "https://github.com/LEI/dot-vim" -d "$HOME/.dot/vim";
dot link -d "$HOME/.dot/git" \
	".gitconfig:$HOME" \
	".gitignore:$HOME";
dot template -d "$HOME/.dot/git" \
	".gitconfig.local.tpl:$HOME";
dot link -d "$HOME/.dot/tmux" \
	"tmux.conf:$HOME/.tmux" \
	"tmux-*.conf:$HOME/.tmux" \
	"tmux.\$OS.conf:$HOME/.tmux" \
	"colors:$HOME/.tmux";
dot link -d "$HOME/.dot/vim" \
	"autoload/*:$HOME/.vim/autoload" \
	"*.vim:$HOME/.vim" \
	"config:$HOME/.vim" \
	"ftdetect:$HOME/.vim" \
	"ftplugin:$HOME/.vim" \
	"plugin:$HOME/.vim";
	#"[^.]*[^\.md]?:$HOME/.vim"
