package users

import (
	"errors"
	"fmt"
)

var ErrFirstNameRequired = errors.New("first name is required")
var ErrLastNameRequired = errors.New("last name is required")

type ErrUserNotFound struct {
	UserId string
}

func (e ErrUserNotFound) Error() string {
	return fmt.Sprintf("User '%s' doesent exist", e.UserId)
}
