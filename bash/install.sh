directory present $HOME/bin

lineinfile present "$HOME/.bashrc" '[[ -n "$PS1" ]] && [[ -f ~/.bash_profile ]] && source ~/.bash_profile'
