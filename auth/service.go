package auth

type Service interface {
	ValidateUser(username, password string) error
	StoreNewUser(UserSignUpReq) error
}

type service struct {
	uR Repository
}

func NewService(uR Repository) Service {
	return &service{uR}
}

func (s *service) ValidateUser(u, p string) error {
	if err := s.uR.ValidateUser(u, p); err != nil {
		return err
	}
	return nil
}

func (s *service) StoreNewUser(u UserSignUpReq) error {
	if err := s.uR.StoreNewUser(u); err != nil {
		return err
	}
	return nil
}
