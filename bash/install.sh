case "$OS" in
  android) apt install -qq -y bash bash-completion ;;
esac

for p in $HOME/{bin}
do [[ -d "$p" ]] || mkdir -p "$p"
done
