#!/usr/bin/env bash

LIB="$(dirname "$BASH_SOURCE")"
source $LIB/utils.bash

find_stow() {
  local action="$1"
  local target="$2"
  local path="$3"
  local opts="${STOW_OPTS:-}" # --verbose
  [[ "$verbose" -ne 0 ]] && opts+="-$(nchar "v" $verbose)"
  case "$action" in
    # install) --stow ;;
    delete) opts+=" --delete" ;;
  esac
  local pkgpath="$path/.pkg"
  local p="$(basename $path)"
  local directory="$(dirname $path)"
  local name="$p"
  [[ "${directory##*/}" != "${ROOT##*/}" ]] && name="${directory##*/}/$p"

  >&2 log "" "$name: $action..."

  unset packages _post_$action
  if [[ -f "$pkgpath/$action" ]]
  then source "$pkgpath/$action" "$ROOT" || return 1
  fi
  if [[ -f "$pkgpath/packages" ]]
  then source "$pkgpath/packages" "$ROOT" \
    && [[ -n "$packages" ]] && pkg_$action $packages \
    || return 1
  fi
  run stow $opts --ignore='.*.tpl' --ignore='.pkg' \
    --dir "$directory" --target "$target" "$p" \
    || return 1
  if hash _post_$action 2>/dev/null
  then _post_$action "$ROOT" && unset _post_$action || return 1
  fi
}
