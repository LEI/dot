for f in $1/lib/{directory,lineinfile}.bash; do source "$f"; done

lineinfile present "$HOME/.tmux.conf" 'source $HOME/.tmux/tmux.conf'

directory present $HOME/.tmux
