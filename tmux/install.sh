case "$OS" in
  android) apt install tmux -qq -y ;;
  linux)
    if has apk
    then apk add -q tmux
    elif has apt-get
    then apt-get install -y tmux
    fi
    ;;
esac

for p in $HOME/{.tmux}
do [[ -d "$p" ]] || mkdir -p "$p"
done
