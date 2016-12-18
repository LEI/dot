packages="vim"
case "$OSTYPE" in
  linux-android) packages+="neovim" ;;
esac

lineinfile absent "$HOME/.vimrc" 'source ~/.vim/init.vim'

_post_delete() {
  directory absent $HOME/.vim/{plugin,settings} $HOME/{.vim,.config/nvim}
}
