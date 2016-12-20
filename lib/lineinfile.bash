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
  log "line '$line' $state $file"
  case "$state" in
    present)
      if [[ -z "$(fgrep -Flx "$line" "$file" 2>/dev/null)" ]]
      then
        [[ "${verbose:-0}" -ne 0 ]] && log "echo $line >> $file"
        [[ -n "${RUN:-}" ]] && [[ "$RUN" -ne 0 ]] && echo "$line" >> "$file"
      fi
      ;;
    absent)
      if [[ -z "$(fgrep -FLx "$line" "$file" 2>/dev/null)" ]]
      then
        # [[ -f "$file" ]] && mv "$file" "$file.backup"
        line="${line//\[/\\\[}"
        line="${line//\]/\\\]}"
        run sed --in-place=.backup "/${line//\//\\\/}/d" "$file"
        # eval sed --in-place \'/${line//\//\\\/}/d\' "$file"
        # local tmp="/tmp/${file##*/}.grep"
        # run eval "grep -Fv '"${line}"' "$file" > "$tmp" && mv "$tmp" "$file""
        # [[ "${verbose:-0}" -ne 0 ]] && log "grep -Fv '$line' $file > $tmp && mv $tmp $file"
        # [[ -n "${RUN:-}" ]] && [[ "$RUN" -ne 0 ]] \
        #   && grep -Fv \'"$line"\' "$file" > "$tmp" && mv "$tmp" "$file" \
        #   || echo "failed to $RUN"
      fi
      ;;
  esac
}
