package storage

import (
	"bytes"
	"encoding/json"
	"fmt"
	"math/rand"
	"time"

	"github.com/VariableExp0rt/dddexample/auth"
	"github.com/VariableExp0rt/dddexample/notes"
	"github.com/boltdb/bolt"
	hash "github.com/mitchellh/hashstructure/v2"
)

type BoltStorage struct {
	DB *bolt.DB
}

func IDGenerator() int {

	s := rand.NewSource(time.Now().UnixNano())
	r := rand.New(s)

	return r.Intn(100000)
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

	//use this package instead to make a hashmap as the key it is stored as, based on the ID
	//https://github.com/mitchellh/hashstructure
	//below in the Get/GetAll/Delete/Update, we can use the given ID, hash it with this package
	//then call a bkt.Get() to find the hash within the key
	n.ID = IDGenerator()
	n.CreatedTime = time.Now().UTC()

	h, err := hash.Hash(n.ID, hash.FormatV2, nil)
	if err != nil {
		return err
	}

	if buf, err := json.Marshal(n); err != nil {
		return err
	} else if err := bkt.Put([]byte(fmt.Sprint(h)), buf); err != nil {
		return err
	}

	if err := tx.Commit(); err != nil {
		return err
	}

	return nil
}

func (s *BoltStorage) Get(id int) (notes.Note, error) {
	var n notes.Note

	err := s.DB.View(func(tx *bolt.Tx) error {

		h, err := hash.Hash(id, hash.FormatV2, nil)
		if err != nil {
			return err
		}

		b := tx.Bucket([]byte("NOTES"))
		result := b.Get([]byte(fmt.Sprint(h)))
		if result == nil {
			return notes.ErrNoteNotFound
		}

		if err := json.Unmarshal(result, &n); err != nil {
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
	err := s.DB.Update(func(tx *bolt.Tx) error {

		bkt := tx.Bucket([]byte("NOTES"))

		h, err := hash.Hash(id, hash.FormatV2, nil)
		if err != nil {
			return err
		}

		if err := bkt.Delete([]byte(fmt.Sprint(h))); err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return err
	}

	return nil
}

func (s *BoltStorage) Update(id int, n notes.Note) error {
	var err error

	err = s.DB.Update(func(tx *bolt.Tx) error {

		//logic for updating here
		bkt := tx.Bucket([]byte("NOTES"))

		if err := bkt.ForEach(func(k, v []byte) error {
			var note notes.Note

			if err := json.NewDecoder(bytes.NewReader(v)).Decode(&note); err != nil {
				return err
			}

			if note.ID == id {
				n.ID = note.ID
				n.CreatedTime = note.CreatedTime
				if buf, err := json.Marshal(n); err != nil {
					return err
				} else if err := bkt.Put(k, buf); err != nil {
					return err
				}
			}

			return nil

		}); err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return err
	}

	return nil
}

func (s *BoltStorage) Validate(username, password string) error {

	//take username and password and compare to stored values
	//if match, return no err and the handler creates token
	if err := s.DB.View(func(tx *bolt.Tx) error {
		bkt, err := tx.CreateBucketIfNotExists([]byte("USERS"))
		if err != nil {
			return err
		}
		cur := bkt.Cursor()

		var u auth.User

		for {
			_, v := cur.Next()
			if err := json.Unmarshal(v, &u); err != nil {
				return err
			}
			if username == u.Username && password == u.Password {
				return nil
			} else if username != u.Username && password != u.Password || username == u.Username && password != u.Password {
				break
			}
			return auth.ErrUserNotFound
		}

		return nil
	}); err != nil {
		return err
	}

	return nil
}

func (s *BoltStorage) Store(auth.User) error {
	return nil
}
