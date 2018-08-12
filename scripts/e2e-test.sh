#!/bin/sh

set -e

if [ -z "$DOT" ]; then
  echo "DOT is undefined"
  exit 1
fi

main() {
  if [ -n "$DOT" ] && [ ! -f "$HOME/.dotrc.yml" ]; then
    ln -sf "$DOT/.dotrc.yml" "$HOME/.dotrc.yml"
  fi

  # # tail_bashrc="$(tail -n1 ~/.bashrc)"
  yes | run dot sync --verbose
  yes | run dot install --verbose
  # # run tmux -2 -u new-session -n test "vim -E -s -u $HOME/.vimrc +Install +qall; exit"
  for f in "$HOME"/.gitconfig; do run test -f "$f"; done
  # # for d in "$HOME"/{.tmux/plugins/tpm,.vim/pack/config}; do run test -d "$d"; done
  # [ "$(tail -n1 ~/.bashrc)" != "$tail_bashrc" ] || exit 1
  # yes | run dot remove # "${dot_args[@]}"
  # # yes | run dot rm --https --empty;
  # [ "$(tail -n1 ~/.bashrc)" = "$tail_bashrc" ] || exit 1
  # # touch ~/{.bashrc,.vim/init.vim}
  # # yes | run dot install -s -f bash,vim -c $DOT/.dotrc.yml
  # # # for d in $HOME/.dot/*; do yes | run dot "${d##*/}"; done'
}

sep() {
  d=80;
  c="${COLUMNS:-$d}";
  [ "$c" -gt 80 ] && d=$d;
  printf "%${c}s" | tr " " "${1:-=}"
  echo
}

log() {
  sep "-"
  echo
  printf '\t%s\n' "$@"
  echo
  sep "-"
}

run() {
  log "\$ $*"
  "$@" || exit $?
}

main "$@" || exit 1
