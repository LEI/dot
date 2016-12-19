for f in $1/lib/{directory,lineinfile}.bash; do source "$f"; done

directory present $HOME/bin

lineinfile present "$HOME/.bashrc" '[[ -n "$PS1" ]] && [[ -f ~/.bash_profile ]] && source ~/.bash_profile'
