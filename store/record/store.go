// Copyright 2021 ChainSafe Systems
// SPDX-License-Identifier: LGPL-3.0-only

// Package record implements db for historic data
package record

import (
	"context"

	"eth2-crawler/graph/model"
	"eth2-crawler/models"
)

// Provider represents store provider interface that can be implemented by different DB engines
type Provider interface {
	Create(ctx context.Context, history *models.History) error
	GetHistory(ctx context.Context, start int64, end int64, peerFilter *model.PeerFilter) ([]*models.HistoryCount, error)
}
