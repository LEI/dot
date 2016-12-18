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

for p in $HOME/.vim
do [[ -d "$p" ]] || mkdir -p "$p"
done

postow() {
  if has nvim
  then
    for p in $HOME/.config
    do [[ -d "$p" ]] || mkdir -p "$p"
    done
    [[ -d "$HOME/.vim" ]] && ln -s $HOME/{.vim,.config/nvim}
    [[ -f "$HOME/.vimrc" ]] && ln -s $HOME/{.vimrc,.config/nvim/init.vim}
  fi
}
