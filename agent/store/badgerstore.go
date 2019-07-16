package store

import "github.com/dgraph-io/badger"

//KVStore describes CRUD operations over KV store
type BadgerStore struct {
	db *badger.DB
}

func NewBadgerStore(db *badger.DB) *BadgerStore {
	return &BadgerStore{db: db}
}

func (bs *BadgerStore) Store(bucket string, k string, v interface{}) error {
	return nil
}
func (bs *BadgerStore) Find(bucket string, k string) (interface{}, error) {
	return "", nil
}
func (bs *BadgerStore) FindString(bucket, k string) (string, error) {
	return "", nil
}
