#!/usr/bin/env bash

add_line_in_file() {
  local file="$1"
  local line="$2"
  # if [[ -z "$(fgrep -lx "$line" "$file" 2>/dev/null)" ]]
  if ! fgrep --files-with-matches --line-regexp --quiet "$line" "$file"
  then [[ "$VERBOSE" -gt 1 ]] && log "$line >> $file"
    if [[ "$DRY_RUN" -eq 0 ]]
    then printf "%s\n" "$line" >> "$file"
    fi
  fi
}

remove_line_in_file() {
  local file="$1"
  local line="$2"
  # [[ -z "$(fgrep -Lx "$line" "$file" 2>/dev/null)" ]]
  if ! fgrep --files-without-matches --line-regexp --quiet "$line" "$file"
  then local tmp="/tmp/$$.${file##*/}.grep"
    line="${line//\[/\\[}"
    line="${line//\]/\\]}"
    [[ "$VERBOSE" -gt 1 ]] && log "grep -v $line $file > $tmp && mv $tmp $file"
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

  if [[ -L "$file" ]] && ! confirm "$file is a symlink, abort?"
  then return 1
  fi

  local line
  for line in "$@"
  do case "$state" in
      $state_install) add_line_in_file "$file" "$line" ;;
      $state_remove) remove_line_in_file "$file" "$line" ;;
    esac
  done
}
