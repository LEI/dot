# ~/.bashrc

TERMUX_CFG="$HOME/storage/shared/termux-config"

export EDITOR="vim -f"

for option in autocd cdspell checkwinsize extglob globstar histappend nocaseglob
do shopt -s "$option" 2> /dev/null
done
unset option

HISTCONTROL=${HISTCONTROL:-erasedups}
HISTSIZE=${HISTSIZE:-10000}
HISTFILESIZE=${HISTFILESIZE:-10000}
HISTTIMEFORMAT='%F %T '

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

tgit() {
  git --git-dir="$HOME/termux-config.git" --work-tree="$TERMUX_CFG" "$@"
}

tstow() {
  stow -d "$TERMUX_CFV" -t "$HOME" "$@"
}

ttpm() {
  termux-fix-shebang ~/.tmux/plugged/**/*
}

#\u@\h
PS1='\w\$ '

for f in {.bash_aliases,.bashrc.local}
do f="$HOME/$f"
  if [[ -f "$f" ]] || [[ -L "$f" ]]
  then source "$f"
  fi
done
unset f
