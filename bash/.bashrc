export EDITOR="vim -f"

#\u@\h
PS1='\w\$ '

alias termux-stow="stow -d $HOME/storage/shared/termux-config -t $HOME"
alias termux-git="git --git-dir=$HOME/termux-config.git --work-tree=$HOME/storage/shared/termux-config"

alias g="git"
alias la="ls -la"
alias mkd="mkdir -p"

command_not_found_handle() {
  $PREFIX/libexec/termux/command-not-found "$1"
}

e() {
  if [[ -z "$EDITOR" ]]
  then
    >&2 printf "%s" "EDITOR is undefined"
    return 1
  fi
  if [[ $# -ne 0 ]]
  then $EDITOR "$@"
  else $EDITOR
  fi
}

t() {
  if [[ $# -ne 0 ]]
  then tmux "$@"
  elif [[ -n "$TMUX" ]]
  then tmux new-session -d
  else tmux attach || tmux new-session
  fi
}

reload() {
  source ~/.bashrc
}

if [[ -r ~/.bashrc.local ]]
then source ~/.bashrc.local
fi
