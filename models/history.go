package models

import (
	"time"

	"github.com/google/uuid"
)

type History struct {
	ID        uuid.UUID `bson:"_id" json:"id"`
	Time      int64     `json:"time" bson:"time"`
	SyncNodes int       `bson:"sync_nodes" json:"sync_nodes"`
	Eth2Nodes int       `bson:"eth_2_nodes" json:"eth_2_nodes"`
}

func NewHistory(syncNodes int, eth2Nodes int) *History {
	t := time.Now()
	return &History{
		ID:        uuid.New(),
		Time:      t.Unix(),
		SyncNodes: syncNodes,
		Eth2Nodes: eth2Nodes,
	}
}
