_install() {
  bash_pkg="bash bash-completion tree"

  case "$OS" in
    android) apt install -qq -y $bash_pkg ;;
    *linux)
      if has apk 2>/dev/null
      then apk add -q $bash_pkg
      elif has apt-get 2>/dev/null
      then apt-get install -y $bash_pkg
      fi
      ;;
  esac

  mkdirs $HOME/bin

  append "$HOME/.bashrc" '[[ -n "$PS1" ]] && [[ -f ~/.bash_profle ]] && source ~/.bash_profile'
}

_delete() {
  erase "$HOME/.bashrc" '[[ -n "$PS1" ]] && [[ -f ~/.bash_profle ]] && source ~/.bash_profile'
}
