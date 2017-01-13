package state

import (
	"time"

	"github.com/boltdb/bolt"
)

type Fsm struct {
	db *bolt.DB
}

func NewFsm() (fsm *Fsm, err error) {
	fsm = &Fsm{}
	fsm.db, err = bolt.Open("state.db", 0600, &bolt.Options{Timeout: 1 * time.Second})
	return
}

func (fsm *Fsm) ApplyState(key string, value string) {
	fsm.db.Update(func(tx *bolt.Tx) error {
		bucket, err := tx.CreateBucketIfNotExists([]byte("State"))
		if err != nil {
			return err
		}

		err = bucket.Put([]byte(key), []byte(value))
		return err
	})
}

func (fsm *Fsm) RetrieveState(key string) (value string, err error) {
	fsm.db.View(func(tx *bolt.Tx) error {
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
