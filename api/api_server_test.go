package api

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/thomasylee/GoRaft/global"
	"github.com/thomasylee/GoRaft/state"
)

func Test_handleAppendEntries_WhenRequestIsEmpty_ReturnsEmpty200(t *testing.T) {
	resetTestEnvironment()

	// An empty append_entries request is just a heartbeat.
	request, err := http.NewRequest("POST", "http://localhost:8000/append_entries", strings.NewReader(""))
	if err != nil {
		t.Fatal(err)
	}

	writer := httptest.NewRecorder()
	handleAppendEntries(writer, request)

	// Make sure we wrote to the timeoutChannel within a second.
	select {
	case <-global.TimeoutChannel:
	case <-time.After(time.Second):
		t.Error("No message sent on timeoutChannel.")
	}

	if writer.Code != 200 {
		t.Error("Response from server was not 200:", writer.Code)
	}

	body := writer.Body.String()
	if err != nil {
		t.Fatalf("%s: %s", err.Error(), body)
	}
	if len(body) > 0 {
		t.Error("The heartbeat response was not empty:", body)
	}
}

func Test_handleAppendEntries_WhenRequestIsNotEmpty_ReturnsValidJson(t *testing.T) {
	resetTestEnvironment()

	jsonRequest := "{\"term\": 1, \"leaderId\": \"abc\", \"prevLogIndex\": -1, \"PrevLogTerm\": -1, \"entries\": [{\"key\": \"a\", \"value\": \"1\"}], \"leaderCommit\": 0}"

	request, err := http.NewRequest("POST", "http://localhost:8000/append_entries", strings.NewReader(jsonRequest))
	if err != nil {
		t.Fatal(err)
	}

	writer := httptest.NewRecorder()
	handleAppendEntries(writer, request)

	// Make sure we wrote to the timeoutChannel within a second.
	select {
	case <-global.TimeoutChannel:
	case <-time.After(time.Second):
		t.Error("No message sent on timeoutChannel.")
	}

	if writer.Code != 200 {
		t.Error("Response from server was not 200:", writer.Code)
	}

	body := writer.Body.String()
	if err != nil {
		t.Fatalf("%s: %s", err.Error(), body)
	}

	response := AppendEntriesResponse{}
	err = json.Unmarshal([]byte(body), &response)
	if err != nil {
		t.Fatalf("%s: %s", err.Error(), body)
	}
	if response.Success != true {
		t.Error("Response success value was false.")
	}
	if response.Term != 1 {
		t.Error("Response term was not 1:", response.Term)
	}
}

func Test_processAppendEntries_WhenRequestTermIsTooOld_ReturnsUnsuccessful(t *testing.T) {
	resetTestEnvironment()

	entry := Entry{Key: "a", Value: "1"}
	entries := []Entry{entry}
	request := AppendEntriesRequest{Term: -1, LeaderId: "b394b092-f840-406f-9284-ec3a6e0a2aa9", PrevLogIndex: -1, PrevLogTerm: -1, Entries: entries, LeaderCommit: -1}

	_, success := processAppendEntries(request)

	if success {
		t.Error("AppendEntriesResponse had success value true.")
	}
}

func Test_processAppendEntries_WhenRequestTermIsMoreRecent_UpdatesCurrentTerm(t *testing.T) {
	resetTestEnvironment()

	entry := Entry{Key: "a", Value: "1"}
	entries := []Entry{entry}
	request := AppendEntriesRequest{Term: 1, LeaderId: "b394b092-f840-406f-9284-ec3a6e0a2aa9", PrevLogIndex: -1, PrevLogTerm: -1, Entries: entries, LeaderCommit: -1}

	currentTerm, success := processAppendEntries(request)

	if !success {
		t.Error("AppendEntriesResponse had success value false.")
	}
	if currentTerm != 1 {
		t.Error("Returned CurrentTerm did not have value 1:", currentTerm)
	}

	retrievedTerm := state.GetNodeState().CurrentTerm()
	if retrievedTerm != 1 {
		t.Error("CurrentTerm saved to node state did not have value 1:", retrievedTerm)
	}
}
