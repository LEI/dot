case "$OS" in
  android) apt install -qqy git ;;
  *linux)
    if has apt-get 2>/dev/null
    then apt-get install -qqy git-core
    elif has apt 2>/dev/null
    then apt install -qqy git-core
    fi
    ;;
esac

template "$BOOTSTRAP/git/.gitconfig.local.tpl" "$HOME/.gitconfig.local" \
  "GIT_AUTHOR_NAME:What is your github full name?" \
  "GIT_AUTHOR_USERNAME:What is your github username?" \
  "GIT_AUTHOR_EMAIL:What is your github email?"
