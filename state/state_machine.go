package state

type StateMachine interface {
	Put(string, string) error
	Get(string) (string, error)
	RetrieveLogEntries(int) ([]LogEntry, error)
}
