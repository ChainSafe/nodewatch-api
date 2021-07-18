package models

// AggregateData represents data of group by queries
type AggregateData struct {
	Name  string `json:"name"`
	Count int    `json:"count"`
}
