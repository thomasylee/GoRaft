package state

type DataStore interface {
	Put(string, string) error
	Get(string) (string, error)
	RetrieveLogEntries(int, int) ([]LogEntry, error)
}
