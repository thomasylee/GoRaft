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
		index uint32
		entry LogEntry
		jsonRep string
	}{
		{1, LogEntry{"a", "A", 0}, "{\"Key\":\"a\",\"Value\":\"A\",\"Term\":0}"},
		{2, LogEntry{"b", "B", 0}, "{\"Key\":\"b\",\"Value\":\"B\",\"Term\":0}"},
		{3, LogEntry{"c", "C", 1}, "{\"Key\":\"c\",\"Value\":\"C\",\"Term\":1}"},
		{4, LogEntry{"d", "D", 2}, "{\"Key\":\"d\",\"Value\":\"D\",\"Term\":2}"},
		{5, LogEntry{"e", "E", 2}, "{\"Key\":\"e\",\"Value\":\"E\",\"Term\":2}"},
	}

	for _, test := range tests {
		node.SetLogEntry(test.index, test.entry)

		// Log entry indices start at 1, but the entryInMem slice indices start at 0,
		// so we must add 1 to make the indices align.
		entryInMem := node.Log(test.index + 1)
		if entryInMem != test.entry {
			t.Error("Log entry in memory doesn't match:", test.index, entryInMem)
		}

		jsonInSM, err := node.NodeStateMachine.Get(strconv.Itoa(int(test.index)))
		if err != nil {
			t.Errorf("Error processing entry %d: %s", test.index, err.Error())
		} else if jsonInSM != test.jsonRep {
			t.Error("Log entry JSON in state machine doesn't match:", test.index, jsonInSM)
		}
	}
}
