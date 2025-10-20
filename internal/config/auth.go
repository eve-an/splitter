package config

import "errors"

type Auth struct {
	Username string
	Password string
}

func (a Auth) Validate() error {
	if a.Username == "" {
		return errors.New("auth: username is required")
	}

	if a.Password == "" {
		return errors.New("auth: password is required")
	}

	return nil
}
