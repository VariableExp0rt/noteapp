package deleting

import (
	"github.com/VariableExp0rt/dddexample/notes"
)

type Service interface {
	DeleteNote(int) error
}

type service struct {
	nR notes.Repository
}

func NewService(nR notes.Repository) Service {
	return &service{nR}
}

func (s *service) DeleteNote(id int) error {
	if err := s.nR.Delete(id); err != nil {
		return err
	}
	return nil
}
