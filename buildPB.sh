#!/bin/bash
git clone https://github.com/gogo/protobuf.git /tmp/protobuf

go install github.com/gogo/protobuf/protoc-gen-gofast@latest

protoc -I=. -I=/tmp/protobuf/gogoproto -I=/tmp/protobuf/protobuf --gofast_out=plugins=grpc:. zoro.proto

rm -rf /tmp/protobuf
