package adding

import (
	"github.com/VariableExp0rt/dddexample/notes"
)

type Service interface {
	AddNote(...notes.Note)
}

type service struct {
	nR notes.Repository
}

func NewService(nR notes.Repository) Service {
	return &service{nR}
}

func (s *service) AddNote(n ...notes.Note) {
	for _, note := range n {
		_ = s.nR.Add(note)
	}
}
