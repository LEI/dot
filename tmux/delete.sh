source $1/lib/{directory,lineinfile}.bash

type lineinfile
lineinfile absent "$HOME/.tmux.conf" "source $HOME/.tmux/tmux.conf"

directory absent $HOME/.tmux/plugins
