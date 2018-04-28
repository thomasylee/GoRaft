#!/usr/bin/env bash
protoc -I rpc rpc/goraft.proto --go_out=plugins=grpc:rpc
