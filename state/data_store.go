package state

// DataStore represents any kind of key-value database.
type DataStore interface {
	Put(string, string) error
	Get(string) (string, error)
	RetrieveLogEntries(int, int) ([]LogEntry, error)
}
