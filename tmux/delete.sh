for f in $1/lib/{directory,lineinfile}.bash; do source "$f"; done

lineinfile absent "$HOME/.tmux.conf" "source $HOME/.tmux/tmux.conf"

directory absent $HOME/.tmux/plugins
