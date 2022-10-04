FROM ubuntu:22.04
RUN apt-get update
RUN apt-get dist-upgrade -y
RUN apt-get install ruby rubygems rpm curl -y
RUN gem install fpm
RUN curl -s -L https://go.dev/dl/go1.19.1.linux-amd64.tar.gz | tar -C /usr/local -xz
RUN mkdir /code
ENV PATH=$PATH:/usr/local/go/bin:/root/go/bin
ENV VERSION=4
WORKDIR /code