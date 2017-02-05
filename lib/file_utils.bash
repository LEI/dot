#!/usr/bin/env bash

backup_file() {
  if [[ ! -e "$dst" ]]
  then return 0
  fi
  if confirm "$src destination already exists, backup to ‘$dst.backup’?"
  then run mv -v "$dst" "$dst.backup"
  fi
}

link_file() {
  local src="$1"
  local dst="$2" # {dst%/*}/
  if [[ -e "$dst" ]]
  then err "$dst: no such file or directory"; return 1
  fi
  # echo "LINK_FILE $src -> $dst"
  if [[ -e "$dst" ]] || [[ -L "$dst" ]]
  then
    if [[ "$(readlink "$dst")" != "$src" ]]
    then
      backup_file "$dst"
      if [[ -e "$dst" ]]
      then >&2 log "$dst != $src"
      else do_link "$src" "$dst" \
        && log "$dst => $src"
      fi
    else log "$dst == $src"
    fi
  else do_link "$src" "$dst" \
    && log "$dst -> $src"
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

create_dir() {
  local opts="--parents"
  if [[ "$VERBOSE" -ne 0 ]]
  then opts+=" --verbose"
  fi
  local d
  for d in "$@"
  do run mkdir $opts "$d"
  done
}

remove_dir() {
  local opts="--ignore-fail-on-non-empty"
  if [[ "$VERBOSE" -ne 0 ]]
  then opts+=" --verbose"
  fi
  local d
  for d in "$@"
  do
    if [[ -d "$d" ]]
    then run rmdir $opts "$1"
    fi
  done
}

directory() {
  local state="$1"
  shift
  local dir
  for dir in "$@"
  do
    case "$state" in
      $state_install) [[ -d "$dir" ]] || run mkdir -p "$dir" ;;
      $state_remove) [[ -d "$dir" ]] && run rmdir "$dir" ;;
    esac
  done
}

line_in_file() {
  local state="$1"
  local file="$2"
  local line="$3"

  if ! hash fgrep 2>/dev/null
  then err "fgrep: command not found"; return 1
  fi

  case "$state" in
    $state_install)
      if [[ -z "$(fgrep -lx "$line" "$file" 2>/dev/null)" ]]
      then [[ "$VERBOSE" -gt 1 ]] && log "$line >> $file"
        if [[ "$DRY_RUN" -eq 0 ]]
        then printf "%s\n" "$line" >> "$file"
        fi
      fi
      ;;
    $state_remove)
      if [[ -z "$(fgrep -Lx "$line" "$file" 2>/dev/null)" ]]
      then local tmp="/tmp/${file##*/}.grep"
        [[ "$VERBOSE" -gt 1 ]] && log "grep -v \'${line}\' "$file" > "$tmp" && mv "$tmp" "$file""
        if [[ "$DRY_RUN" -eq 0 ]]
        then eval grep -v \'${line}\' "$file" > "$tmp" && mv "$tmp" "$file"
        fi # eval sed --in-place \'/${line//\//\\\/}/d\' "$file"
      fi
      ;;
  esac
}
