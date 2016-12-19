source $1/lib/{directory,lineinfile}.bash

lineinfile absent "$HOME/.tmux.conf" 'source $HOME/.tmux/tmux.conf'

directory absent $HOME/.tmux/plugins
