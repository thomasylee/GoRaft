package state

import (
	"strconv"
)

type NodeStateTestImpl struct {
	currentTerm int
	nodeStateMachine StateMachine
	storageStateMachine StateMachine
	logLength int
}

func NewNodeStateTestImpl() *NodeStateTestImpl {
	return &NodeStateTestImpl{currentTerm: 0, nodeStateMachine: newStateMachineTestImpl(), storageStateMachine: newStateMachineTestImpl(), logLength: 0}
}

func (nodeState *NodeStateTestImpl) SetCurrentTerm(newCurrentTerm int) {
	nodeState.currentTerm = newCurrentTerm
}
func (nodeState NodeStateTestImpl) CurrentTerm() int {
	return nodeState.currentTerm
}
func (nodeState NodeStateTestImpl) SetVotedFor(newVotedFor string) {}
func (nodeState NodeStateTestImpl) VotedFor() string { return "" }
func (nodeState NodeStateTestImpl) AppendEntryToLog(index int, entry LogEntry) error {
	nodeState.nodeStateMachine.Put(strconv.Itoa(index), "")
	nodeState.logLength++
	return nil
}
func (nodeState NodeStateTestImpl) LogLength() int { return nodeState.logLength }
func (nodeState NodeStateTestImpl) Log(index int) LogEntry { return LogEntry{} }

func (nodeState NodeStateTestImpl) NodeStateMachine() StateMachine {
	return nodeState.nodeStateMachine
}
func (nodeState NodeStateTestImpl) StorageStateMachine() StateMachine {
	return nodeState.storageStateMachine
}
func (nodeState NodeStateTestImpl) CommitIndex() int { return 0 }
func (nodeState NodeStateTestImpl) LastApplied() int { return 0 }
func (nodeState NodeStateTestImpl) NextIndex() map[int]int { return map[int]int{} }
func (nodeState NodeStateTestImpl) MatchIndex() map[int]int { return map[int]int{} }
