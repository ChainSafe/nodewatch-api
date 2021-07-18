// Package config holds the config file parsing functionality
package config

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"

	"gopkg.in/yaml.v2"
)

// Configuration holds data necessary for configuring application
type Configuration struct {
	Server   *Server   `yaml:"server,omitempty"`
	Database *Database `yaml:"database,omitempty"`
	Resolver *Resolver `yaml:"resolver,omitempty"`
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

// Resolver provides config for resolver
type Resolver struct {
	APIKey  string `yaml:"-"`
	Timeout int    `yaml:"request_timeout_sec"`
}

func loadDatabaseURI() (string, error) {
	mongoURI := os.Getenv("MONGODB_URI")
	if mongoURI == "" {
		return "", errors.New("MONGODB_URI is required")
	}
	return mongoURI, nil
}

func loadResolverAPIKey() (string, error) {
	resolverAPIKey := os.Getenv("RESOLVER_API_KEY")
	if resolverAPIKey == "" {
		return "", errors.New("RESOLVER_API_KEY is required")
	}
	return resolverAPIKey, nil
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
	cfg.Database.URI, err = loadDatabaseURI()
	if err != nil {
		return nil, err
	}
	cfg.Resolver.APIKey, err = loadResolverAPIKey()
	if err != nil {
		return nil, err
	}

	return cfg, nil
}
