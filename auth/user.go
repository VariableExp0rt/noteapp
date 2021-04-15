package auth

import "errors"

var (
	ErrUserNotFound = errors.New("Invalid username and/or password.")
)

type User struct {
	ID           int
	Username     string
	Password     string
	EmailAddress string
}

type Repository interface {
	Validate(username, password string) error
	Store(User) error
}
