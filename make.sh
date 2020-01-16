#!/usr/bin/env bash

DOCKER_HOST=${DOCKER_HOST:-}
DOCKER_TLS_VERIFY=${DOCKER_TLS_VERIFY:-0}

function build_all() {
  TlsAppendArg=""
  if [ "$DOCKER_TLS_VERIFY" == "1" ]; then
    TlsAppendArg="--tls"
  fi

  docker $TlsAppendArg run --rm --net=host -v $(pwd):/project -w /project \
    admin-docker.artifactory.semrush.net/golang:1.13 \
    /bin/bash -c "cd application && go build -o ../bin/${1}_lx ./cmd/${1}"

  check_exit $?
}

function check_exit() {
  if [ $1 -ne 0 ]; then
    exit 1
  fi
}

case "$1" in
server)
  build_all "server"
  ;;
*)
  echo $"Usage: $0 {action}"
  echo "Actions: "
  echo "	- server"
  exit 1
  ;;
esac
