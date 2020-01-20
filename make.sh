#!/usr/bin/env bash

DOCKER_TLS_VERIFY=${DOCKER_TLS_VERIFY:-0}

function rand() {
  LC_CTYPE=C tr -dc 'a-zA-Z0-9' </dev/urandom | fold -w 10 | head -n 1
}

function build() {
  TlsAppendArg=""
  if [ "$DOCKER_TLS_VERIFY" == "1" ]; then
    TlsAppendArg="--tls"
  fi

  docker $TlsAppendArg run --rm --net=host -v "$(pwd)":/project -w /project \
    admin-docker.artifactory.semrush.net/golang:1.13 \
    /bin/bash -c "cd application && go build -o ../bin/server_lx ./cmd/server" || exit 1
}

function pack() {
  mkdir -p ./build
  mkdir -p ./build/bin

  cp -R ./static ./build/static
  cp -R ./templates ./build/templates
  cp ./bin/server_lx ./build/bin/server

  tar -cvzf ./tourtoster.tar.gz ./build

  rm -rf ./build
}

function deploy() {
  scp ./tourtoster.tar.gz srv1:/tmp/ | exit 2

  temp_dir="/tmp/deploy-$(rand)"

  ssh srv1 "mkdir -p ${temp_dir}"
  ssh srv1 "tar -xzf /tmp/tourtoster.tar.gz -C ${temp_dir}"

  ssh srv1 "cp -R ${temp_dir}/build/static ~/web/tourtoster.truewebber.com/app/"
  ssh srv1 "cp -R ${temp_dir}/build/templates ~/web/tourtoster.truewebber.com/app/"
  ssh srv1 "cp ${temp_dir}/build/bin/server ~/web/tourtoster.truewebber.com/bin/"

  ssh srv1 "rm -rf ${temp_dir}"
  ssh srv1 "rm -f /tmp/tourtoster.tar.gz"
}

case "$1" in
build)
  build
  ;;
pack)
  build
  pack
  ;;
deploy)
  build
  pack
  deploy
  ;;
*)
  echo $"Usage: $0 {action}"
  echo "Actions: "
  echo "	- build"
  echo "	- pack"
  echo "	- deploy"
  exit 1
  ;;
esac
