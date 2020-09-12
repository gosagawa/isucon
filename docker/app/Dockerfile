FROM golang:latest

ENV GOPATH /go
ENV PATH $PATH:$GOPATH/bin

RUN apt-get update
RUN apt-get install -y sudo
RUN apt-get install -y vim
RUN apt-get install -y default-mysql-client
RUN apt-get install -y python-pydot python-pydot-ng graphviz

RUN go get -u bitbucket.org/liamstask/goose/cmd/goose
RUN git clone https://github.com/vishnubob/wait-for-it
RUN go get -u github.com/tockins/realize
RUN go get -u github.com/tsenart/vegeta

