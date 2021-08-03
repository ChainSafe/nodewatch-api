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
