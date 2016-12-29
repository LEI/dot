#!/usr/bin/env bash

# Deamonized: docker run -d -P <img> <cmd>
# Interactive: docker run -i -t -P <img> <cmd>

alias d="_docker"
# alias dcompose="_docker_compose"
# alias dmachine="_docker_machine"

_docker() {
  local cmd="$1"
  shift
  case "$cmd" in
    ''|a|pa|psa) docker ps --all ;; # List all containers
    all) _docker_all "$@" ;;
    *-all) _docker_all "${cmd%-all}" "$@" ;;
    b) docker build -t "$1" "${2:-.}" "${@:3}" ;;
    bash|sh) _docker_exec "$1" "$cmd" "${@:2}" ;;
    c|compose) _docker_compose "$@" ;;
    clean) _docker_clean "$@" ;;
    dangling) docker images --all --quiet --filter "dangling=${1:-true}" "${@:2}" ;;
    env) _docker_env "$@" ;; # env | grep DOCKER_
    i) docker images "$@" ;;
    id) docker ps --all --quiet --filter "name=$1" "${@:2}" ;;
    ip) _docker_ip "$@" ;;
    l) docker logs --follow --timestamps "$@" ;; # --since, --tail=all
    last) docker ps -l --quiet "$@" ;; # Latest container ID
    m|machine) _docker_machine "$@" ;;
    p) docker pull "$@" ;; # --all-tags
    r) docker run -i -t -v "$(pwd):${1:-/app}" -w "${1:-/app}" --rm "${@:2}" ;; # --name ? -p 80:80
    *) docker "$cmd" "$@" ;;
  esac
}

_docker_all() {
  local cmd="$1"
  shift
  case "$cmd" in
    ''|ps) docker ps --all ;;
    rmi) docker rmi $(docker images --quiet) "$@" ;;
    rm|start|stop|*) docker $cmd $(docker ps --all --quiet) "$@" ;;
  esac
}

_docker_clean() {
  local cmd="$1"
  shift
  case "$cmd" in # d i | awk '/<none>/ {print $3}/'
    ''|images) local d="$(_docker dangling)"; if [[ -n "$d" ]]; then docker rmi $d "$@"; fi ;;
    exited) docker rm $(docker ps --all | awk '/Exited \([0-9]+\)/ {print $1}') ;;
  esac
}

_docker_compose() {
  local cmd="$1"
  shift
  case "$cmd" in
    '') docker-compose ps ;;
    l) docker-compose logs --follow --timestamps ;; # --tail=all
    u) docker-compose up -d ;; # --{force,no}-recreate --{,no-}build
    *) docker-compose "$cmd" "$@" ;;
  esac
}

_docker_exec() {
  # local c="$(_docker id "$1" || _docker last)"
  local c="${1:-$(_docker last)}"
  [[ -n "$c" ]] && docker exec -i -t "$c" "${2:-bash}" "${@:3}"
}

_docker_ip() {
  local c="${1:-$(_docker last)}"
  [[ -n "$c" ]] && docker inspect --format "{{ .NetworkSettings.IPAddress }}" "$c" "${@:2}"
}

_docker_machine() {
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

_docker_env() {
  if [[ $# -ne 0 ]]
  then _docker_machine env "$@"
  else
    local v
    for v in "${!DOCKER_@}"
    do printf "%s=\"%s\"\n" "$v" "${!v}"
    done
  fi
}
