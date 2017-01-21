package api

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
	"github.com/op/go-logging"

	"github.com/thomasylee/GoRaft/state"
)

/**
 * The leveled Logger to use in the api package.
 */
var Log *logging.Logger

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
	CommitIndex int
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
func handleAppendEntries(writer http.ResponseWriter, request *http.Request, timeoutChannel chan<- bool) {
	// Indicate that a message has been received so we don't time out.
	timeoutChannel <- true

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

	nodeState := state.GetNodeState(Log)

	currentTerm := nodeState.CurrentTerm()
	success := false

	// Don't append entries for a stale leader.
	if appendEntriesRequest.Term < currentTerm {
		writeAppendEntriesResponse(writer, currentTerm, success)
		return
	} else if appendEntriesRequest.Term > currentTerm {
		// Update currentTerm if the supplied term is newer.
		nodeState.SetCurrentTerm(appendEntriesRequest.Term)
		currentTerm = appendEntriesRequest.Term
	}

	// Make sure prevLogTerm matches the term of the entry at prevLogIndex, unless
	// this is the first entry to be applied.
	if (appendEntriesRequest.PrevLogIndex != -1 || nodeState.LogLength() != 0) && (
		nodeState.LogLength() <= appendEntriesRequest.PrevLogIndex ||
		appendEntriesRequest.PrevLogTerm != nodeState.Log(appendEntriesRequest.PrevLogIndex).Term) {
		writeAppendEntriesResponse(writer, currentTerm, success)
		return
	}

	requestEntriesIndex := 0
	success = true
	prevLogIndex := appendEntriesRequest.PrevLogIndex

	// Save all the log entries that were received, but trust that ones with the
	// same term don't need to be updated.
	for i := prevLogIndex + 1; i <= prevLogIndex + len(appendEntriesRequest.Entries); i++ {
		if nodeState.LogLength() <= i || nodeState.Log(i).Term != currentTerm {
			entry := appendEntriesRequest.Entries[i]

			logEntry := state.LogEntry{Key: entry.Key, Value: entry.Value, Term: currentTerm}
			err = nodeState.AppendEntryToLog(i, logEntry)
			if err != nil {
				Log.Error(err.Error())
				success = false
			}
		}
		requestEntriesIndex += 1
	}

	// Remove all log entries that existed in this node's state past the last
	// index given by the request.
	for i := prevLogIndex + len(appendEntriesRequest.Entries) + 1; ; i++ {
		key := strconv.Itoa(prevLogIndex)
		value, err := nodeState.NodeStateMachine.Get(key)
		if err != nil {
			Log.Error(err.Error())
			break
		}
		if value == "" {
			break
		}
		nodeState.NodeStateMachine.Put(key, value)
	}

	writeAppendEntriesResponse(writer, currentTerm, success)
}

/**
 * Converts the currentTerm and success to JSON and returns them with a 200 status code.
 */
func writeAppendEntriesResponse(writer http.ResponseWriter, currentTerm int, success bool) {
	appendEntriesResponse := AppendEntriesResponse{Term: currentTerm, Success: success}
	json.NewEncoder(writer).Encode(appendEntriesResponse)
}

/**
 * Runs the API server on the port specified in config.yaml.
 */
func RunServer(logger *logging.Logger, timeoutChannel chan<- bool, port int) {
	Log = logger
	router := mux.NewRouter()

	router.HandleFunc("/append_entries", func(writer http.ResponseWriter, request *http.Request) {
		handleAppendEntries(writer, request, timeoutChannel)
	})

	server := &http.Server{
		Handler: router,
		Addr: "127.0.0.1:" + strconv.Itoa(port),
		WriteTimeout: 10 * time.Second,
		ReadTimeout: 10 * time.Second,
	}

	Log.Notice(server.ListenAndServe())
}
