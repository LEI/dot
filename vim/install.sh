packages="vim"
case "$OSTYPE" in
  linux-android) packages+="neovim" ;;
esac

# if !filereadable(expand("~/.vim/init.vim"))
lineinfile present "$HOME/.vimrc" 'source ~/.vim/init.vim'

directory present $HOME/.vim/{plugin,settings}

if has nvim 2>/dev/null
then directory present $HOME/.config
  [[ -d "$HOME/.vim" ]] && ln -s $HOME/{.vim,.config/nvim}
fi
