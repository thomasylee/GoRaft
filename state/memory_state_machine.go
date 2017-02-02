package state

import (
	"encoding/json"
	"strconv"
)

type MemoryStateMachine struct {
	values map[string]string
}

func NewMemoryStateMachine() MemoryStateMachine {
	return MemoryStateMachine{values: make(map[string]string)}
}

func (sm MemoryStateMachine) Put(key string, value string) error {
	sm.values[key] = value
	return nil
}

func (sm MemoryStateMachine) Get(key string) (string, error) {
	return sm.values[key], nil
}

func (sm MemoryStateMachine) RetrieveLogEntries(firstIndex int, lastIndex int) ([]LogEntry, error) {
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
