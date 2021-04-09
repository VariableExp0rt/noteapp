package listing

import (
	"github.com/VariableExp0rt/dddexample/notes"
)

type Service interface {
	GetNotes() ([]notes.Note, error)
	GetNote(int) (notes.Note, error)
}

type service struct {
	nR notes.Repository
}

func NewService(nR notes.Repository) Service {
	return &service{nR}
}

func (s *service) GetNotes() ([]notes.Note, error) {
	return s.nR.GetAll()
}

func (s *service) GetNote(id int) (notes.Note, error) {
	return s.nR.Get(id)
}
