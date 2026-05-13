FROM golang:1.23.3-alpine

# add user group
# RUN addgroup -S nonroot \
#     && adduser -S nonroot -G nonroot

# USER nonroot

# Required because go requires gcc to build
RUN apk add build-base

RUN apk add inotify-tools

RUN echo $GOPATH

RUN go install github.com/go-delve/delve/cmd/dlv@latest

COPY . /clean_web

RUN apk add --no-cache curl \
    && curl -sSf https://atlasgo.sh | sh \
    && mv /atlas /usr/bin/atlas

WORKDIR /clean_web

RUN go install github.com/swaggo/swag/cmd/swag@latest && \
    export PATH=$PATH:$(go env GOPATH)/bin

ENV GOFLAGS -buildvcs=false

RUN go mod tidy

CMD sh /clean_web/docker/run.sh
