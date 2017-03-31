package config

import (
	"errors"
	"flag"
	"os"
)

type Config struct {
	DBHost string
	DBName string
}

func New(arguments []string) (*Config, error) {
	config := new(Config)
	f := flag.NewFlagSet(os.Args[0], flag.ContinueOnError)

	f.StringVar(&config.DBHost, "db-host", config.DBHost, "The hostname of the database server")
	f.StringVar(&config.DBName, "db-name", config.DBName, "The database name")

	if err := f.Parse(arguments); err != nil {
		return nil, err
	}

	if containDefaultValues(config) {
		f.PrintDefaults()
		return nil, errors.New("Some fields are not set.")
	}

	return config, nil
}

func containDefaultValues(c *Config) bool {
	return c.DBHost == "" || c.DBName == ""
}
