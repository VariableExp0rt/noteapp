package auth

import "errors"

var (
	ErrUserNotFound            = errors.New("Invalid username and/or password.")
	ErrInvalidPasswordCriteria = errors.New("Password does not meet complexity criteria.")
)

type UserLoginReq struct {
	Username string
	Password string
}

type UserSignUpReq struct {
	Username        string
	NewPassword     string
	ConfirmPassword string
}
type Repository interface {
	ValidateUser(username, password string) error
	StoreNewUser(UserSignUpReq) error
}
