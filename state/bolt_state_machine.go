package state

import (
	"bytes"
	"encoding/json"
	"log"
	"strconv"
	"time"

	"github.com/boltdb/bolt"

	"github.com/thomasylee/GoRaft/errors"
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
			log.Fatal("Create bucket error:", err)
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
 *
 */
func (boltSM BoltStateMachine) RetrieveLogEntries(lastIndex int) (entries []LogEntry, err error) {
	err = boltSM.db.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(bucket)).Cursor()

		min := []byte("0")
		max := []byte(strconv.Itoa(lastIndex))

		var entries []LogEntry
		var err error
		for k, v := bucket.Seek(min); k != nil && bytes.Compare(k, max) <= 0; k, v = bucket.Next() {
			var entry LogEntry
			err = json.Unmarshal(v, &entry)
			if errors.HandleError("Failed to unmarshall JSON LogEntry:", err, false) {
				return err
			}

			intKey, err := strconv.Atoi(string(k))
			if errors.HandleError("Failed to convert LogEntry key to int:", err, false) {
				return err
			}
			entries[intKey] = entry
		}
		return nil
	})
	return
}
