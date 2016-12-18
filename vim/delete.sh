lineinfile absent "$HOME/.vimrc" 'source ~/.vim/init.vim'

# [[ -L "$HOME/.config/nvim" ]] && rm $HOME/.config/nvim
[[ -d "$HOME/.config/nvim" ]] && directory absent $HOME/.config/nvim

case "$OS" in
  android) apt remove -qqy vim neovim ;;
  *linux)
    if has apt-get 2>/dev/null
    then apt-get remove -qqy vim
    elif has apt 2>/dev/null
    then apt remove -qqy vim
    fi
    ;;
esac

_post_delete() {
  directory absent $HOME/.vim/{plugin,settings}
}
