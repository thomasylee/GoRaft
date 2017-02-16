package state

import (
	"encoding/json"
	"strconv"

	"github.com/thomasylee/GoRaft/global"
)

// A LogEntry represents a state machine command (Put {Key},{Value}) and the
// term when the entry was received by the leader.
type LogEntry struct {
	Key string
	Value string
	Term uint32
}

// Define constants for important keys in the Bolt database.
const (
	currentTerm string = "CurrentTerm"
	votedFor string = "VotedFor"
	logEntries string = "LogEntries"
)

// A NodeState contains both the persistent and volatile state that a node
// needs to function in a Raft cluster.
//
// Note that currentTerm, votedFor, and log are not exported because they need
// to be persistent, so to preserve consistency, they should be updated and
// retrieved using the respective NodeState methods.
type NodeState struct {
	// The latest term the node has seen.
	currentTerm uint32

	// The candidateId (UUID) that was voted for the current term.
	votedFor string

	// Log entries, each containing a command for the state machine and the term
	// when the entry was received from the leader.
	log *[]LogEntry

	// The node id of the current leader.
	LeaderId string

	// Node state and stored state are stored by different state machines.
	NodeStateMachine StateMachine
	StorageStateMachine StateMachine

	// Index of the highest log entry known to be committed.
	CommitIndex uint32

	// Index of the highest log entry applied to this node's storage state machine.
	LastApplied uint32

	// (Leader only) For each node, the index of the next log entry to send to.
	NextIndex map[string]uint32

	// (Leader only) For each node, the index of the highest log entry known to be replicated on
	// the node.
	MatchIndex map[string]uint32
}

var Node *NodeState

// Returns the global NodeState variable, initializing it if it hasn't yet been
// instantiated.
func GetNodeState() *NodeState {
	if Node != nil {
		return Node
	}

	nodeStateMachine, err := NewBoltStateMachine("node_state.db")
	if err != nil {
		global.Log.Panic("Failed to initialize nodeStateMachine:", err.Error())
	}
	storageStateMachine, err := NewBoltStateMachine("state.db")
	if err != nil {
		global.Log.Panic("Failed to initialize storageStateMachine:", err.Error())
	}
	Node = NewNodeState(nodeStateMachine, storageStateMachine)
	return Node
}

// Returns a NodeState based on values in the node state Bolt database, using
// default values if the database does not exist or have any values in it.
func NewNodeState(nodeStateMachine StateMachine, storageStateMachine StateMachine) *NodeState {
	var currentTermValue uint32
	retrievedCurrentTerm, err := nodeStateMachine.Get(currentTerm)
	if err != nil {
		global.Log.Panic("Failed to retrieve CurrentTerm:", err.Error())
	}
	if retrievedCurrentTerm == "" {
		currentTermValue = 0
	} else {
		currentTermValueInt, err := strconv.Atoi(retrievedCurrentTerm)
		if err != nil {
			global.Log.Panic("Failed to convert CurrentTerm to int:", err.Error())
		}
		currentTermValue = uint32(currentTermValueInt)
	}

	votedForValue, err := nodeStateMachine.Get(votedFor)
	if err != nil {
		global.Log.Panic("Failed to retrieve VotedFor:", err.Error())
	}

	logEntries, err := nodeStateMachine.RetrieveLogEntries(1, 1000000)
	if err != nil {
		global.Log.Panic("Failed to retrieve log entries:", err.Error())
	}
	global.Log.Debug("Pre-existing log entries:", len(logEntries))

	var node *NodeState
	node = &NodeState{
		NodeStateMachine: nodeStateMachine,
		StorageStateMachine: storageStateMachine,
	}
	node.SetCurrentTerm(currentTermValue)
	node.SetVotedFor(votedForValue)
	node.log = &logEntries

	return node
}

// Sets the current term in memory and in the node state machine.
func (state *NodeState) SetCurrentTerm(newCurrentTerm uint32) {
	state.currentTerm = newCurrentTerm
	global.Log.Debugf("CurrentTerm updated: %d", newCurrentTerm)
	state.NodeStateMachine.Put(currentTerm, strconv.Itoa(int(newCurrentTerm)))
}

// Returns the current term recognized by the node.
func (state *NodeState) CurrentTerm() uint32 {
	return state.currentTerm;
}

// Sets VotedFor in memory and in the node state machine.
func (state *NodeState) SetVotedFor(newVotedFor string) {
	state.votedFor = newVotedFor
	state.NodeStateMachine.Put(votedFor, newVotedFor)
}

// Returns the node's VotedFor.
func (state *NodeState) VotedFor() string {
	return state.votedFor;
}

// Sets the log entry in the NodeState's log at the given index.
// Note that this method does not do any safety checking to prevent overwriting
// existing entries; that check should be done by the caller beforehand.
func (state *NodeState) SetLogEntry(index uint32, entry LogEntry) error {
	jsonValue, err := json.Marshal(entry)
	if err != nil {
		return err
	}

	err = state.NodeStateMachine.Put(strconv.Itoa(int(index)), string(jsonValue))
	if err != nil {
		return err
	}

	for i := uint32(len(*state.log)); i < index - 1; i++ {
		*state.log = append(*state.log, LogEntry{})
	}
	*state.log = append(*state.log, entry)
	return nil
}

// Returns the number of entries in the node's log.
func (state NodeState) LogLength() uint32 {
	return uint32(len(*state.log))
}

// Returns the LogEntry at the specified index. Note that log indices start
// at 1, but the slice indices start at 0.
func (state NodeState) Log(index uint32) LogEntry {
	return (*state.log)[index - 1]
}
