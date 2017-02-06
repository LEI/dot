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
  local check_mark="✓"
  if [[ ! -e "$src" ]]
  then err "$dst: no such source"; return 1
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
        && log "$check_mark $dst => $src"
      fi
    else log "$check_mark $dst == $src"
    fi
  else do_link "$src" "$dst" \
    && log "$check_mark $dst -> $src"
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

add_line_in_file() {
  local file="$1"
  local line="$2"
  if [[ -z "$(fgrep -lx "$line" "$file" 2>/dev/null)" ]]
  then [[ "$VERBOSE" -gt 1 ]] && log "$line >> $file"
    if [[ "$DRY_RUN" -eq 0 ]]
    then printf "%s\n" "$line" >> "$file"
    fi
  fi
}

remove_line_in_file() {
  local file="$1"
  local line="$2"
  if [[ -z "$(fgrep -Lx "$line" "$file" 2>/dev/null)" ]]
  then local tmp="/tmp/$$.${file##*/}.grep"
    line="${line//\[/\\[}"
    line="${line//\]/\\]}"
    [[ "$VERBOSE" -gt 1 ]] && log "grep -v "$line" "$file" > "$tmp" && mv "$tmp" "$file""
    if [[ "$DRY_RUN" -eq 0 ]]
    then grep -v "$line" "$file" > "$tmp" && mv "$tmp" "$file"
    fi # eval sed --in-place \'/${line//\//\\\/}/d\' "$file"
  fi
}

line_in_file() {
  local state="$1"
  shift
  local file="$1"
  shift

  if ! hash fgrep 2>/dev/null
  then err "fgrep: command not found"; return 1
  fi

  local line
  for line in "$@"
  do case "$state" in
      $state_install) add_line_in_file "$file" "$line" ;;
      $state_remove) remove_line_in_file "$file" "$line" ;;
    esac
  done
}
