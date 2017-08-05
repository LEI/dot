// https://github.com/boltdb/bolt#mobile-use-iosandroid

package db

import (
	"fmt"
	"log"

	"github.com/boltdb/bolt"
)

func NewBoltDB(filepath string) *BoltDB {
	db, err := bolt.Open(filepath+"/demo.db", 0600, nil)
	if err != nil {
		log.Fatal(err)
	}

	return &BoltDB{db}
}

type BoltDB struct {
	db *bolt.DB
}

func (b *BoltDB) Path() string {
	return b.db.Path()
}

func (b *BoltDB) Close() {
	b.db.Close()
}

func (b *BoltDB) Put(name string, key string, value string) {
	b.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(name))
		err := b.Put([]byte(key), []byte(value))
		return err
	})
}

func (b *BoltDB) Get(name string, key string) string {
	var value string
	b.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(name))
		v := b.Get([]byte(key))
		// fmt.Printf("The answer is: %s\n", v)
		value = copy(v)
		return nil
	})
	return value
}

func (b *BoltDB) Delete(name string, key string) {
	b.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(name))
		err := b.Delete([]byte(key))
		return err
	})
}

// type BoltBucket struct {
// 	name string
// 	db *bolt.DB
// 	// bucket *bold.Bucket
// }

// func (b *BoltDB) NewBucket(name string) *BoltBucket {
// 	b.db.Update(func(tx *bolt.Tx) error {
// 		b, err := tx.CreateBucket([]byte(name))
// 		if err != nil {
// 			return fmt.Errorf("create bucket '%s': %s", name, err)
// 		}
// 		return nil
// 	})

// 	return &BoltBucket{name, b.db}
// }
