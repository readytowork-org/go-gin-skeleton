#!/bin/bash

ginWatch() {
  gin -a $1 -i run .
}

if which gin; then
  ginWatch $1
else
  go install github.com/codegangsta/gin@latest
  ginWatch $1
fi
