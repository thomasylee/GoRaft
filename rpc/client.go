package rpc

import (
	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

// SendAppendEntries sends an AppendEntries request to the specified address.
func SendAppendEntries(address string, request *AppendEntriesRequest) (*AppendEntriesResponse, error) {
	conn, err := grpc.Dial(address, grpc.WithInsecure())
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	client := NewGoRaftClient(conn)

	response, err := client.AppendEntries(context.Background(), request)
	if err != nil {
		return nil, err
	}

	return response, nil
}

// SendRequestVote sends a RequestVote request to the specified address.
func SendRequestVote(address string, request *RequestVoteRequest) (*RequestVoteResponse, error) {
	conn, err := grpc.Dial(address, grpc.WithInsecure())
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	client := NewGoRaftClient(conn)

	response, err := client.RequestVote(context.Background(), request)
	if err != nil {
		return nil, err
	}

	return response, nil
}
