case "$OS" in
  android|macos) packages="git" ;;
  debian) packages="git-core" ;;
  *) packages="git" ;;
esac
