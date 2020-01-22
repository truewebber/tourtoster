FROM ubuntu:18.04

WORKDIR /work

RUN apt-get update
RUN apt-get install -y build-essential software-properties-common curl
RUN add-apt-repository ppa:longsleep/golang-backports
RUN apt-get update
RUN apt-get install -y golang-go
RUN go version

RUN curl -sL https://deb.nodesource.com/setup_13.x | bash -
RUN apt-get install -y nodejs
RUN node -v
RUN npm -v
