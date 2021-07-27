// Copyright 2021 ChainSafe Systems
// SPDX-License-Identifier: LGPL-3.0-only

package models

// AggregateData represents data of group by queries
type AggregateData struct {
	Name  string `json:"name"`
	Count int    `json:"count"`
}

// HeatmapData represents data for heatmap
type HeatmapData struct {
	NetworkType string  `json:"network_type"`
	ClientType  string  `json:"client_type"`
	SyncStatus  string  `json:"sync_status"`
	Latitude    float64 `json:"latitude"`
	Longitude   float64 `json:"longitude"`
}
