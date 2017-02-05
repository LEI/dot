# ~/.bash_profile

load() {
  local path
  for path in "$@"
  do
    if [[ -d "$path" ]]
    then load $path/*
    elif [[ -r "$path" ]] && [[ -f "$path" ]] # || [[ -L "$f" ]]
    then source "$path"
    # else >&2 printf "%s\n" "$path: No such file or directory"
    fi
  done
}

main() {
  # [[ -z "$PS1" ]] && PS1='\u at \h in \w\n\$ '

  for option in autocd cdspell checkwinsize extglob globstar histappend nocaseglob
  do shopt -s "$option" 2> /dev/null
  done
  unset option

  load $HOME/.bash{_aliases,_exports,_functions}
  # OS="$(uname -o 2>/dev/null || uname -s | to lower)"

  local file f
  for file in $HOME/.bash.d/*.bash
  do f="${file##*/}"
    hash ${file%.bash} 2>/dev/null && source $file
  done

  load $HOME/.bashrc.local
}

main "$@"
