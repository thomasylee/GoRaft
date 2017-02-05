package rpc

import (
	"golang.org/x/net/context"
	"google.golang.org/grpc"

	pb "github.com/thomasylee/GoRaft/rpc/proto"
)

func SendAppendEntries(address string, request *pb.AppendEntriesRequest) (*pb.AppendEntriesResponse, error) {
	conn, err := grpc.Dial(address, grpc.WithInsecure())
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	client := pb.NewGoRaftClient(conn)

	response, err := client.AppendEntries(context.Background(), request)
	if err != nil {
		return nil, err
	}

	return response, nil
}
