bash_pkg="bash bash-completion tree"

case "$OS" in
  android) apt remove -qqy $bash_pkg ;;
  *linux)
    if has apt-get 2>/dev/null
    then apt-get remove -qqy $bash_pkg
    elif has apt 2>/dev/null
    then apt remove -qqy $bash_pkg
    fi
    ;;
esac

directory absent $HOME/bin

lineinfile absent "$HOME/.bashrc" '[[ -n "$PS1" ]] && [[ -f ~/.bash_profile ]] && source ~/.bash_profile'
