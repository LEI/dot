case "$OS" in
  android) apt install -qq -y vim neovim ;;
  linux)
    if has apk
    then apk add -q vim
    elif has apt-get
    then apt-get install -y vim
    fi
    ;;
esac

for p in $HOME/.vim
do [[ -d "$p" ]] || mkdir -p "$p"
done

if has nvim
then
  for p in $HOME/.config/nvim
  do [[ -d "$p" ]] || mkdir -p "$p"
  done
fi
