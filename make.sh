#!/usr/bin/env bash

function rand() {
  LC_CTYPE=C tr -dc 'a-zA-Z0-9' </dev/urandom | fold -w 10 | head -n 1
}

function build_backend() {
  docker run --rm --net=host -v "$(pwd)":/project -w /project \
    tourtoster/builder:latest \
    /bin/bash -c "cd application && go build -o ../bin/server_lx ./cmd/server" || exit 1
}

function build_frontend() {
  docker run --rm --net=host -v "$(pwd)":/project -w /project \
    tourtoster/builder:latest \
    /bin/bash -c "cd front/tools && yarn install && yarn build" || exit 1
}

function front_watcher() {
  cd front/tools || exit 1
  yarn watch
}

function build_frontend_local() {
  cd front/tools || exit 1
  yarn install || exit 1
  yarn build || exit 1
  cd ../../ || exit 1
}

function move_static() {
  if [ ! -d "./static" ]; then
    ln -s ./front/dist/assets ./static
  fi
}

function pack() {
  mkdir -p ./build
  mkdir -p ./build/bin
  mkdir -p ./build/static

  cp -R ./front/dist/assets/css ./build/static/css
  cp -R ./front/dist/assets/js ./build/static/js
  cp -R ./front/dist/assets/media ./build/static/media
  cp -R ./front/dist/assets/plugins ./build/static/plugins

  cp -R ./templates ./build/templates
  cp ./bin/server_lx ./build/bin/server

  tar -cvzf ./tourtoster.tar.gz ./build

  rm -rf ./build

  rm -f ./bin/server_lx
  rm -rf ./front/tools/node_modules
  rm -rf ./front/dist
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
watch)
  front_watcher
  ;;
static)
  build_frontend_local
  move_static
  ;;
build)
  build_backend
  build_frontend
  ;;
pack)
  build_backend
  build_frontend
  pack
  ;;
deploy)
  build_backend
  build_frontend
  pack
  deploy
  ;;
*)
  echo $"Usage: $0 {action}"
  echo "Actions: "
  echo "	- static"
  echo "	- build"
  echo "	- pack"
  echo "	- deploy"
  exit 1
  ;;
esac
