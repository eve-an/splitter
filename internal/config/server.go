package config

import (
	"errors"
	"time"
)

type Server struct {
	Address     string
	ReadTimeout time.Duration
}

func (c Server) Validate() error {
	if c.Address == "" {
		return errors.New("empty server address")
	}

	if c.ReadTimeout == 0 {
		return errors.New("empty read timeout")
	}

	return nil
}
