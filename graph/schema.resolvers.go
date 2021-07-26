// Copyright 2021 ChainSafe Systems
// SPDX-License-Identifier: LGPL-3.0-only

package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"eth2-crawler/graph/generated"
	"eth2-crawler/graph/model"
	"eth2-crawler/store"
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

// Query returns generated.QueryResolver implementation.
func (r *Resolver) Query() generated.QueryResolver {
	return &queryResolver{
		Resolver:  r,
		peerStore: r.peerStore,
	}
}

type queryResolver struct {
	*Resolver
	peerStore store.Provider
}
