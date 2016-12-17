case "$OS" in
  android) apt install tmux -qq -y ;;
  linux)
    if has apk
    then apk add -q tmux
    fi
    ;;
esac

for p in $HOME/{.tmux}
do [[ -d "$p" ]] || mkdir -p "$p"
done
