package state

import (
	"encoding/json"
	"strconv"

	"github.com/thomasylee/GoRaft/global"
)

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

type NodeState interface {
	SetCurrentTerm(newCurrentTerm int)
	CurrentTerm() int
	SetVotedFor(newVotedFor string)
	VotedFor() string
	AppendEntryToLog(index int, entry LogEntry) error
	LogLength() int
	Log(index int) LogEntry

	NodeStateMachine() StateMachine
	StorageStateMachine() StateMachine
	CommitIndex() int
	LastApplied() int
	NextIndex() map[int]int
	MatchIndex() map[int]int
}

/**
 * A NodeStateImpl contains both the persistent and volatile state that a node
 * needs to function in a Raft cluster.
 *
 * Note that currentTerm, votedFor, and log are not exported because they need
 * to be persistent, so to preserve consistency, they should be updated and
 * retrieved using the respective NodeState methods.
 */
type NodeStateImpl struct {
	// The latest term the node has seen.
	currentTerm int

	// The candidateId (UUID) that was voted for the current term.
	votedFor string

	// Log entries, each containing a command for the state machine and the term
	// when the entry was received from the leader.
	log *[]LogEntry

	// Node state and stored state are stored by different state machines.
	nodeStateMachine StateMachine
	storageStateMachine StateMachine

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

var Node NodeState
var initializedNodeState bool = false

/**
 * Returns the global NodeState variable, initializing it if it hasn't yet been
 * instantiated.
 */
func GetNodeState() NodeState {
	if Node != nil {
		return Node
	}
	Node = NewNodeState()
	return Node
}

/**
 * Returns a NodeState based on values in the node state Bolt database, using
 * default values if the database does not exist or have any values in it.
 */
func NewNodeState() NodeState {
	if initializedNodeState {
		return Node
	}
	initializedNodeState = true

	nodeStateMachine, err := NewBoltStateMachine("node_state.db")
	if err != nil {
		global.Log.Panic("Failed to initialize nodeStateMachine:", err.Error())
	}
	storageStateMachine, err := NewBoltStateMachine("state.db")
	if err != nil {
		global.Log.Panic("Failed to initialize storageStateMachine:", err.Error())
	}

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

	var node NodeState
	node = &NodeStateImpl{nodeStateMachine: nodeStateMachine, storageStateMachine: storageStateMachine}
	node.SetCurrentTerm(currentTermValue)
	node.SetVotedFor(votedForValue)
	node.(*NodeStateImpl).log = &logEntries

	return node
}

/**
 * Sets the current term in memory and in the node state machine.
 */
func (state *NodeStateImpl) SetCurrentTerm(newCurrentTerm int) {
	state.currentTerm = newCurrentTerm
	global.Log.Debugf("CurrentTerm updated: %d", newCurrentTerm)
	state.nodeStateMachine.Put(currentTerm, strconv.Itoa(newCurrentTerm))
}

/**
 * Returns the current term recognized by the node.
 */
func (state NodeStateImpl) CurrentTerm() int {
	return state.currentTerm;
}

/**
 * Sets VotedFor in memory and in the node state machine.
 */
func (state *NodeStateImpl) SetVotedFor(newVotedFor string) {
	state.votedFor = newVotedFor
	state.nodeStateMachine.Put(votedFor, newVotedFor)
}

/**
 * Returns the node's VotedFor.
 */
func (state NodeStateImpl) VotedFor() string {
	return state.votedFor;
}

/**
 * Append the log entry to the NodeState's log.
 * Note that this method does not do any safety checking to prevent overwriting
 * existing entries; that check should be done by the caller beforehand.
 */
func (state *NodeStateImpl) AppendEntryToLog(index int, entry LogEntry) error {
	jsonValue, err := json.Marshal(entry)
	if err != nil {
		return err
	}

	err = state.nodeStateMachine.Put(strconv.Itoa(index), string(jsonValue))
	if err != nil {
		return err
	}

	*state.log = append(*state.log, entry)
	return nil
}

/**
 * Returns the number of entries in the node's log.
 */
func (state NodeStateImpl) LogLength() int {
	return len(*state.log)
}

/**
 * Returns the LogEntry at the specified index.
 */
func (state NodeStateImpl) Log(index int) LogEntry {
	return (*state.log)[index]
}

func (state *NodeStateImpl) NodeStateMachine() StateMachine {
	return state.nodeStateMachine
}

func (state *NodeStateImpl) StorageStateMachine() StateMachine {
	return state.storageStateMachine
}

func (state *NodeStateImpl) CommitIndex() int {
	return state.commitIndex
}

func (state *NodeStateImpl) LastApplied() int {
	return state.lastApplied
}

func (state *NodeStateImpl) NextIndex() map[int]int {
	return state.nextIndex
}

func (state *NodeStateImpl) MatchIndex() map[int]int {
	return state.matchIndex
}
