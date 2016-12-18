case "$OSTYPE" in
  linux-android) packages="git" ;;
  *linux) packages="git-core" ;;
  *) packages="git" ;;
esac
