# ~/.bashrc

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

  case "$(uname -o 2>/dev/null)" in
    Android) OS="android" ;;
    *) OS="$(uname -s | to lower)" ;;
  esac

  load $HOME/.bash{_aliases,_exports,_functions,rc.local}
  load $HOME/.$OS.d/*.bash
  load $HOME/.bashrc.local
}

main "$@"
