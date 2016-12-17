case "$OS" in
  android) apt install -qq -y bash bash-completion ;;
  linux)
    if has apk
    then apk add -q bash bash-completion
    fi
    ;;
esac

for p in $HOME/{bin}
do [[ -d "$p" ]] || mkdir -p "$p"
done
