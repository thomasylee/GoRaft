package rpc

import (
	"os"
	"testing"
	"time"

	"github.com/thomasylee/GoRaft/global"
	"github.com/thomasylee/GoRaft/state"
)

const port string = "8000"

func TestMain(m *testing.M) {
	global.SetUpLogger()
	global.SetLogLevel("debug")

	global.TimeoutChannel = make(chan bool, 1)

	go RunServer(port)
	// Give the server a few seconds to start.
	time.Sleep(3 * time.Second)
	global.Log.Debug("gRPC server for tests has been started.")

	os.Exit(m.Run())
}

func resetTestEnvironment() {
	close(global.TimeoutChannel)
	global.TimeoutChannel = make(chan bool, 1)

	state.Node = state.NewNodeState(
		state.NewMemoryDataStore(),
		state.NewMemoryDataStore())
}

func Test_AppendEntries_WhenRequestHasNoEntries_ReturnsSuccessTrue(t *testing.T) {
	resetTestEnvironment()

	request := &AppendEntriesRequest{
		Term:         0,
		LeaderId:     "123",
		PrevLogIndex: 0,
		PrevLogTerm:  0,
		// Heartbeat request since entries is empty.
		Entries:      []*AppendEntriesRequest_Entry{},
		LeaderCommit: 0,
	}

	response, err := SendAppendEntries("127.0.0.1:"+port, request)
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
		Term:         0,
		LeaderId:     "123",
		PrevLogIndex: 0,
		PrevLogTerm:  0,
		Entries: []*AppendEntriesRequest_Entry{
			&AppendEntriesRequest_Entry{
				Key:   "a",
				Value: "A",
			},
		},
		LeaderCommit: 0,
	}

	response, err := SendAppendEntries("127.0.0.1:"+port, request)
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
		Term:         1,
		LeaderId:     "123",
		PrevLogIndex: 1,
		// Request should fail since term for entry 1 is 1, not 0.
		PrevLogTerm: 0,
		Entries: []*AppendEntriesRequest_Entry{
			&AppendEntriesRequest_Entry{
				Key:   "b",
				Value: "B",
			},
		},
		LeaderCommit: 0,
	}

	response, err := SendAppendEntries("127.0.0.1:"+port, request)
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
		Term:         1,
		LeaderId:     "123",
		PrevLogIndex: 1,
		// Request should fail since term for entry 1 is 1, not 0.
		PrevLogTerm: 0,
		Entries: []*AppendEntriesRequest_Entry{
			&AppendEntriesRequest_Entry{
				Key:   "b",
				Value: "B",
			},
		},
		LeaderCommit: 2,
	}

	response, err := SendAppendEntries("127.0.0.1:"+port, request)
	if err != nil {
		t.Fatal(err)
	}

	if !response.Success {
		t.Error("Success was true")
	}

	if response.Term != 1 {
		t.Error("Term was not 1:", response.Term)
	}

	if state.Node.CommitIndex != 2 {
		t.Error("CommitIndex was not 2:", state.Node.CommitIndex)
	}
}

func Test_RequestVote(t *testing.T) {
	resetTestEnvironment()

	state.Node.SetVotedFor("1")

	request := &AppendEntriesRequest{
		Term:         1,
		LeaderId:     "123",
		PrevLogIndex: 0,
		PrevLogTerm:  0,
		Entries: []*AppendEntriesRequest_Entry{
			&AppendEntriesRequest_Entry{
				Key:   "a",
				Value: "A",
			},
		},
		LeaderCommit: 0,
	}

	_, err := SendAppendEntries("127.0.0.1:"+port, request)
	if err != nil {
		t.Fatal(err)
	}

	var tests = []struct {
		request  RequestVoteRequest
		response RequestVoteResponse
	}{
		{
			// Term is too old.
			RequestVoteRequest{
				Term:         0,
				CandidateId:  "1",
				LastLogIndex: 1,
				LastLogTerm:  0,
			},
			RequestVoteResponse{
				Term:        1,
				VoteGranted: false,
			},
		},
		{
			// Voted for a different candidate.
			RequestVoteRequest{
				Term:         1,
				CandidateId:  "2",
				LastLogIndex: 1,
				LastLogTerm:  1,
			},
			RequestVoteResponse{
				Term:        1,
				VoteGranted: false,
			},
		},
		{
			// The log is out of date.
			RequestVoteRequest{
				Term:         1,
				CandidateId:  "1",
				LastLogIndex: 0,
				LastLogTerm:  0,
			},
			RequestVoteResponse{
				Term:        1,
				VoteGranted: false,
			},
		},
		{
			// The log is in conflict with this node's log.
			RequestVoteRequest{
				Term:         1,
				CandidateId:  "1",
				LastLogIndex: 1,
				LastLogTerm:  2,
			},
			RequestVoteResponse{
				Term:        1,
				VoteGranted: false,
			},
		},
		{
			// Everything is good!
			RequestVoteRequest{
				Term:         1,
				CandidateId:  "1",
				LastLogIndex: 1,
				LastLogTerm:  1,
			},
			RequestVoteResponse{
				Term:        1,
				VoteGranted: true,
			},
		},
	}

	for _, test := range tests {
		response, err := SendRequestVote("127.0.0.1:"+port, &test.request)
		if err != nil {
			t.Error(err)
			continue
		}

		if *response != test.response {
			t.Errorf("RequestVote test case failed %v. Expected %v but found %v",
				test.request, test.response, *response)
		}
	}
}
