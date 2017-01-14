package state

import (
	"time"

	"github.com/boltdb/bolt"
)

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
	return
}

/**
 * Writes a key-value pair to the Bolt database.
 */
func (boltWrapper *BoltWrapper) Put(key string, value string) {
	boltWrapper.db.Update(func(tx *bolt.Tx) error {
		bucket, err := tx.CreateBucketIfNotExists([]byte("State"))
		if err != nil {
			return err
		}

		err = bucket.Put([]byte(key), []byte(value))
		return err
	})
}

/**
 * Returns the value of the specified key stored in the Bolt database.
 */
func (boltWrapper *BoltWrapper) Get(key string) (value string) {
	boltWrapper.db.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte("State"))
		data := bucket.Get([]byte(key))
		if data != nil {
			valueBytes := make([]byte, len(data))
			copy(valueBytes, data)
			value = string(valueBytes)
		}
		return nil
	})

	return
}
