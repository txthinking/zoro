#!/bin/bash
protoc -I=. -I=$GOPATH/src --gofast_out=. mr2.proto
