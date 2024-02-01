#!/bin/sh

PATH=$PATH:$GOPATH/bin
genpath=$(pwd)/genproto/.

protoc \
--proto_path=./protos/posts \
--go_out=$genpath \
--go_opt=Mcloudbees/protos/posts/posts.proto=cloudbees/genproto/posts \
--go-grpc_out=$genpath \
--go-grpc_opt=Mcloudbees/protos/posts/posts.proto=cloudbees/genproto/posts \
 ./protos/posts/posts.proto
