package updating

import "github.com/VariableExp0rt/dddexample/notes"

type Service interface {
	Update(id int, n notes.Note) error
}

type service struct {
	nR notes.Repository
}

func NewService(nR notes.Repository) Service {
	return &service{nR}
}

func (s *service) Update(id int, n notes.Note) error {
	if err := s.nR.Update(id, n); err != nil {
		return err
	}
	return nil
}
