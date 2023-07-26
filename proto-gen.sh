#!/bin/bash

set -e

# https://github.com/protocolbuffers/protobuf/releases/tag/v23.4
which protoc

# go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.31.0
which protoc-gen-go

find -name '*.proto'

protoc --go_out=. $(find -name '*.proto')
