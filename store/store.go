// Copyright 2021 ChainSafe Systems
// SPDX-License-Identifier: LGPL-3.0-only

// Package store represents the data Store service
package store

import (
	"context"

	"eth2-crawler/models"

	"github.com/libp2p/go-libp2p-core/peer"
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
	AggregateByNetworkType(ctx context.Context) ([]*models.AggregateData, error)
}
