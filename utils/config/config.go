// Package config holds the config file parsing functionality
package config

import (
	"fmt"
	"io/ioutil"
	"os"

	"gopkg.in/yaml.v2"
)

// Configuration holds data necessary for configuring application
type Configuration struct {
	Server *Server   `yaml:"server,omitempty"`
	DB     *Database `yaml:"database,omitempty"`
}

// Server holds data necessary for server configuration
type Server struct {
	Port         string `yaml:"port,omitempty"`
	ReadTimeout  int    `yaml:"read_timeout_seconds,omitempty"`
	WriteTimeout int    `yaml:"write_timeout_seconds,omitempty"`
}

// Database is a MongoDB config
type Database struct {
	URI         string `yaml:"-"`
	Timeout     int    `yaml:"request_timeout_sec"`
	Database    string `yaml:"database"`
	Collection  string `yaml:"collection"`
	Concurrency int    `yaml:"concurrency"`
	MaxRetries  int    `yaml:"max_retries"`
}

func loadDatabaseURI() string {
	return os.Getenv("CSF_MONGODB_CONN")
}

// Load returns Configuration struct
func Load(path string) (*Configuration, error) {
	bytes, err := ioutil.ReadFile(path)

	if err != nil {
		return nil, fmt.Errorf("error reading config file, %w", err)
	}

	var cfg = new(Configuration)

	if err := yaml.Unmarshal(bytes, cfg); err != nil {
		return nil, fmt.Errorf("unable to decode into struct, %w", err)
	}

	// load envs
	cfg.DB.URI = loadDatabaseURI()

	return cfg, nil
}
