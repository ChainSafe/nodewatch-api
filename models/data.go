// Copyright 2021 ChainSafe Systems
// SPDX-License-Identifier: LGPL-3.0-only

package models

// AggregateData represents data of group by queries
type AggregateData struct {
	Name  string `json:"name"`
	Count int    `json:"count"`
}

// ClientVersionAggregation represents aggregation data for client and client version
type ClientVersionAggregation struct {
	Client   string           `json:"client"`
	Count    int              `json:"count"`
	Versions []*AggregateData `json:"versions"`
}

// HistoryRequest represent the required field to view history
type HistoryRequest struct {
	LastDays int `json:"last_days"`
}

type HistoryCount struct {
	Time        int64 `json:"time"`
	TotalNodes  int   `json:"total_nodes"`
	SyncedNodes int   `json:"synced_nodes"`
}
