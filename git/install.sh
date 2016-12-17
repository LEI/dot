case "$OS" in
  android) apt install -qq -y git
    template "$BOOTSTRAP/git/.gitconfig.local.tpl" "$HOME/.gitconfig.local" \
      "GIT_NAME:What is your github full name?" \
      "GIT_USERNAME:What is your github username?" \
      "GIT_EMAIL:What is your github email?"
    ;;
  linux)
    if has apk
    then apk add -q git
    elif has apt-get
    then apt-get install -y git-core
    fi
    ;;
esac
