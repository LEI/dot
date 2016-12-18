_install() {
  case "$OS" in
    android) apt install -qq -y git ;;
    *linux)
      if has apk 2>/dev/null
      then apk add -q git
      elif has apt-get 2>/dev/null
      then apt-get install -y git-core
      fi
      ;;
  esac

  template "$BOOTSTRAP/git/.gitconfig.local.tpl" "$HOME/.gitconfig.local" \
    "GIT_AUTHOR_NAME:What is your github full name?" \
    "GIT_AUTHOR_USERNAME:What is your github username?" \
    "GIT_AUTHOR_EMAIL:What is your github email?"
}
