FROM       ubuntu:16.04

RUN apt-get update
RUN apt-get -y upgrade
RUN apt-get -y install golang curl

EXPOSE 7777
ENV GOPATH=/usr/local/challenge
RUN mkdir -p $GOPATH/src/github.com/dacamp/challenge
ADD . /usr/local/challenge/src/github.com/dacamp/challenge

WORKDIR $GOPATH/src/github.com/dacamp/challenge
RUN mkdir bin
RUN go build -o bin/counter
CMD bin/counter