source $1/lib/{directory,lineinfile}.bash

lineinfile present "$HOME/.tmux.conf" 'source $HOME/.tmux/tmux.conf'

directory present $HOME/.tmux
