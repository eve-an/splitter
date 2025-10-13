package config

import (
	"errors"
	"os"
	"time"

	"github.com/eve-an/splitter/internal/http"
)

type Config struct {
	LogLevel     string
	ServerConifg http.Config
}

func (c Config) Validate() error {
	var errs []error

	if c.LogLevel == "" {
		errs = append(errs, errors.New("empty log level"))
	}

	errs = append(errs, c.ServerConifg.Validate())

	return errors.Join(errs...)
}

func Load() (c Config, err error) {
	c.LogLevel = "debug"
	c.ServerConifg.Address = ":8080"
	c.ServerConifg.ReadTimeout = time.Second * 5

	if addr := os.Getenv("SPLITTER_ADDR"); addr != "" {
		c.ServerConifg.Address = addr
	}

	return c, c.Validate()
}
