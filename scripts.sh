#!/bin/bash

function build(){
  go build -o tok
}

function test(){
  ginkgo test ./...
}

# Run a function name in the context of this script
eval "$@"