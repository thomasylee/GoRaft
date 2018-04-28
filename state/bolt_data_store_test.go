package state

import (
	"os"
	"testing"
)

func Test_PutAndGet_WithCreatedBucket_PutsAndGetsCorrectly(t *testing.T) {
	dataStoreFile := "test_temp_db"

	bolt, err := NewBoltDataStore(dataStoreFile)
	if err != nil {
		t.Error("Creating BoltDataStore failed:", err)
	}

	bolt.Put("key", "value")

	value, _ := bolt.Get("key")
	if value != "value" {
		t.Error("Retrieved value did not equal expected:", "value", value)
	}

	os.Remove(dataStoreFile)
}
