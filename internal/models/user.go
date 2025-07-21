package models

import "template-api/internal/util/validation"

type LoginUser struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func (u LoginUser) Validate() error {
	return validation.NewValidator(u).
		Add(u.Username != "", "username must not be empty").
		Add(u.Password != "", "password must not be empty").
		Validate()
}
