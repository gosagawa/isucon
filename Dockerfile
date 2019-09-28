FROM golang:latest

ENV GOPATH /go
ENV PATH $PATH:$GOPATH/bin

RUN apt-get update
RUN apt-get install -y sudo
RUN apt-get install -y vim
