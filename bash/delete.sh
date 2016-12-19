source $1/lib/{directory,lineinfile}.bash

directory absent $HOME/bin

lineinfile absent "$HOME/.bashrc" '[[ -n "$PS1" ]] && [[ -f ~/.bash_profile ]] && source ~/.bash_profile'
