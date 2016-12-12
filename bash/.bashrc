# ~/.bashrc

for option in autocd cdspell checkwinsize extglob globstar histappend nocaseglob
do shopt -s "$option" 2> /dev/null
done
unset option

HISTCONTROL=${HISTCONTROL:-erasedups}
HISTSIZE=${HISTSIZE:-10000}
HISTFILESIZE=${HISTFILESIZE:-10000}
HISTTIMEFORMAT='%F %T '

# Ctrl-D
if [[ "$IGNOREEOF" -lt 10 ]]
then IGNOREEOF=10
fi

#\u@\h
PS1='\w\$ '

export EDITOR="vim -f"

alias g="git"
alias la="ls -la"
alias mkd="mkdir -p"

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
  # bash --login
  source ~/.bashrc
}

command_not_found_handle() {
  $PREFIX/libexec/termux/command-not-found "$1"
}

termux-fix-tpm() {
  termux-fix-shebang ~/.tmux/plugged/**/*
}

termux-git() {
  git --git-dir="$HOME/termux-config.git" --work-tree="$HOME/storage/shared/termux-config" "$@"
}

termux-stow() {
  stow -d "$HOME/storage/shared/termux-config" -t "$HOME" "$@"
}

if [[ -r ~/.bashrc.local ]]
then source ~/.bashrc.local
fi
