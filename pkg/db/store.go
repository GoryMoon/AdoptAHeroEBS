package db

import (
	"github.com/dgraph-io/badger/v3"
	"google.golang.org/protobuf/proto"
	"time"
)

type Store struct {
	Db   *KvDB
	Name string
}

func (s Store) getKey(key string) []byte {
	return []byte(s.Name + "/" + key)
}

func (s *Store) GetValue(key string, out proto.Message) error {
	db, err := s.Db.GetDB()
	if err != nil {
		return err
	}

	// Read the data from the db
	return db.View(func(txn *badger.Txn) error {
		item, err := txn.Get(s.getKey(key))
		if err != nil {
			return err
		}

		// Read the value and convert it to a proto structure
		return item.Value(func(val []byte) error {
			return proto.Unmarshal(val, out)
		})
	})
}

func (s *Store) SetValue(key string, value proto.Message, ttl time.Duration) error {
	db, err := s.Db.GetDB()
	if err != nil {
		return err
	}

	// Convert data to bytes
	data, err := proto.Marshal(value)
	if err != nil {
		return err
	}

	// Save message bytes in the db
	return db.Update(func(txn *badger.Txn) error {
		entry := badger.NewEntry(s.getKey(key), data)

		if ttl.Seconds() > 0 {
			entry.WithTTL(ttl)
		}

		return txn.SetEntry(entry)
	})
}

func (s *Store) BulkSetValue(keys []string, fn func(key string) proto.Message, ttl time.Duration) error {
	db, err := s.Db.GetDB()
	if err != nil {
		return err
	}

	// Start a new batch to improve bulk settings
	wb := db.NewWriteBatch()
	defer wb.Cancel()

	for _, v := range keys {
		data, err := proto.Marshal(fn(v))
		if err != nil {
			return err
		}

		entry := badger.NewEntry(s.getKey(v), data)

		// Don't force a ttl on the value
		if ttl.Seconds() > 0 {
			entry.WithTTL(ttl)
		}

		err = wb.SetEntry(entry)
		if err != nil {
			return err
		}
	}

	return wb.Flush()
}

func (s *Store) Delete(key string) error {
	db, err := s.Db.GetDB()
	if err != nil {
		return err
	}

	dbKey := s.getKey(key)
	return db.Update(func(txn *badger.Txn) error {
		// Gets the key to verify it exists before trying to delete it
		_, err := txn.Get(dbKey)
		if err != nil {
			return err
		}

		return txn.Delete(dbKey)
	})
}

func (s *Store) BulkDelete(keys []string) error {
	db, err := s.Db.GetDB()
	if err != nil {
		return err
	}

	wb := db.NewWriteBatch()
	defer wb.Cancel()

	for _, v := range keys {
		err := wb.Delete(s.getKey(v))
		if err != nil {
			return err
		}
	}

	return wb.Flush()
}
