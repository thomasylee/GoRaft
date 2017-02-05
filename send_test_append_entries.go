package main

import (
	"math/rand"
	"strconv"

	"github.com/thomasylee/GoRaft/global"
	"github.com/thomasylee/GoRaft/rpc"
)

// The previous log index for the AppendEntries rpc call.
const prevLogIndex uint32 = 0

// This program can be modified and used to test AppendEntries requests.
func main() {
	global.SetUpLogger()

	value := rand.Intn(100)
	global.Log.Info("Generated value:", value)

	request := &rpc.AppendEntriesRequest{
		Term: 0,
		LeaderId: "123",
		PrevLogIndex: prevLogIndex,
		PrevLogTerm: 0,
		Entries: []*rpc.AppendEntriesRequest_Entry{{Term: 0, Key: "a", Value: strconv.Itoa(value)}},
		LeaderCommit: 0,
	}

	response, err := rpc.SendAppendEntries("127.0.0.1:8000", request)
	if err != nil {
		panic(err)
	}
	global.Log.Info("Success!")
	global.Log.Info(response)
}
