FROM golang:latest

MAINTAINER Mahesh Vangala <vangalamaheshh@gmail.com>

WORKDIR /usr/local/

RUN set -ex \
  && git clone https://github.com/dst-umms/Synergist-BQ-API.git

RUN set -ex \
  && go get -u cloud.google.com/go/bigquery \
  && go get -u golang.org/x/net/context

CMD ["go", "run", "Synergist-BQ-API/main.go"]
