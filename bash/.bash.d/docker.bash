#!/usr/bin/env bash

d() {
  local cmd="$1"
  shift
  case "$cmd" in
    ''|a|pa|psa) docker ps --all ;; # List all containers
    *-all) d_all "${cmd%-all}" "$@" ;;
    b) docker build -t "$1" "${2:-.}" "${@:3}" ;;
    bash|sh) d_exec "${1:-$(d last)}" "$cmd" "${@:2}" ;;
    c|compose) d_compose "$@" ;;
    clean) d_clean "$@" ;;
    dangling) docker images --all --quiet --filter "dangling=${1:-true}" "${@:2}" ;;
    e) d_exec "${1:-$(d last)}" "$2" "${@:3}" ;;
    e:*) d_exec "${1:-$(d last)}" "${cmd#e:}" "${@:2}" ;;
    env) d_env "$@" ;; # env | grep DOCKER_
    i|img) docker images "$@" ;; # --format "table {{.ID}}\t{{.Repository}}\t{{.Tag}}\t{{.CreatedSince}}\t{{.Size}}"
    id) docker ps --all --quiet --filter "name=$1" "${@:2}" ;;
    ip) d_ip "${1:-$(d last)}" "${@:2}" ;;
    l) docker logs --follow --timestamps "$@" ;; # --since, --tail=all
    last) docker ps -l --quiet "$@" ;; # Latest container ID
    m|machine) d_machine "$@" ;;
    p) docker pull "$@" ;; # --all-tags
    r) d_run "$@" ;; # "$c" "${2:-/app}" "${3:-$PWD}"
    *) docker "$cmd" "$@" ;;
  esac
}

d_all() {
  local cmd="$1"
  shift
  case "$cmd" in
    ''|ps) docker ps --all ;;
    rmi) docker rmi $(docker images --quiet) "$@" ;;
    rm|start|stop|*) docker $cmd $(docker ps --all --quiet) "$@" ;;
  esac
}

d_clean() {
  local cmd="$1"
  shift
  case "$cmd" in # d i | awk '/<none>/ {print $3}/'
    ''|images) local d="$(d dangling)"; if [[ -n "$d" ]]; then docker rmi $d "$@"; fi ;;
    exited) docker rm $(docker ps --all | awk '/Exited \([0-9]+\)/ {print $1}') ;;
  esac
}

d_compose() {
  local cmd="$1"
  shift
  case "$cmd" in
    '') docker-compose ps ;;
    l) docker-compose logs --follow --timestamps ;; # --tail=all
    u) docker-compose up -d ;; # --{force,no}-recreate --{,no-}build
    *) docker-compose "$cmd" "$@" ;;
  esac
}

d_env() {
  if [[ $# -ne 0 ]]
  then d_machine env "$@"
  else
    local v
    for v in "${!DOCKER_@}"
    do printf "%s=\"%s\"\n" "$v" "${!v}"
    done
  fi
}

d_exec() {
  # local c="$(d id "$1" || d last)"
  [[ -n "$1" ]] && docker exec --interactive --tty \
    "$1" "${2:-bash}" "${@:3}"
}

d_ip() {
  [[ -n "$1" ]] && docker inspect \
    --format "{{.NetworkSettings.IPAddress}}" \
    "$1" "${@:2}"
}

d_machine() {
  # if ! hash docker-machine 2>/dev/null
  # then return 127
  # fi
  local cmd="$1"
  shift
  case "$cmd" in
    '') docker-machine ls ;;
    c) docker-machine create --driver "${1:-virtualbox}" "${@:2}" ;;
    e) docker-machine env "$@" ;;
    eval) eval "$(docker-machine env "${1:-default}")" "${@:2}" ;;
    rs) docker-machine restart ;;
    s) docker-machine status ;;
    *) docker-machine "$cmd" "$@" ;;
  esac
}

# Deamonized: docker run -d -P <img> <cmd>
# Interactive: docker run -i -t -P <img> <cmd>

# d_run <img> <vol> [<dir>] [args..]
d_run() {
  local i="$1" # Image
  shift
  local v="${1:-/app}"
  shift
  local d="${1:-$(pwd)}"
  shift
  [[ -n "$i" ]] && [[ -n "$v" ]] \
    && docker run \
      --interactive --tty --rm \
      --volume "$d:$v" \
      --workdir "$v" \
      "$i" "$@"
      # --name "${1:-${v##*/}}"
      # --user $(id -u):$(id -g)
      # --publish 80:80
}

if hash _docker 2>/dev/null
then complete -F _docker d
fi

if hash _docker-machine 2>/dev/null
then complete -F _docker-machine _d_machine
fi
