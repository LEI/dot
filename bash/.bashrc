# ~/.bashrc

# Work tree path
TERMUX_CFG="$HOME/storage/shared/termux-config"
# Git directory
TERMUX_GIT="$HOME/termux-config.git"

# \u@\h
PS1='\w\$ '

for option in autocd cdspell checkwinsize extglob globstar histappend nocaseglob
do shopt -s "$option" 2> /dev/null
done
unset option

for f in $HOME/.bash{_aliases,_exports,_functions,rc.local}
do [[ -r "$f" ]] && [[ -f "$f" ]] && source "$f" # || [[ -L "$f" ]]
done
unset f
