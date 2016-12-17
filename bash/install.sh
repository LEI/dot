bash_pkg="bash bash-completion"

case "$OS" in
  android) apt install -qq -y $bash_pkg ;;
  linux)
    if has apk
    then apk add -q $bash_pkg
    elif has apt-get
    then apt-get install -y $bash_pkg
    fi
    ;;
esac

for p in $HOME/{bin}
do [[ -d "$p" ]] || mkdir -p "$p"
done

p='[[ -n "$PS1" ]] && source ~/.bash_profile'
if ! fgrep -x "$p"
then printf "%s\n" "$p" >> ~/.bashrc
fi
