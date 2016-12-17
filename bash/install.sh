bash_pkg="bash bash-completion tree"

case "$OS" in
  android) apt install -qq -y $bash_pkg ;;
  *linux)
    if has apk 2>/dev/null
    then apk add -q $bash_pkg
    elif has apt-get 2>/dev/null
    then apt-get install -y $bash_pkg
    fi
    ;;
esac

for p in $HOME/bin
do [[ -d "$p" ]] || mkdir -p "$p"
done

p='[[ -n "$PS1" ]] && [[ -f ~/.bash_profle ]] && source ~/.bash_profile'
if ! fgrep -x "$p" "$HOME/.bashrc" &>/dev/null
then printf "%s\n" "" "$p" >> "$HOME/.bashrc"
fi
