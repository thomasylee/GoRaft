package state

import (
	"strconv"
	"testing"

	"github.com/thomasylee/GoRaft/global"
)

func createNodeState() *NodeState {
	return NewNodeState(NewMemoryStateMachine(), NewMemoryStateMachine())
}

func Test_AppendEntryToLog_WithValidNodeAndParams_AppendsEntryToMemAndStateMachine(t *testing.T) {
	global.SetUpLogger()
	global.SetLogLevel("critical")

	node := createNodeState()
	node.log = &[]LogEntry{}

	var tests = []struct {
		index int
		entry LogEntry
		jsonRep string
	}{
		{0, LogEntry{"a", "A", 0}, "{\"Key\":\"a\",\"Value\":\"A\",\"Term\":0}"},
		{1, LogEntry{"b", "B", 0}, "{\"Key\":\"b\",\"Value\":\"B\",\"Term\":0}"},
		{2, LogEntry{"c", "C", 1}, "{\"Key\":\"c\",\"Value\":\"C\",\"Term\":1}"},
		{3, LogEntry{"d", "D", 2}, "{\"Key\":\"d\",\"Value\":\"D\",\"Term\":2}"},
		{4, LogEntry{"e", "E", 2}, "{\"Key\":\"e\",\"Value\":\"E\",\"Term\":2}"},
	}

	for _, test := range tests {
		node.AppendEntryToLog(test.index, test.entry)

		entryInMem := node.Log(test.index)
		if entryInMem != test.entry {
			t.Error("Log entry in memory doesn't match:", test.index, entryInMem)
		}

		jsonInSM, err := node.NodeStateMachine.Get(strconv.Itoa(test.index))
		if err != nil {
			t.Errorf("Error processing entry %d: %s", test.index, err.Error())
		} else if jsonInSM != test.jsonRep {
			t.Error("Log entry JSON in state machine doesn't match:", test.index, jsonInSM)
		}
	}
}
