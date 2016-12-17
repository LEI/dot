case "$OS" in
  android)
    apt install git
    if [[ -e "$HOME/.gitconfig.local" ]]
    then template "$BOOTSTRAP/git/.gitconfig.local.tpl" \
        "GIT_NAME:What is your github full name?" \
        "GIT_USERNAME:What is your github username?" \
        "GIT_EMAIL:What is your github email?" # > "$HOME/.gitconfig.local"
    fi
    ;;
esac
