cd rpc
protoc -I . goraft.proto --go_out=plugins=grpc:proto
