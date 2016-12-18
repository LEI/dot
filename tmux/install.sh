case "$OS" in
  android) apt install tmux -qqy ;;
  *linux)
    if has apt-get 2>/dev/null
    then apt-get install -qqy tmux
    elif has apt 2>/dev/null
    then apt install -qqy tmux
    fi
    ;;
esac

lineinfile present "$HOME/.tmux.conf" 'source $HOME/.tmux/tmux.conf'

directory present $HOME/.tmux
