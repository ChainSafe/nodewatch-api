// Package store represents the data Store service
package store

import (
	"context"

	"eth2-crawler/models"
	"eth2-crawler/store/mongo"
	"eth2-crawler/utils/config"
)

// Provider represents store provider interface that can be implemented by different DB engines
type Provider interface {
	Upsert(ctx context.Context, peer *models.Peer) error
	View(ctx context.Context, peerID string) (*models.Peer, error)
	AggregateByAgentName(ctx context.Context) ([]*models.AggregateData, error)
	AggregateByOperatingSystem(ctx context.Context) ([]*models.AggregateData, error)
	AggregateByCountry(ctx context.Context) ([]*models.AggregateData, error)
}

// New creates new instance of Entry Store Provider based on provided config
func New(cfg *config.Database) (Provider, error) {
	return mongo.New(cfg)
}
