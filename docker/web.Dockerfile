FROM golang:1.18-alpine

# Required because go requires gcc to build
RUN apk add build-base

RUN apk add inotify-tools

RUN echo $GOPATH

RUN go install github.com/go-delve/delve/cmd/dlv@latest

RUN go install github.com/rubenv/sql-migrate/...@latest

COPY . /clean_web

WORKDIR /clean_web

ENV GOFLAGS -buildvcs=false

RUN go mod tidy

CMD sh /clean_web/docker/run.sh