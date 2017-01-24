package state

import (
	"encoding/json"
	"strconv"

	"github.com/op/go-logging"
)

/**
 * The leveled Logger to use in the state package.
 */
var Log *logging.Logger

/**
 * A LogEntry represents a state machine command (Put {Key},{Value}) and the
 * term when the entry was received by the leader.
 */
type LogEntry struct {
	Key string
	Value string
	Term int
}

/**
 * Define constants for important keys in the Bolt database.
 */
const (
	currentTerm string = "CurrentTerm"
	votedFor string = "VotedFor"
	logEntries string = "LogEntries"
)

/**
 * A NodeState contains both the persistent and volatile state that a node
 * needs to function in a Raft cluster.
 *
 * Note that currentTerm, votedFor, and log are not exported because they need
 * to be persistent, so to preserve consistency, they should be updated and
 * retrieved using the respective NodeState methods.
 */
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
	CommitIndex int

	// Index of the highest log entry applied to this node's storage state machine.
	LastApplied int

	// (Leader only) For each node, the index of the next log entry to send to.
	NextIndex map[int]int

	// (Leader only) For each node, the index of the highest log entry known to be replicated on
	// the node.
	MatchIndex map[int]int
}

var nodeState NodeState
var initializedNodeState bool = false

/**
 * Return a NodeState based on values in the node state Bolt database, using
 * default values if the database does not exist or have any values in it.
 */
func GetNodeState(logger *logging.Logger) NodeState {
	if initializedNodeState {
		return nodeState
	}
	initializedNodeState = true

	Log = logger

	nodeStateMachine, err := NewBoltStateMachine("node_state.db")
	if err != nil {
		Log.Panic("Failed to initialize nodeStateMachine:", err.Error())
	}
	storageStateMachine, err := NewBoltStateMachine("state.db")
	if err != nil {
		Log.Panic("Failed to initialize storageStateMachine:", err.Error())
	}

	var currentTermValue int
	retrievedCurrentTerm, err := nodeStateMachine.Get(currentTerm)
	if err != nil {
		Log.Panic("Failed to retrieve CurrentTerm:", err.Error())
	}
	if retrievedCurrentTerm == "" {
		currentTermValue = 0
	} else {
		currentTermValue, err = strconv.Atoi(retrievedCurrentTerm)
		if err != nil {
			Log.Panic("Failed to convert CurrentTerm to int:", err.Error())
		}
	}

	votedForValue, err := nodeStateMachine.Get(votedFor)
	if err != nil {
		Log.Panic("Failed to retrieve VotedFor:", err.Error())
	}

	logEntries, err := nodeStateMachine.RetrieveLogEntries(0, 1000000)
	if err != nil {
		Log.Panic("Failed to retrieve log entries:", err.Error())
	}

	nodeState = NodeState{NodeStateMachine: nodeStateMachine, StorageStateMachine: storageStateMachine}
	nodeState.SetCurrentTerm(currentTermValue)
	nodeState.SetVotedFor(votedForValue)
	nodeState.log = &logEntries

	return nodeState
}

/**
 * Sets the current term in memory and in the node state machine.
 */
func (state *NodeState) SetCurrentTerm(newCurrentTerm int) {
	state.currentTerm = newCurrentTerm
	Log.Debugf("CurrentTerm updated: %d", newCurrentTerm)
	state.NodeStateMachine.Put(currentTerm, strconv.Itoa(newCurrentTerm))
}

/**
 * Returns the current term recognized by the node.
 */
func (state NodeState) CurrentTerm() int {
	return state.currentTerm;
}

/**
 * Sets VotedFor in memory and in the node state machine.
 */
func (state NodeState) SetVotedFor(newVotedFor string) {
	state.votedFor = newVotedFor
	state.NodeStateMachine.Put(votedFor, newVotedFor)
}

/**
 * Returns the node's VotedFor.
 */
func (state NodeState) VotedFor() string {
	return state.votedFor;
}

/**
 * Append the log entry to the NodeState's log.
 * Note that this method does not do any safety checking to prevent overwriting
 * existing entries; that check should be done by the caller beforehand.
 */
func (state NodeState) AppendEntryToLog(index int, entry LogEntry) error {
	jsonValue, err := json.Marshal(entry)
	if err != nil {
		return err
	}

	err = state.NodeStateMachine.Put(strconv.Itoa(index), string(jsonValue))
	if err == nil {
		*state.log = append(*state.log, entry)
		return nil
	} else {
		return err
	}
}

/**
 * Returns the number of entries in the node's log.
 */
func (state NodeState) LogLength() int {
	return len(*state.log)
}

/**
 * Returns the LogEntry at the specified index.
 */
func (state NodeState) Log(index int) LogEntry {
	return (*state.log)[index]
}
