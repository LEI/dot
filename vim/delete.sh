for f in $1/lib/{directory,lineinfile}.bash; do source "$f"; done

lineinfile absent "$HOME/.vimrc" "source ~/.vim/init.vim"

_post_delete() {
  directory absent $HOME/.vim/{plugin,settings} $HOME/{.vim,.config/nvim}
}
