case "$OS" in
  android) apt install tmux -qq -y ;;
esac

for p in $HOME/{.tmux}
do [[ -d "$p" ]] || mkdir -p "$p"
done
