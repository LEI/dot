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
  then ln -s $HOME/{.vim,.config/nvim}
    ln -s $HOME/{.vimrc,.config/nvim/init.vim}
  fi
}
