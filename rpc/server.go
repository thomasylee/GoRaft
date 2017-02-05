package rpc

import (
	"net"
	"strconv"

	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	"github.com/thomasylee/GoRaft/global"
	"github.com/thomasylee/GoRaft/state"
)

// server is used to implement the GoRaft gRPC server.
type server struct{}

func (s *server) AppendEntries(ctx context.Context, request *AppendEntriesRequest) (*AppendEntriesResponse, error) {
	// Indicate that a message has been received so we don't time out.
	global.TimeoutChannel <- true

	// If the entries attribute is empty, it's just a heartbeat.
	if len(request.Entries) == 0 {
		return &AppendEntriesResponse{Term: state.Node.CurrentTerm(), Success: true}, nil
	}

	nodeState := state.GetNodeState()
	response := &AppendEntriesResponse{
		Term: nodeState.CurrentTerm(),
		Success: false,
	}
	prevLogIndex := request.PrevLogIndex
	logLength := nodeState.LogLength()

	// Don't append entries for a stale leader.
	if request.Term < response.Term {
		global.Log.Debug("success = false due to term being too old:", request.Term)
		return response, nil
	} else if request.Term > response.Term {
		// Update CurrentTerm if the supplied term is newer.
		response.Term = request.Term
		nodeState.SetCurrentTerm(response.Term)
	}

	global.Log.Debug("LogLength =", logLength)
	if prevLogIndex > logLength {
		global.Log.Debug("success = false due to PrevLogIndex > log length:", request.PrevLogIndex, nodeState.LogLength())
		return response, nil
	}

	// Make sure PrevLogTerm matches the term of the entry at PrevLogIndex, unless
	// the request considers these the first entries added to the log.
	global.Log.Debug("PrevLogIndex =", prevLogIndex)
	if prevLogIndex > 0 && request.PrevLogTerm != nodeState.Log(prevLogIndex).Term {

		global.Log.Debug("success = false due to PrevLogIndex mismatch:", prevLogIndex)
		return response, nil
	}

	requestEntriesIndex := 0
	response.Success = true
	var err error

	// Save all the log entries that were received, but trust that ones with the
	// same term don't need to be updated.
	for i := prevLogIndex + 1; i <= prevLogIndex + uint32(len(request.Entries)); i++ {
		if nodeState.LogLength() < i || nodeState.Log(i).Term != response.Term {
			entry := request.Entries[requestEntriesIndex]
			logEntry := state.LogEntry{
				Key: entry.Key,
				Value: entry.Value,
				Term: response.Term,
			}

			err = nodeState.SetLogEntry(i, logEntry)
			if err != nil {
				global.Log.Error(err.Error())
				response.Success = false
			}
		}
		requestEntriesIndex += 1
	}

	if !response.Success {
		return response, err
	}

	// Remove all log entries that existed in this node's state past the last
	// index given by the request.
	for i := prevLogIndex + uint32(len(request.Entries)) + 1; ; i++ {
		key := strconv.Itoa(int(i))
		value, err := nodeState.NodeStateMachine.Get(key)
		if err != nil {
			global.Log.Error(err.Error())
			break
		}
		if value == "" {
			break
		}
		nodeState.NodeStateMachine.Put(key, "")
	}

	// Update the node's commitIndex when the request's LeaderCommit index is
	// higher.
	if request.LeaderCommit > nodeState.CommitIndex() {
		if request.LeaderCommit > nodeState.LogLength() {
			nodeState.SetCommitIndex(nodeState.LogLength() - 1)
		} else {
			nodeState.SetCommitIndex(request.LeaderCommit)
		}
	}

	global.Log.Debug("success = true")
	return response, err
}

func RunServer(port string) {
	lis, err := net.Listen("tcp", ":" + port)
	if err != nil {
		global.Log.Panicf("Failed to listen: %v", err)
	}

	s := grpc.NewServer()
	RegisterGoRaftServer(s, &server{})

	// Register reflection service on gRPC server.
	reflection.Register(s)

	err = s.Serve(lis)
	if err != nil {
		global.Log.Panicf("Failed to serve: %v", err)
	}
}
