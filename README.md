[![Build Status](https://travis-ci.org/thomasylee/GoRaft.svg?branch=master)](https://travis-ci.org/thomasylee/GoRaft)
[![Codacy Coverage Badge](https://api.codacy.com/project/badge/Coverage/17f8f1370cae4b05a2677a85213deb81)](https://www.codacy.com/app/thomasylee/GoRaft?utm_source=github.com&utm_medium=referral&utm_content=thomasylee/GoRaft&utm_campaign=Badge_Coverage)
[![Codacy Grade Badge](https://api.codacy.com/project/badge/Grade/17f8f1370cae4b05a2677a85213deb81)](https://www.codacy.com/app/thomasylee/GoRaft?utm_source=github.com&amp;utm_medium=referral&amp;utm_content=thomasylee/GoRaft&amp;utm_campaign=Badge_Grade)
[![Go Report Card](https://goreportcard.com/badge/github.com/thomasylee/GoRaft)](https://goreportcard.com/report/github.com/thomasylee/GoRaft)

# GoRaft
GoRaft is an implementation of the Raft consensus algorithm written in Go.

The goal of this project is for me to learn Go and become better acquainted with Raft's inner workings. The implementation mostly follows the Raft definition found in [In Search of an Understandable Consensus Algorithm](https://raft.github.io/raft.pdf) by Diego Ongaro and John Ousterhout, though some parts of the code may stray as I play with different ways to get Go to fulfill various functionalities of Raft.

## Table of Contents
- [Dependencies](#dependencies)
- [Getting Started](#getting-started)
 - [Installing](#installing)
 - [Configuring](#configuring)
 - [Running](#running)

## Dependencies
GoRaft depends on the following external packages:
* [github.com/boltdb/bolt](https://github.com/boltdb/bolt) - To store persistent state
* [github.com/go-yaml/yaml](https://github.com/go-yaml/yaml) - To parse the YAML config file
* [github.com/golang/protobuf](https://github.com/golang/protobuf) - To serialize gRPC requests/responses
* [github.com/op/go-logging](https://github.com/op/go-logging) - For leveled logging
* [golang.org/x/net/context](https://godoc.org/golang.org/x/net/context) - For context handling across gRPC requests
* [google.golang.org/grpc](https://godoc.org/google.golang.org/grpc) - To make/receive gRPC requests/responses

This package uses gRPC to send/receive rpc requests between nodes. \*.proto files will need to be compiled into Go code using protoc to be usable by gRPC. The executable file generate_protobuf.sh runs the necessary commands to compile all \*.proto files.

## Getting Started

### Installing
To install GoRaft, first install Go and then run "go get":
```sh
$ go get github.com/thomasylee/GoRaft
```

### Configuring
Configurations can be found in [config.yaml](https://github.com/thomasylee/GoRaft/blob/master/config.yaml). The details of each configuration is explained in the comments above the relevant key-value pairs.

### Testing

```sh
$ go test -cover -v ./...
```

### Running
Run using "go run" from the source directory, or run "go run" on the main.go file itself:
```sh
$ go run main.go
```

For now, the send_test_append_entries.go program can be used to append new entries to the node logs. It must be edited before being run to include the correct request values.
```sh
# Rename, since two files with main() methods will break the test setup.
mv send_test_append_entries.go2 send_test_append_entries.go
go run send_test_append_entries.go
```

The inspect_bolt.go file can be used to alter and inspect Bolt database files.
```sh
# Rename, since two files with main() methods will break the test setup.
mv run inspect_bolt.go2 inspect_bolt.go
go run inspect_bolt.go
```
