package rpc

import (
	"os"
	"testing"

	"github.com/thomasylee/GoRaft/global"
	"github.com/thomasylee/GoRaft/state"
)

const port string = "8000"

func TestMain(m *testing.M) {
	global.SetUpLogger()
	global.SetLogLevel("debug")

	global.TimeoutChannel = make(chan bool, 1)

	go RunServer(port)
	global.Log.Debug("gRPC server for tests has been started.")

	os.Exit(m.Run())
}

func resetTestEnvironment() {
	state.Node = state.NewNodeState(
		state.NewMemoryStateMachine(),
		state.NewMemoryStateMachine())
}

func Test_AppendEntries_WhenRequestHasNoEntries_ReturnsSuccess(t *testing.T) {
	resetTestEnvironment()

	request := &AppendEntriesRequest{
		Term: 0,
		LeaderId: "123",
		PrevLogIndex: 0,
		Entries: []*AppendEntriesRequest_Entry{},
		LeaderCommit: 0,
	}

	response, err := SendAppendEntries("127.0.0.1:" + port, request)
	if err != nil {
		t.Fatal(err)
	}

	if !response.Success {
		t.Error("Success was false")
	}
}
