package state

import (
	"log"
	"time"

	"github.com/boltdb/bolt"
)

// Key-value pairs are stored in the State bucket for all Bolt databases.
const bucket string = "State"

/**
 * A BoltWrapper has a Bolt database that it uses to store key-value pairs and
 * retrieve values by their keys.
 */
type BoltWrapper struct {
	db *bolt.DB
}

/**
 * Returns a new instance of the BoltWrapper type.
 */
func NewBoltWrapper(dbFile string) (boltWrapper *BoltWrapper, err error) {
	boltWrapper = &BoltWrapper{}
	boltWrapper.db, err = bolt.Open(dbFile, 0600, &bolt.Options{Timeout: 1 * time.Second})

	boltWrapper.CreateBucketIfNotExists(bucket)

	return
}

/**
 * Ensures that a bucket with the given name exists by creating it if it doesn't
 * already exist.
 */
func (boltWrapper *BoltWrapper) CreateBucketIfNotExists(name string) error {
	return boltWrapper.db.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists([]byte(bucket))
		if err != nil {
			log.Fatal("Create bucket error: %s", err)
			return err
		}
		return nil
	})
}

/**
 * Writes a key-value pair to the Bolt database.
 */
func (boltWrapper *BoltWrapper) Put(key string, value string) {
	boltWrapper.db.Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(bucket))

		err := bucket.Put([]byte(key), []byte(value))
		return err
	})
}

/**
 * Returns the value of the specified key stored in the Bolt database.
 */
func (boltWrapper *BoltWrapper) Get(key string) string {
	var value string

	boltWrapper.db.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte("State"))
		log.Println("Bucket:", bucket)
		data := bucket.Get([]byte(key))
		if data != nil {
			valueBytes := make([]byte, len(data))
			copy(valueBytes, data)
			value = string(valueBytes)
		}
		return nil
	})

	return value
}
