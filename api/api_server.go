package api

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
)

type AppendEntriesRequest struct {
	Term int
	LeaderId string
	PrevLogIndex int
	PrevLogTerm int
	Entries []string
	CommitIndex int
}

type AppendEntriesResponse struct {
	Term int
	Success bool
}

func handleAppendEntries(writer http.ResponseWriter, request *http.Request, timeoutChannel chan<- bool) {
	// Indicate that a message has been received so we don't time out.
	timeoutChannel <- true

	// If the request is empty, it's just a heartbeat.
	body, _ := ioutil.ReadAll(request.Body)
	if len(body) == 0 {
		return
	}

	var appendEntriesRequest AppendEntriesRequest
	json.NewDecoder(request.Body).Decode(appendEntriesRequest)

	// TODO: Add the entries to the log and return the correct Term
	// and Success values.

	appendEntriesResponse := AppendEntriesResponse{Term: 0, Success: true}
	json.NewEncoder(writer).Encode(appendEntriesResponse)
}

func RunServer(timeoutChannel chan<- bool, port int) {
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

	log.Fatal(server.ListenAndServe())
}
