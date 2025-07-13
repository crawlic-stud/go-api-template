package models

import "validation-api/internal/util/validation"

type User struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func (u User) Validate() error {
	return validation.NewValidator(u).
		Add(u.Username != "", "username must not be empty").
		Add(u.Password != "", "password must not be empty").
		Validate()
}
