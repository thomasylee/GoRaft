package state

import (
	"encoding/json"
	"strconv"
)

// MemoryDataStore stores key-value pairs in memory.
type MemoryDataStore struct {
	values map[string]string
}

// NewMemoryDataStore constructs a new empty MemoryDataStore.
func NewMemoryDataStore() MemoryDataStore {
	return MemoryDataStore{values: make(map[string]string)}
}

// Put adds a new key-value pair to the data store.
func (sm MemoryDataStore) Put(key string, value string) error {
	sm.values[key] = value
	return nil
}

// Get retrieves a value based on its key.
func (sm MemoryDataStore) Get(key string) (string, error) {
	return sm.values[key], nil
}

// RetrieveLogEntries returns the LogEntries found between the specified indices.
func (sm MemoryDataStore) RetrieveLogEntries(firstIndex int, lastIndex int) ([]LogEntry, error) {
	entries := []LogEntry{}
	for i := firstIndex; i <= lastIndex; i++ {
		jsonValue, err := sm.Get(strconv.Itoa(i))
		if err != nil {
			return nil, err
		} else if jsonValue == "" {
			// As soon as we reach an empty record, we know there are no more entries.
			return entries, nil
		}

		var entry LogEntry
		err = json.Unmarshal([]byte(jsonValue), &entry)
		if err != nil {
			return nil, err
		}
		entries = append(entries, entry)
	}
	return entries, nil
}
