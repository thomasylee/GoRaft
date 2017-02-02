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
	Term int
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
	currentTerm int

	// The candidateId (UUID) that was voted for the current term.
	votedFor string

	// Log entries, each containing a command for the state machine and the term
	// when the entry was received from the leader.
	log *[]LogEntry

	// Node state and stored state are stored by different state machines.
	NodeStateMachine StateMachine
	StorageStateMachine StateMachine

	// Index of the highest log entry known to be committed.
	commitIndex int

	// Index of the highest log entry applied to this node's storage state machine.
	lastApplied int

	// (Leader only) For each node, the index of the next log entry to send to.
	nextIndex map[int]int

	// (Leader only) For each node, the index of the highest log entry known to be replicated on
	// the node.
	matchIndex map[int]int
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
	var currentTermValue int
	retrievedCurrentTerm, err := nodeStateMachine.Get(currentTerm)
	if err != nil {
		global.Log.Panic("Failed to retrieve CurrentTerm:", err.Error())
	}
	if retrievedCurrentTerm == "" {
		currentTermValue = 0
	} else {
		currentTermValue, err = strconv.Atoi(retrievedCurrentTerm)
		if err != nil {
			global.Log.Panic("Failed to convert CurrentTerm to int:", err.Error())
		}
	}

	votedForValue, err := nodeStateMachine.Get(votedFor)
	if err != nil {
		global.Log.Panic("Failed to retrieve VotedFor:", err.Error())
	}

	logEntries, err := nodeStateMachine.RetrieveLogEntries(0, 1000000)
	if err != nil {
		global.Log.Panic("Failed to retrieve log entries:", err.Error())
	}

	var node *NodeState
	node = &NodeState{NodeStateMachine: nodeStateMachine, StorageStateMachine: storageStateMachine}
	node.SetCurrentTerm(currentTermValue)
	node.SetVotedFor(votedForValue)
	node.log = &logEntries

	return node
}

// Sets the current term in memory and in the node state machine.
func (state *NodeState) SetCurrentTerm(newCurrentTerm int) {
	state.currentTerm = newCurrentTerm
	global.Log.Debugf("CurrentTerm updated: %d", newCurrentTerm)
	state.NodeStateMachine.Put(currentTerm, strconv.Itoa(newCurrentTerm))
}

// Returns the current term recognized by the node.
func (state *NodeState) CurrentTerm() int {
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

// Append the log entry to the NodeState's log.
// Note that this method does not do any safety checking to prevent overwriting
// existing entries; that check should be done by the caller beforehand.
func (state *NodeState) AppendEntryToLog(index int, entry LogEntry) error {
	jsonValue, err := json.Marshal(entry)
	if err != nil {
		return err
	}

	err = state.NodeStateMachine.Put(strconv.Itoa(index), string(jsonValue))
	if err != nil {
		return err
	}

	*state.log = append(*state.log, entry)
	return nil
}

// Returns the number of entries in the node's log.
func (state *NodeState) LogLength() int {
	return len(*state.log)
}

// Returns the LogEntry at the specified index.
func (state *NodeState) Log(index int) LogEntry {
	return (*state.log)[index]
}

// Returns the node's CommitIndex.
func (state *NodeState) CommitIndex() int {
	return state.commitIndex
}

// Returns the node's LastApplied.
func (state *NodeState) LastApplied() int {
	return state.lastApplied
}

// Returns the node's NextIndex (only valid for leader nodes).
func (state *NodeState) NextIndex() map[int]int {
	return state.nextIndex
}

// Returns the node's MatchIndex (only valid for leader nodes).
func (state *NodeState) MatchIndex() map[int]int {
	return state.matchIndex
}
