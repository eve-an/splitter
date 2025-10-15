package config

import (
	"errors"
	"os"
	"time"

	"github.com/eve-an/splitter/internal/http"
)

type DatabaseConfig struct {
	URL             string
	MaxOpenConns    int
	MaxIdleConns    int
	ConnMaxLifetime time.Duration
}

func (c DatabaseConfig) Validate() error {
	var errs []error
	if c.URL == "" {
		errs = append(errs, errors.New("database url is required"))
	}
	if c.MaxOpenConns < 0 {
		errs = append(errs, errors.New("max open connections cannot be negative"))
	}
	if c.MaxIdleConns < 0 {
		errs = append(errs, errors.New("max idle connections cannot be negative"))
	}
	if c.ConnMaxLifetime < 0 {
		errs = append(errs, errors.New("connection max lifetime cannot be negative"))
	}

	return errors.Join(errs...)
}

type Config struct {
	LogLevel     string
	ServerConifg http.Config
	Database     DatabaseConfig
}

func (c Config) Validate() error {
	var errs []error

	if c.LogLevel == "" {
		errs = append(errs, errors.New("empty log level"))
	}

	errs = append(errs, c.ServerConifg.Validate())
	errs = append(errs, c.Database.Validate())

	return errors.Join(errs...)
}

func Load() (c Config, err error) {
	c.LogLevel = "debug"
	c.ServerConifg.Address = ":8080"
	c.ServerConifg.ReadTimeout = time.Second * 5

	c.Database.MaxOpenConns = 10
	c.Database.MaxIdleConns = 5
	c.Database.ConnMaxLifetime = time.Minute * 5

	if addr := os.Getenv("SPLITTER_ADDR"); addr != "" {
		c.ServerConifg.Address = addr
	}

	if dbURL := os.Getenv("SPLITTER_DB_URL"); dbURL != "" {
		c.Database.URL = dbURL
	}

	return c, c.Validate()
}
