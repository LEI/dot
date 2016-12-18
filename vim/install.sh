case "$OS" in
  android) apt install -qq -y vim neovim ;;
  *linux)
    if has apk 2>/dev/null
    then apk add -q vim
    elif has apt-get 2>/dev/null
    then apt-get install -y vim
    fi
    ;;
esac

create_dirs $HOME/.vim

# if !filereadable(expand("~/.vim/init.vim"))
append "$HOME/.vimrc" 'source ~/.vim/init.vim'

postow() {
  if has nvim
  then
    create_dirs $HOME/.config
    # [[ -d "$HOME/.vim" ]]
    ln -s $HOME/{.vim,.config/nvim}
    # [[ -f "$HOME/.vimrc" ]] && ln -s $HOME/{.vimrc,.config/nvim/init.vim}
  fi
}
