case "$OS" in
  android) apt remove tmux -qqy ;;
  *linux)
    if has apt-get 2>/dev/null
    then apt-get remove -qqy tmux
    elif has apt 2>/dev/null
    then apt remove -qqy tmux
    fi
    ;;
esac

lineinfile absent "$HOME/.tmux.conf" 'source $HOME/.tmux/tmux.conf'

directory absent $HOME/.tmux/plugins
