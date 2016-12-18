case "$OS" in
  android) apt install -qqy vim neovim ;;
  *linux)
    if has apt-get 2>/dev/null
    then apt-get install -qqy vim
    elif has apt 2>/dev/null
    then apt install -qqy vim
    fi
    ;;
esac

# if !filereadable(expand("~/.vim/init.vim"))
lineinfile present "$HOME/.vimrc" 'source ~/.vim/init.vim'

directory present $HOME/.vim/{plugin,settings}

if has nvim
then
  directory present $HOME/.config
  # [[ -d "$HOME/.vim" ]]
  ln -s $HOME/{.vim,.config/nvim}
  # [[ -f "$HOME/.vimrc" ]] && ln -s $HOME/{.vimrc,.config/nvim/init.vim}
fi
