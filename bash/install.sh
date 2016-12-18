bash_pkg="bash bash-completion tree"

case "$OS" in
  android) apt install -qqy $bash_pkg ;;
  *linux)
    if has apt-get 2>/dev/null
    then apt-get install -qqy $bash_pkg
    elif has apt 2>/dev/null
    then apt install -qqy $bash_pkg
    fi
    ;;
esac

directory present $HOME/bin

lineinfile present "$HOME/.bashrc" '[[ -n "$PS1" ]] && [[ -f ~/.bash_profile ]] && source ~/.bash_profile'
