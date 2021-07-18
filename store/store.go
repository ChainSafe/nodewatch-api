// Package store represents the data Store service
package store

import (
	"context"
	"eth2-crawler/crawler/peer"
	"eth2-crawler/store/mongo"
	"eth2-crawler/utils/config"
)

// Provider represents store provider interface that can be implemented by different DB engines
type Provider interface {
	Create(ctx context.Context, peer *peer.Peer) error
}

// New creates new instance of Entry Store Provider based on provided config
func New(cfg *config.Database) (Provider, error) {
	return mongo.New(cfg)
}
