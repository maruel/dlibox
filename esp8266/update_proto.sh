#!/bin/bash

set -eu

cd "$(dirname $0)"

protoc --plugin=../proto/nanopb/generator/protoc-gen-nanopb \
  --nanopb_out=--no-timestamp:src --proto_path=src \
  -I ../proto \
  -I ../proto/nanopb/generator/proto \
  ../proto/anim1d_msg.proto
