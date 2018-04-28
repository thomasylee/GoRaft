package state

import (
	"testing"
)

func Test_RetrieveLogEntries_WhenNoEntriesExist_ReturnsEmptySlice(t *testing.T) {
	entries, err := NewMemoryDataStore().RetrieveLogEntries(1, 5)
	if err != nil {
		t.Fatal(err)
	}
	if len(entries) != 0 {
		t.Error("Number of returned entries is not 0:", len(entries))
	}
}

func Test_RetrieveLogEntries_WhenWholeRangeExists_ReturnsLogEntriesInRange(t *testing.T) {
	var tests = []struct {
		key     string
		entry   LogEntry
		jsonRep string
	}{
		{"1", LogEntry{"a", "A", 0}, "{\"Key\":\"a\",\"Value\":\"A\",\"Term\":0}"},
		{"2", LogEntry{"b", "B", 0}, "{\"Key\":\"b\",\"Value\":\"B\",\"Term\":0}"},
		{"3", LogEntry{"c", "C", 1}, "{\"Key\":\"c\",\"Value\":\"C\",\"Term\":1}"},
		{"4", LogEntry{"d", "D", 2}, "{\"Key\":\"d\",\"Value\":\"D\",\"Term\":2}"},
	}

	stateMachine := NewMemoryDataStore()
	for _, test := range tests {
		stateMachine.Put(test.key, test.jsonRep)
	}

	retrievedEntries, err := stateMachine.RetrieveLogEntries(1, len(tests))
	if err != nil {
		t.Fatal(err)
	}
	if len(retrievedEntries) != len(tests) {
		t.Fatalf("Number of returned entries is not %d: %d", len(tests), len(retrievedEntries))
	}
	for i := 0; i < len(tests); i++ {
		if retrievedEntries[i] != tests[i].entry {
			t.Errorf("The retrieved entry at index %d does not match expected %v: %v", i, tests[i].entry, retrievedEntries[i])
		}
	}
}
