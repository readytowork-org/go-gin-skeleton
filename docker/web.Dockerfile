FROM golang:1.16-alpine

# Required because go requires gcc to build
RUN apk add build-base

RUN apk add inotify-tools

RUN echo $GOPATH

COPY . /clean_web

RUN go get github.com/rubenv/sql-migrate/...

WORKDIR /clean_web

RUN go mod tidy

RUN go get github.com/go-delve/delve/cmd/dlv

CMD sh /clean_web/docker/run.sh