package storage

import (
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"github.com/VariableExp0rt/dddexample/notes"
	"github.com/boltdb/bolt"
)

type BoltStorage struct {
	DB *bolt.DB
}

//Move this to main
func NewBoltStorage(path string) (*BoltStorage, error) {

	fmt.Println("Opening...")
	db, err := bolt.Open(path, 600, &bolt.Options{Timeout: 1 * time.Second})
	if err != nil {
		return &BoltStorage{}, err
	}
	return &BoltStorage{db}, nil
}

func (s *BoltStorage) Add(n notes.Note) error {
	tx, err := s.DB.Begin(true)
	if err != nil {
		return fmt.Errorf("unable to return transaction: %v", err)
	}

	defer tx.Rollback()

	bkt, err := tx.CreateBucketIfNotExists([]byte("NOTES"))
	if err != nil {
		return err
	}

	id, err := bkt.NextSequence()
	if err != nil {
		return err
	}

	if buf, err := json.Marshal(n); err != nil {
		return err
	} else if err := bkt.Put([]byte(strconv.FormatInt(int64(id), 10)), buf); err != nil {
		return err
	}

	if err := tx.Commit(); err != nil {
		return err
	}

	return nil
}

func (s *BoltStorage) Get(id int) (notes.Note, error) {
	n := notes.Note{}

	err := s.DB.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("NOTES"))
		v := b.Get([]byte(strconv.FormatInt(int64(id), 10)))
		if err := json.Unmarshal(v, &n); err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return notes.Note{}, err
	}

	return n, nil
}

func (s *BoltStorage) GetAll() ([]notes.Note, error) {
	ns := []notes.Note{}

	err := s.DB.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("NOTES"))

		err := b.ForEach(func(k, v []byte) error {
			//For each note, append to list of notes above
			n := notes.Note{}
			if err := json.Unmarshal(v, &n); err != nil {
				return err
			}

			ns = append(ns, n)

			return nil
		})
		if err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return ns, nil
}

func (s *BoltStorage) Delete(id int) error {
	s.DB.Update(func(tx *bolt.Tx) error {
		return nil
	})

	return nil
}

func (s *BoltStorage) Update(n notes.Note) error {
	s.DB.Update(func(tx *bolt.Tx) error {
		return nil
	})

	return nil
}
