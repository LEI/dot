_install() {
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
  append "$HOME/.vimrc" 'source ~/.vim/init.vim'

  mkdirs $HOME/.vim/{plugin,settings}

  if has nvim
  then
    mkdirs $HOME/.config
    # [[ -d "$HOME/.vim" ]]
    ln -s $HOME/{.vim,.config/nvim}
    # [[ -f "$HOME/.vimrc" ]] && ln -s $HOME/{.vimrc,.config/nvim/init.vim}
  fi
}

_delete() {
  erase "$HOME/.vimrc" 'source ~/.vim/init.vim'

  if has nvim
  then
    [[ -L "$HOME/.config/nvim" ]] && rm $HOME/.config/nvim
    rmdirs $HOME/.config/nvim
  fi

  case "$OS" in
    android) apt remove -qqy vim neovim ;;
    *linux)
      if has apt-get 2>/dev/null
      then apt-get remove -qqy vim
      elif has apt 2>/dev/null
      then apt remove -qqy vim
      fi
      ;;
  esac
}

_post_delete() {
  rmdirs $HOME/.vim/{plugin,settings}
}
