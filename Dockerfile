FROM       ubuntu:16.04

EXPOSE 7777
RUN apt-get update
RUN apt-get -y upgrade
RUN apt-get -y install wget

RUN wget https://storage.googleapis.com/golang/go1.7.linux-amd64.tar.gz
RUN tar -xvf go1.7.linux-amd64.tar.gz
RUN mv go /usr/local

ENV GOROOT=/usr/local/go
ENV PATH=$GOROOT/bin:$PATH
ENV GOPATH=/usr/local/challenge

RUN mkdir -p $GOPATH/src/github.com/dacamp/challenge
ADD . /usr/local/challenge/src/github.com/dacamp/challenge

WORKDIR $GOPATH/src/github.com/dacamp/challenge
RUN mkdir bin
RUN go build -o bin/counter
CMD bin/counter