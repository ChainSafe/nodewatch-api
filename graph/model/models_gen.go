// Code generated by github.com/99designs/gqlgen, DO NOT EDIT.

package model

type AggregateData struct {
	Name  string `json:"name"`
	Count int    `json:"count"`
}

type ClientVersionAggregation struct {
	Client   string           `json:"client"`
	Count    int              `json:"count"`
	Versions []*AggregateData `json:"versions"`
}

type HeatmapData struct {
	NetworkType string  `json:"networkType"`
	ClientType  string  `json:"clientType"`
	SyncStatus  string  `json:"syncStatus"`
	Latitude    float64 `json:"latitude"`
	Longitude   float64 `json:"longitude"`
	City        string  `json:"city"`
	Country     string  `json:"country"`
}

type NextHardforkAggregation struct {
	Version string `json:"version"`
	Epoch   string `json:"epoch"`
	Count   int    `json:"count"`
}

type NodeStats struct {
	TotalNodes             int     `json:"totalNodes"`
	NodeSyncedPercentage   float64 `json:"nodeSyncedPercentage"`
	NodeUnsyncedPercentage float64 `json:"nodeUnsyncedPercentage"`
}

type NodeStatsOverTime struct {
	Time          float64 `json:"time"`
	TotalNodes    int     `json:"totalNodes"`
	SyncedNodes   int     `json:"syncedNodes"`
	UnsyncedNodes int     `json:"unsyncedNodes"`
}

type PeerFilter struct {
	ForkDigest *string `json:"forkDigest"`
}

type RegionalStats struct {
	TotalParticipatingCountries int     `json:"totalParticipatingCountries"`
	HostedNodePercentage        float64 `json:"hostedNodePercentage"`
	NonhostedNodePercentage     float64 `json:"nonhostedNodePercentage"`
}
