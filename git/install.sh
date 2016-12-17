case "$OS" in
  android)
    has git || apt install git
    if [[ ! -e "$HOME/.gitconfig.local" ]]
    then > "$HOME/.gitconfig.local" \
        template "$BOOTSTRAP/git/.gitconfig.local.tpl" \
        "GIT_NAME:What is your github full name?" \
        "GIT_USERNAME:What is your github username?" \
        "GIT_EMAIL:What is your github email?"
    fi
    ;;
esac
