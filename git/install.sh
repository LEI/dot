case "$OS" in
  android)
    has git || apt install -qq -y git
    template "$BOOTSTRAP/git/.gitconfig.local.tpl" \
      "$HOME/.gitconfig.local" \
      "GIT_NAME:What is your github full name?" \
      "GIT_USERNAME:What is your github username?" \
      "GIT_EMAIL:What is your github email?"
    ;;
esac
