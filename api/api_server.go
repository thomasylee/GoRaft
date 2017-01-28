package api

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"

	"github.com/thomasylee/GoRaft/global"
	"github.com/thomasylee/GoRaft/state"
)

/**
 * An Entry holds a log entry in an AppendEntriesRequest. It is different from
 * the LogEntry type in github.com/thomasylee/GoRaft/state since it does not
 * contain a Term.
 */
type Entry struct {
	Key string
	Value string
}

/**
 * An AppendEntriesRequest contains all the parameters that should be present in
 * append_entries calls that aren't heartbeats.
 */
type AppendEntriesRequest struct {
	Term int
	LeaderId string
	PrevLogIndex int
	PrevLogTerm int
	Entries []Entry
	LeaderCommit int
}

/**
 * An AppendEntriesResponse contains the expected return values from an append_entries call.
 */
type AppendEntriesResponse struct {
	Term int
	Success bool
}

/**
 * Handles append_entries calls, including heartbeats.
 */
func handleAppendEntries(writer http.ResponseWriter, request *http.Request) {
	// Indicate that a message has been received so we don't time out.
	global.TimeoutChannel <- true

	// If the request is empty, it's just a heartbeat.
	body, _ := ioutil.ReadAll(request.Body)
	if len(body) == 0 {
		return
	}

	// Unmarshal the JSON request body into an AppendEntriesRequest.
	var appendEntriesRequest AppendEntriesRequest
	err := json.Unmarshal(body, &appendEntriesRequest)
	if err != nil {
		http.Error(writer, err.Error(), http.StatusBadRequest)
		return
	}

	currentTerm, success := processAppendEntries(appendEntriesRequest, state.GetNodeState())

	appendEntriesResponse := AppendEntriesResponse{Term: currentTerm, Success: success}
	json.NewEncoder(writer).Encode(appendEntriesResponse)
}

/**
 * Processes the body of a non-heartbeat append_entries call.
 * The return values are the currentTerm and the success boolean.
 */
func processAppendEntries(request AppendEntriesRequest, nodeState state.NodeState) (int, bool) {
	currentTerm := nodeState.CurrentTerm()
	success := false

	// Don't append entries for a stale leader.
	if request.Term < currentTerm {
		global.Log.Debug("success = false due to term being too old:", currentTerm)
		return currentTerm, success
	} else if request.Term > currentTerm {
		// Update currentTerm if the supplied term is newer.
		nodeState.SetCurrentTerm(request.Term)
		currentTerm = request.Term
	}

	global.Log.Debug("LogLength =", nodeState.LogLength())
	if nodeState.LogLength() != 0 && request.PrevLogIndex == -1 {
		global.Log.Debug("success = false due to PrevLogIndex = -1:", request.PrevLogIndex)
		return currentTerm, success
	}

	// Make sure prevLogTerm matches the term of the entry at prevLogIndex, unless
	// this is the first entry to be applied.
	global.Log.Debug("PrevLogIndex =", request.PrevLogIndex)
	if (request.PrevLogIndex != -1 || nodeState.LogLength() != 0) && (
		nodeState.LogLength() <= request.PrevLogIndex ||
		request.PrevLogTerm != nodeState.Log(request.PrevLogIndex).Term) {

		global.Log.Debug("success = false due to PrevLogIndex mismatch:", request.PrevLogIndex)
		return currentTerm, success
	}

	requestEntriesIndex := 0
	success = true
	prevLogIndex := request.PrevLogIndex

	// Save all the log entries that were received, but trust that ones with the
	// same term don't need to be updated.
	for i := prevLogIndex + 1; i <= prevLogIndex + len(request.Entries); i++ {
		if nodeState.LogLength() <= i || nodeState.Log(i).Term != currentTerm {
			entry := request.Entries[i]

			logEntry := state.LogEntry{Key: entry.Key, Value: entry.Value, Term: currentTerm}
			err := nodeState.AppendEntryToLog(i, logEntry)
			if err != nil {
				global.Log.Error(err.Error())
				success = false
			}
		}
		requestEntriesIndex += 1
	}

	if !success {
		return currentTerm, success
	}

	// Remove all log entries that existed in this node's state past the last
	// index given by the request.
	for i := prevLogIndex + len(request.Entries) + 1; ; i++ {
		key := strconv.Itoa(prevLogIndex)
		value, err := nodeState.NodeStateMachine().Get(key)
		if err != nil {
			global.Log.Error(err.Error())
			break
		}
		if value == "" {
			break
		}
		nodeState.NodeStateMachine().Put(key, value)
	}

	global.Log.Debug("success = true")
	return currentTerm, success
}

/**
 * Runs the API server on the port specified in config.yaml.
 */
func RunServer(port int) {
	router := mux.NewRouter()

	router.HandleFunc("/append_entries", handleAppendEntries)
/*
	router.HandleFunc("/append_entries", func(writer http.ResponseWriter, request *http.Request) {
		handleAppendEntries(writer, request, timeoutChannel)
	})
*/

	server := &http.Server{
		Handler: router,
		Addr: "127.0.0.1:" + strconv.Itoa(port),
		WriteTimeout: 10 * time.Second,
		ReadTimeout: 10 * time.Second,
	}

	global.Log.Warning(server.ListenAndServe())
}
