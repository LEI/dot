case "$OS" in
  android) apt install -qq -y vim neovim ;;
  linux)
    if has apk
    then apk add -q vim
    fi
    ;;
esac

for p in $HOME/{.vim,.config/nvim}
do [[ -d "$p" ]] || mkdir -p "$p"
done
