case "$OS" in
  android) apt remove -qqy git ;;
  *linux)
    if has apt-get 2>/dev/null
    then apt-get remove -qqy git-core
    elif has apt 2>/dev/null
    then apt remove -qqy git-core
    fi
    ;;
esac

# rm? "$HOME/.gitconfig.local"
