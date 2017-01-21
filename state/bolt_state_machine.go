package state

import (
	"encoding/json"
	"strconv"
	"time"

	"github.com/boltdb/bolt"
)

// Key-value pairs are stored in the State bucket for all Bolt databases.
const bucket string = "State"

/**
 * A BoltStateMachine has a Bolt database that it uses to store key-value pairs and
 * retrieve values by their keys.
 */
type BoltStateMachine struct {
	db *bolt.DB
}

/**
 * Returns a new instance of the BoltStateMachine type.
 */
func NewBoltStateMachine(dbFile string) (boltSM *BoltStateMachine, err error) {
	boltSM = &BoltStateMachine{}
	boltSM.db, err = bolt.Open(dbFile, 0600, &bolt.Options{Timeout: 1 * time.Second})

	boltSM.CreateBucketIfNotExists(bucket)

	return
}

/**
 * Ensures that a bucket with the given name exists by creating it if it doesn't
 * already exist.
 */
func (boltSM *BoltStateMachine) CreateBucketIfNotExists(name string) error {
	return boltSM.db.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists([]byte(bucket))
		if err != nil {
			Log.Error("Create bucket error:", err)
			return err
		}
		return nil
	})
}

/**
 * Writes a key-value pair to the Bolt database.
 */
func (boltSM BoltStateMachine) Put(key string, value string) error {
	return boltSM.db.Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(bucket))

		err := bucket.Put([]byte(key), []byte(value))
		return err
	})
}

/**
 * Returns the value of the specified key stored in the Bolt database.
 */
func (boltSM BoltStateMachine) Get(key string) (string, error) {
	var value string

	err := boltSM.db.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(bucket))
		data := bucket.Get([]byte(key))
		if data != nil {
			valueBytes := make([]byte, len(data))
			copy(valueBytes, data)
			value = string(valueBytes)
		}
		return nil
	})

	return value, err
}

/**
 * Returns log entries within the specified key range.
 *
 * TODO: Use a more efficient method than querying each index one at a time.
 */
func (boltSM BoltStateMachine) RetrieveLogEntries(firstIndex int, lastIndex int) ([]LogEntry, error) {
	entries := []LogEntry{}
	for i := firstIndex; i <= lastIndex; i++ {
		jsonValue, err := boltSM.Get(strconv.Itoa(i))
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
