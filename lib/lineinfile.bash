#!/usr/bin/env bash

LIB="$(dirname "$BASH_SOURCE")"
source $LIB/utils.bash

lineinfile() {
  local state="$1"
  local file="$2"
  local line="$3"
  if ! has fgrep
  then return 1
  fi
  case "$state" in
    present)
      if [[ -z "$(fgrep -lx "$line" "$file" 2>/dev/null)" ]]
      then log "line>file: $line >> $file"
        run "echo "$line" >> "$file""
      fi
      ;;
    absent)
      if [[ -z "$(fgrep -Lx "$line" "$file" 2>/dev/null)" ]]
      then log "line<file: $line << $file"
        local tmp="/tmp/${file##*/}.grep"
        # eval sed --in-place \'/${line//\//\\\/}/d\' "$file"
        run "eval grep -v \'${line}\' "$file" > "$tmp" && mv "$tmp" "$file""
      fi
      ;;
  esac
}
