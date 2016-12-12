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

#\u@\h
PS1='\w\$ '

for f in {_aliases,_functions,rc.local}
do f="$HOME/.bash$f"
  if [[ -f "$f" ]] || [[ -L "$f" ]]
  then source "$f"
  fi
done
unset f
