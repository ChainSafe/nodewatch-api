package record

import (
	"context"
	"eth2-crawler/models"
)

// Provider represents store provider interface that can be implemented by different DB engines
type Provider interface {
	Create(ctx context.Context, history *models.History) error
	GetHistory(ctx context.Context, request *models.HistoryRequest) ([]*models.HistoryCount, error)
}
