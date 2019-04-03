package store

//KVStore describes CRUD operations over KV store
type KVStore interface {
	Store(bucket string, k string, v interface{}) error
	Find(bucket string, k string) (interface{}, error)
	FindString(bucket, k string) (string, error)
}
