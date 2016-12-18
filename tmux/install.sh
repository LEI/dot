case "$OS" in
  android) apt install tmux -qq -y ;;
  *linux)
    if has apk 2>/dev/null
    then apk add -q tmux
    elif has apt-get 2>/dev/null
    then apt-get install -qq -y tmux
    fi
    ;;
esac

create_dirs $HOME/.tmux

append "$HOME/.tmux.conf" 'source $HOME/.tmux/tmux.conf'
