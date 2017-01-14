package state

import (
	"strconv"
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
type NodeState struct {
	// The latest term the node has seen.
	CurrentTerm int

	// The candidateId that was voted for the current term.
	VotedFor int

	// Log entries, each containing a command for the state machine and the term
	// when the entry was received from the leader.
	Log []LogEntry

	// Node state and state machine state are stored in separate Bolt
	// databases.
	NodeState *BoltWrapper
	StateMachine *BoltWrapper

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

/**
 * Return a NodeState based on values in the node state Bolt database, using
 * default values if the database does not exist or have any values in it.
 */
func NewNodeState() NodeState {
	nodeStateBolt, _ := NewBoltWrapper("node_state.db")
	stateMachineBolt, _ := NewBoltWrapper("state.db")

	var currentTermValue int
	if nodeStateBolt.Get(currentTerm) == "" {
		nodeStateBolt.Put(currentTerm, "0")
		currentTermValue = 0
	} else {
		currentTermValue, _ = strconv.Atoi(nodeStateBolt.Get(currentTerm))
	}

	votedForValue, _ := strconv.Atoi(nodeStateBolt.Get(votedFor))

	nodeState := NodeState{NodeState: nodeStateBolt, StateMachine: stateMachineBolt, CurrentTerm: currentTermValue, VotedFor: votedForValue}

	return nodeState
}
