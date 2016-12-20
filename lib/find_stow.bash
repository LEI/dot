#!/usr/bin/env bash

LIB="$(dirname "$BASH_SOURCE")"
source $LIB/utils.bash

find_names() {
  [[ $# -eq 0 ]] && return 1
  local names=("$@")
  local exp=("(")
  local i
  for i in "${!names[@]}"
  do [[ "$i" -ne 0 ]] && exp+=("-o")
    exp+=("-name" "${names[$i]}")
  done
  exp+=(")")
  printf "%s\n" "${exp[@]}"

}

find_stow() {
  local i ignore=(".git" "os_*" "lib")
  local state="$1"
  shift
  local action="install"
  case "$state" in
    present) action="install" ;;
    absent) action="delete" ;;
  esac
  local stow_opts="${STOW_OPTS:---verbose}"
  case "$action" in
    # install) --stow ;;
    delete) stow_opts+=" --delete" ;;
  esac
  local path="$ROOT"
  local target="$HOME"
  local find_args=("$path")
  [[ -d "$path/os_$OS" ]] && find_args+=("$path/os_$OS")
  find_args+=("-mindepth" "1" "-maxdepth" "1" "-type" "d")
  if [[ "${#ignore[@]}" -ne 0 ]]
  then find_args+=("!")
    while read -r i
    do find_args+=("$i")
    done < <(find_names "${ignore[@]}")
  fi
  if [[ $# -ne 0 ]]
  then find_args+=("-a")
    while read -r i
    do find_args+=("$i")
    done < <(find_names "$@")
  fi
  # local find_opts=()
  # while [[ $# -ne 0 ]]
  # do [[ -z "$find_opts" ]] && find_opts+=("(") || find_opts+=("-o")
  #   [[ -n "$1" ]] && find_opts+=("-name" "$1")
  #   [[ $# -eq 1 ]] && find_opts+=(")")
  #   shift
  # done
  # [[ "${#find_opts[@]}" -ne 0 ]] && find_args+=("-a" "${find_opts[@]}")

  if [[ "${verbose:-0}" -ne 0 ]]
  then printf "%s\n" "find ${find_args[*]}" $(find "${find_args[@]}" -print)
  fi

  local d
  while read -d '' -r d
  do local pkgpath="$d/.pkg"
    local p="$(basename $d)"
    local dir="$(dirname $d)"
    local name="$p"
    [[ "${dir##*/}" != "${ROOT##*/}" ]] && name="${dir##*/}/$p"
    log "$name: $action..."
    unset packages _post_$action
    [[ -f "$pkgpath/$action" ]] && "$pkgpath/$action" "$ROOT"
    [[ -f "$pkgpath/packages" ]] && "$pkgpath/packages" "$ROOT" \
      && [[ -n "$packages" ]] && pkg_$action $packages
    run stow $stow_opts --ignore='.*.tpl' --ignore='.pkg' --dir "$dir" --target "$target" "$p"
    hash _post_$action 2>/dev/null && _post_$action && unset _post_$action
  done < <(find "${find_args[@]}" -print0)
}

