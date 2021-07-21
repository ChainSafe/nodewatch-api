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
	// Sending dummy data
	// return []*model.AggregateData{
	// 	&model.AggregateData{Name: "Cortex", Count: 356},
	// 	&model.AggregateData{Name: "Lighthouse", Count: 250},
	// 	&model.AggregateData{Name: "Lodestar", Count: 211},
	// 	&model.AggregateData{Name: "Nimbus", Count: 189},
	// 	&model.AggregateData{Name: "Prysm", Count: 115},
	// 	&model.AggregateData{Name: "Teku", Count: 112},
	// 	&model.AggregateData{Name: "Trinity", Count: 89},
	// }, nil

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
	// Sending dummy data
	// return []*model.AggregateData{
	// 	&model.AggregateData{Name: "United States", Count: 558},
	// 	&model.AggregateData{Name: "Germany", Count: 429},
	// 	&model.AggregateData{Name: "China", Count: 412},
	// 	&model.AggregateData{Name: "France", Count: 378},
	// 	&model.AggregateData{Name: "Singapore", Count: 349},
	// 	&model.AggregateData{Name: "United Kingdom", Count: 187},
	// 	&model.AggregateData{Name: "Canada", Count: 173},
	// 	&model.AggregateData{Name: "Netherlands", Count: 113},
	// 	&model.AggregateData{Name: "Japan", Count: 104},
	// 	&model.AggregateData{Name: "Finland", Count: 23},
	// 	&model.AggregateData{Name: "South Korea", Count: 12},
	// }, nil

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
	// Sending dummy data
	// return []*model.AggregateData{
	// 	&model.AggregateData{Name: "Linux", Count: 1023},
	// 	&model.AggregateData{Name: "Windows", Count: 294},
	// 	&model.AggregateData{Name: "MacOS", Count: 138},
	// }, nil

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
	// Sending dummy data
	return []*model.AggregateData{
		&model.AggregateData{Name: "Hosted", Count: 1002},
		&model.AggregateData{Name: "Residential", Count: 445},
		&model.AggregateData{Name: "Business", Count: 76},
	}, nil
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
