package notes

import (
	"errors"
	"time"
)

type Note struct {
	ID          int       `json:"id"`
	Name        string    `json:"name"`
	Content     string    `json:"content"`
	CreatedTime time.Time `json:"created_at"`
}

var (
	ErrNoteNotFound  = errors.New("unknown note")
	ErrDuplicateNote = errors.New("note already exists")
)

type Repository interface {
	GetAll() ([]Note, error)
	Get(ID int) (Note, error)
	Add(n Note) error
	Delete(ID int) error
	Update(id int, n Note) error
}
