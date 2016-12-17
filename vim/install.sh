case "$OS" in
  android) apt install -qq -y vim ;;
esac

for p in $HOME/{.vim,.config/nvim}
do [[ -d "$p" ]] || mkdir -p "$p"
done
