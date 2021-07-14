// Package store represents the data Store service
package store

import (
	"eth2-crawler/store/mongo"
	"eth2-crawler/utils/config"
)

// Provider represents store provider interface that can be implemented by different DB engines
type Provider interface {
}

// New creates new instance of Entry Store Provider based on provided config
func New(cfg *config.Database) (Provider, error) {
	return mongo.New(cfg)
}
