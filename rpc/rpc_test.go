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
	close(global.TimeoutChannel)
	global.TimeoutChannel = make(chan bool, 1)

	state.Node = state.NewNodeState(
		state.NewMemoryStateMachine(),
		state.NewMemoryStateMachine())
}

func Test_AppendEntries_WhenRequestHasNoEntries_ReturnsSuccessTrue(t *testing.T) {
	resetTestEnvironment()

	request := &AppendEntriesRequest{
		Term: 0,
		LeaderId: "123",
		PrevLogIndex: 0,
		PrevLogTerm: 0,
		// Heartbeat request since entries is empty.
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

func Test_AppendEntries_WhenTermIsOlderThanCurrentTerm_ReturnsSuccessFalse(t *testing.T) {
	resetTestEnvironment()

	state.Node.SetCurrentTerm(1)

	request := &AppendEntriesRequest{
		// Request should fail since CurrentTerm is 1, which is > 0.
		Term: 0,
		LeaderId: "123",
		PrevLogIndex: 0,
		PrevLogTerm: 0,
		Entries: []*AppendEntriesRequest_Entry{
			&AppendEntriesRequest_Entry{
				Term: 0,
				Key: "a",
				Value: "A",
			},
		},
		LeaderCommit: 0,
	}

	response, err := SendAppendEntries("127.0.0.1:" + port, request)
	if err != nil {
		t.Fatal(err)
	}

	if response.Success {
		t.Error("Success was true")
	}

	if response.Term != 1 {
		t.Error("Term was not 1:", response.Term)
	}
}

func Test_AppendEntries_WhenLogAtPrevLogIndexDoesNotMatch_ReturnsSuccessFalse(t *testing.T) {
	resetTestEnvironment()

	state.Node.SetCurrentTerm(1)

	entry := state.LogEntry{Term: 1, Key: "a", Value: "A"}
	state.Node.SetLogEntry(1, entry)

	request := &AppendEntriesRequest{
		Term: 1,
		LeaderId: "123",
		PrevLogIndex: 1,
		// Request should fail since term for entry 1 is 1, not 0.
		PrevLogTerm: 0,
		Entries: []*AppendEntriesRequest_Entry{
			&AppendEntriesRequest_Entry{
				Term: 1,
				Key: "b",
				Value: "B",
			},
		},
		LeaderCommit: 0,
	}

	response, err := SendAppendEntries("127.0.0.1:" + port, request)
	if err != nil {
		t.Fatal(err)
	}

	if response.Success {
		t.Error("Success was true")
	}

	if response.Term != 1 {
		t.Error("Term was not 1:", response.Term)
	}
}

func Test_AppendEntries_WhenLeaderCommitIsGreaterThanCommitIndex_IncreasesCommitIndex(t *testing.T) {
	resetTestEnvironment()

	state.Node.SetCurrentTerm(0)

	entry := state.LogEntry{Term: 0, Key: "a", Value: "A"}
	state.Node.SetLogEntry(1, entry)

	request := &AppendEntriesRequest{
		Term: 1,
		LeaderId: "123",
		PrevLogIndex: 1,
		// Request should fail since term for entry 1 is 1, not 0.
		PrevLogTerm: 0,
		Entries: []*AppendEntriesRequest_Entry{
			&AppendEntriesRequest_Entry{
				Term: 1,
				Key: "b",
				Value: "B",
			},
		},
		LeaderCommit: 2,
	}

	response, err := SendAppendEntries("127.0.0.1:" + port, request)
	if err != nil {
		t.Fatal(err)
	}

	if !response.Success {
		t.Error("Success was true")
	}

	if response.Term != 1 {
		t.Error("Term was not 1:", response.Term)
	}

	if state.Node.CommitIndex() != 2 {
		t.Error("CommitIndex was not 2:", state.Node.CommitIndex())
	}
}
