FROM       ubuntu:16.04

EXPOSE 7777
RUN apt-get update
RUN apt-get -y upgrade

COPY challenge.tar.gz /tmp

RUN mkdir -p /opt
RUN tar -xzvf /tmp/challenge.tar.gz -C /opt/

WORKDIR /opt/challenge
CMD bin/challenge-executable