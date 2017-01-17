package state

import (
	"encoding/json"
	"strconv"

	"github.com/thomasylee/GoRaft/errors"
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
	lastIndex string = "LastIndex"
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

	// The candidateId that was voted for the current term.
	votedFor int

	// Log entries, each containing a command for the state machine and the term
	// when the entry was received from the leader.
	log []LogEntry

	// The index of the latest value added to the log.
	// Note that this is not in the Raft manual; it exists purely to help find
	// the last index without expensive database queries.
	LastIndex int

	// Node state and stored state are stored by different state machines.
	NodeStateMachine StateMachine
	StorageStateMachine StateMachine

	// Index of the highest log entry known to be committed.
	CommitIndex int

	// Index of the highest log entry applied to the local state machine.
	LastApplied int

	// (Leader only) For each node, the index of the next log entry to send to.
	NextIndex map[int]int

	// (Leader only) For each node, the index of the highest log entry known to be replicated on
	// the node.
	MatchIndex map[int]int
}

var nodeState NodeState

/**
 * Return a NodeState based on values in the node state Bolt database, using
 * default values if the database does not exist or have any values in it.
 */
func GetNodeState() NodeState {
	if &nodeState != nil {
		return nodeState
	}

	nodeStateMachine, err := NewBoltStateMachine("node_state.db")
	errors.HandleError("Failed to initialize nodeStateMachine:", err, true)
	storageStateMachine, err := NewBoltStateMachine("state.db")
	errors.HandleError("Failed to initialize storageStateMachine:", err, true)

	var currentTermValue int
	retrievedCurrentTerm, err := nodeStateMachine.Get(currentTerm)
	errors.HandleError("Failed to retrieve CurrentTerm:", err, true)
	if retrievedCurrentTerm == "" {
		currentTermValue = 0
	} else {
		currentTermValue, err = strconv.Atoi(retrievedCurrentTerm)
		errors.HandleError("Failed to convert CurrentTerm to int:", err, true)
	}

	retrievedVotedFor, err := nodeStateMachine.Get(votedFor)
	errors.HandleError("Failed to retrieve VotedFor:", err, true)
	votedForValue, err := strconv.Atoi(retrievedVotedFor)
	errors.HandleError("Failed to convert VotedFor to int:", err, true)

	retrievedLastIndex, err := nodeStateMachine.Get(lastIndex)
	errors.HandleError("Failed to retrieve LastIndex:", err, true)
	lastIndexValue, err := strconv.Atoi(retrievedLastIndex)
	errors.HandleError("Failed to convert LastIndex to int:", err, true)

	logEntries, err := nodeStateMachine.RetrieveLogEntries(lastIndexValue)
	errors.HandleError("Failed to retrieve log entries:", err, true)

	nodeState := NodeState{NodeStateMachine: nodeStateMachine, StorageStateMachine: storageStateMachine}
	nodeState.SetCurrentTerm(currentTermValue)
	nodeState.SetVotedFor(votedForValue)
	nodeState.log = logEntries

	return nodeState
}

/**
 * Append log entries to the NodeState's log.
 * Note that this method does not do any safety checking to prevent overwriting
 * existing entries; that check should be done by the caller beforehand.
 */
func (state NodeState) AppendEntry(index int, term int, key string, value string) {
	state.log[index] = LogEntry{Key: key, Value: value, Term: term}
	// state.NodeStateMachine.Put(strconv.Itoa(index), jsonEntry)
}

func (state NodeState) SetCurrentTerm(newCurrentTerm int) {
	state.currentTerm = newCurrentTerm
	state.NodeStateMachine.Put(currentTerm, strconv.Itoa(newCurrentTerm))
}

func (state NodeState) CurrentTerm() int {
	return state.currentTerm;
}

func (state NodeState) SetVotedFor(newVotedFor int) {
	state.votedFor = newVotedFor
	state.NodeStateMachine.Put(votedFor, strconv.Itoa(newVotedFor))
}

func (state NodeState) VotedFor() int {
	return state.votedFor;
}

func (state NodeState) AppendEntryToLog(entry LogEntry) error {
	jsonValue, err := json.Marshal(entry)
	if err != nil {
		return err
	}
	err = state.NodeStateMachine.Put(entry.Key, string(jsonValue))
	if err == nil {
		state.log[state.LastIndex + 1] = entry
		return nil
	} else {
		return err
	}
}

func (state NodeState) AllLog() []LogEntry {
	return state.log;
}

func (state NodeState) Log(index int) LogEntry {
	return state.log[index]
}

/**
 * Saves the NodeState's persistentState to the nodeStateDb Bolt database.
 *
func (state *NodeState) SavePersistentState() error {
	jsonValue, err := json.Marshal(state.Persistent)
	if err != nil {
		log.Fatal("Failed to save node state:", err.Error())
		return err
	}

	state.NodeStateDb.Put("NodeState", jsonValue)
	return
}
*/
