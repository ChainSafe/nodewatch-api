// Copyright 2021 ChainSafe Systems
// SPDX-License-Identifier: LGPL-3.0-only

package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"eth2-crawler/graph/generated"
	"eth2-crawler/graph/model"
)

func (r *queryResolver) AggregateByAgentName(ctx context.Context) ([]*model.AggregateData, error) {
	aggregateData, err := r.peerStore.AggregateByAgentName(ctx)
	if err != nil {
		return nil, err
	}

	result := []*model.AggregateData{}
	for i := range aggregateData {
		result = append(result, &model.AggregateData{
			Name:  aggregateData[i].Name,
			Count: aggregateData[i].Count,
		})
	}
	return result, nil
}

func (r *queryResolver) AggregateByCountry(ctx context.Context) ([]*model.AggregateData, error) {
	aggregateData, err := r.peerStore.AggregateByCountry(ctx)
	if err != nil {
		return nil, err
	}

	result := []*model.AggregateData{}
	for i := range aggregateData {
		result = append(result, &model.AggregateData{
			Name:  aggregateData[i].Name,
			Count: aggregateData[i].Count,
		})
	}
	return result, nil
}

func (r *queryResolver) AggregateByOperatingSystem(ctx context.Context) ([]*model.AggregateData, error) {
	aggregateData, err := r.peerStore.AggregateByOperatingSystem(ctx)
	if err != nil {
		return nil, err
	}

	result := []*model.AggregateData{}
	for i := range aggregateData {
		result = append(result, &model.AggregateData{
			Name:  aggregateData[i].Name,
			Count: aggregateData[i].Count,
		})
	}
	return result, nil
}

func (r *queryResolver) AggregateByNetwork(ctx context.Context) ([]*model.AggregateData, error) {
	aggregateData, err := r.peerStore.AggregateByNetworkType(ctx)
	if err != nil {
		return nil, err
	}

	result := []*model.AggregateData{}
	for i := range aggregateData {
		result = append(result, &model.AggregateData{
			Name:  aggregateData[i].Name,
			Count: aggregateData[i].Count,
		})
	}
	return result, nil
}

func (r *queryResolver) GetHeatmapData(ctx context.Context) ([]*model.HeatmapData, error) {
	peers, err := r.peerStore.ViewAll(ctx)
	if err != nil {
		return nil, err
	}

	result := []*model.HeatmapData{}
	for i := range peers {
		if peers[i].GeoLocation != nil &&
			(peers[i].GeoLocation.Latitude != 0 ||
				peers[i].GeoLocation.Longitude != 0) {
			result = append(result, &model.HeatmapData{
				NetworkType: string(peers[i].GeoLocation.ASN.Type),
				ClientType:  string(peers[i].UserAgent.Name),
				Latitude:    peers[i].GeoLocation.Latitude,
				Longitude:   peers[i].GeoLocation.Longitude,
			})
		}
	}
	return result, nil
}

// Query returns generated.QueryResolver implementation.
func (r *Resolver) Query() generated.QueryResolver { return &queryResolver{r} }

type queryResolver struct{ *Resolver }
