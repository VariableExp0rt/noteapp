package deleting

import (
	"fmt"

	"github.com/VariableExp0rt/dddexample/notes"
)

type Service interface {
	DeleteNote(int)
}

type service struct {
	nR notes.Repository
}

func NewService(nR notes.Repository) Service {
	return &service{nR}
}

func (s *service) DeleteNote(id int) {
	if err := s.nR.Delete(id); err != nil {
		fmt.Printf("Error deleting note with ID %v: %v", id, err)
		return
	}
}
