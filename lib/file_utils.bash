#!/usr/bin/env bash

backup_file() {
  local file path
  for file in "$@"
  do path="$file.backup"
    [[ ! -e "$file" ]] && return 0
    if confirm "$file already exists, backup to ‘$path’?"
    then run mv -v "$file" "$path"
    fi

  done
}

link_file() {
  local src="$1"
  local dst="$2" # {dst%/*}/
  local check_mark="✓"
  if [[ ! -e "$src" ]]
  then err "$dst: no such source"; return 1
  fi
  # echo "LINK_FILE $src -> $dst"
  if [[ -e "$dst" ]] || [[ -L "$dst" ]]
  then local link
    link="$(readlink "$dst")"
    if [[ "$link" != "$src" ]]
    then backup_file "$dst"
      if [[ -e "$dst" ]]
      then err "$dst != $src"; return 1
      elif [[ -L "$dst" ]] && [[ ! -e "$link" ]]
      then log "$dst is a broken link, removing"; rm "$dst"
      fi
      do_link "$src" "$dst" && log "$check_mark $dst => $src"
    else log "$check_mark $dst == $src"
    fi
  else do_link "$src" "$dst" && log "$check_mark $dst -> $src"
  fi
}

do_link() {
  local src="$1"
  local dst="$2"
  run ln -s "$src" "$dst"
}

remove_link() {
  local src="$1"
  local dst="$2"
  if [[ -e "$dst" ]] && [[ "$(readlink "$dst")" == "$src" ]]
  then run rm "$dst" \
    && log "Removed ‘$dst‘"
  fi
}
