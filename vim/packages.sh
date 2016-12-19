source $1/lib/detect_os.bash

case "$OS" in
  android|macos) packages="vim neovim" ;;
  *) packages="vim" ;;
esac
