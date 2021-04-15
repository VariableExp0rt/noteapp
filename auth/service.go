package auth

type Service interface {
	ValidateUser(username, password string) error
	StoreNewUser(User) error
}

type service struct {
	uR Repository
}

func NewService(uR Repository) Service {
	return &service{uR}
}

func (s *service) ValidateUser(u, p string) error {
	if err := s.uR.Validate(u, p); err != nil {
		return err
	}
	return nil
}

func (s *service) StoreNewUser(u User) error {
	if err := s.uR.Store(u); err != nil {
		return err
	}
	return nil
}
