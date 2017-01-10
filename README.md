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
* [github.com/boltdb/bolt](https://github.com/boltdb/bolt) - To act as the persistent finite state machine (FSM)
* [github.com/gorilla/mux](https://github.com/gorilla/mux) - For simple REST API construction, though I may switch from mux to gRPC later to stay true to the Raft article

## Getting Started

### Installing
To install GoRaft, first install Go and then run "go get":
```sh
$ go get github.com/thomasylee/GoRaft
```

### Configuring
Configurations can be found in [config.yaml](https://github.com/thomasylee/GoRaft/blob/master/config.yaml). The details of each configuration is explained in the comments above the relevant key-value pairs.

### Running
Run using "go run" from the source directory, or run "go run" on the main.go file itself:
```sh
$ go run main.go
```
