// Package store represents the data Store service
package store

import (
	"context"

	"eth2-crawler/models"

	peer "github.com/libp2p/go-libp2p-peer"
)

// Provider represents store provider interface that can be implemented by different DB engines
type Provider interface {
	Create(ctx context.Context, peer *models.Peer) error
	Update(ctx context.Context, peer *models.Peer) error
	Upsert(ctx context.Context, peer *models.Peer) error
	View(ctx context.Context, peerID peer.ID) (*models.Peer, error)
	AggregateByAgentName(ctx context.Context) ([]*models.AggregateData, error)
	AggregateByOperatingSystem(ctx context.Context) ([]*models.AggregateData, error)
	AggregateByCountry(ctx context.Context) ([]*models.AggregateData, error)
}
