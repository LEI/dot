#!/usr/bin/env bash

if hash nvim 2>/dev/null
then EDITOR="nvim"
else EDITOR="vim -f"
fi

export EDITOR
export VISUAL="$EDITOR"

export HISTSIZE="32768"
export HISTFILESIZE="${HISTSIZE}"
export HISTCONTROL="ignoreboth"
export HISTTIMEFORMAT="%F %T "

# export LESS_TERMCAP_md="${yellow}"
# export MANPAGER="less -X"
