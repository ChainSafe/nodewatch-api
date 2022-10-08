package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"eth2-crawler/graph/generated"
	"eth2-crawler/graph/model"
	svcModels "eth2-crawler/models"
)

// AggregateByAgentName is the resolver for the aggregateByAgentName field.
func (r *queryResolver) AggregateByAgentName(ctx context.Context, peerFilter *model.PeerFilter) ([]*model.AggregateData, error) {
	aggregateData, err := r.peerStore.AggregateByAgentName(ctx, peerFilter)
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

// AggregateByCountry is the resolver for the aggregateByCountry field.
func (r *queryResolver) AggregateByCountry(ctx context.Context, peerFilter *model.PeerFilter) ([]*model.AggregateData, error) {
	aggregateData, err := r.peerStore.AggregateByCountry(ctx, peerFilter)
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

// AggregateByOperatingSystem is the resolver for the aggregateByOperatingSystem field.
func (r *queryResolver) AggregateByOperatingSystem(ctx context.Context, peerFilter *model.PeerFilter) ([]*model.AggregateData, error) {
	aggregateData, err := r.peerStore.AggregateByOperatingSystem(ctx, peerFilter)
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

// AggregateByNetwork is the resolver for the aggregateByNetwork field.
func (r *queryResolver) AggregateByNetwork(ctx context.Context, peerFilter *model.PeerFilter) ([]*model.AggregateData, error) {
	aggregateData, err := r.peerStore.AggregateByNetworkType(ctx, peerFilter)
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

// AggregateByHardforkSchedule is the resolver for the aggregateByHardforkSchedule field.
func (r *queryResolver) AggregateByHardforkSchedule(ctx context.Context, peerFilter *model.PeerFilter) ([]*model.NextHardforkAggregation, error) {
	allPeers, err := r.peerStore.ViewAll(ctx, peerFilter)
	if err != nil {
		return nil, err
	}

	result := []*model.NextHardforkAggregation{}

	for _, group := range model.GroupByHardforkSchedule(allPeers) {
		result = append(result, &model.NextHardforkAggregation{
			Version: group.Version,
			Epoch:   group.Epoch,
			Count:   group.Count,
		})
	}
	return model.SortByCount(result), nil
}

// AggregateByClientVersion is the resolver for the aggregateByClientVersion field.
func (r *queryResolver) AggregateByClientVersion(ctx context.Context, peerFilter *model.PeerFilter) ([]*model.ClientVersionAggregation, error) {
	aggregateData, err := r.peerStore.AggregateByClientVersion(ctx, peerFilter)
	if err != nil {
		return nil, err
	}

	result := []*model.ClientVersionAggregation{}
	for i := range aggregateData {
		versions := []*model.AggregateData{}
		for j := range aggregateData[i].Versions {
			versions = append(versions, &model.AggregateData{
				Name:  aggregateData[i].Versions[j].Name,
				Count: aggregateData[i].Versions[j].Count,
			})
		}
		result = append(result, &model.ClientVersionAggregation{
			Client:   aggregateData[i].Client,
			Count:    aggregateData[i].Count,
			Versions: versions,
		})
	}
	return result, nil
}

// GetHeatmapData is the resolver for the getHeatmapData field.
func (r *queryResolver) GetHeatmapData(ctx context.Context, peerFilter *model.PeerFilter) ([]*model.HeatmapData, error) {
	peers, err := r.peerStore.ViewAll(ctx, peerFilter)
	if err != nil {
		return nil, err
	}

	result := []*model.HeatmapData{}
	for i := range peers {
		if peers[i].GeoLocation != nil &&
			(peers[i].GeoLocation.Latitude != 0 ||
				peers[i].GeoLocation.Longitude != 0) {
			var syncStatus string
			if peers[i].Sync != nil {
				syncStatus = peers[i].Sync.String()
			}
			result = append(result, &model.HeatmapData{
				NetworkType: string(peers[i].GeoLocation.ASN.Type),
				ClientType:  string(peers[i].UserAgent.Name),
				SyncStatus:  syncStatus,
				Latitude:    peers[i].GeoLocation.Latitude,
				Longitude:   peers[i].GeoLocation.Longitude,
				City:        peers[i].GeoLocation.City,
				Country:     peers[i].GeoLocation.Country,
			})
		}
	}
	return result, nil
}

// GetNodeStats is the resolver for the getNodeStats field.
func (r *queryResolver) GetNodeStats(ctx context.Context, peerFilter *model.PeerFilter) (*model.NodeStats, error) {
	aggregateData, err := r.peerStore.AggregateBySyncStatus(ctx, peerFilter)
	if err != nil {
		return nil, err
	}
	return &model.NodeStats{
		TotalNodes:             aggregateData.Total,
		NodeSyncedPercentage:   (float64(aggregateData.Synced) / float64(aggregateData.Total)) * 100,
		NodeUnsyncedPercentage: (float64(aggregateData.Unsynced) / float64(aggregateData.Total)) * 100,
	}, nil
}

// GetNodeStatsOverTime is the resolver for the getNodeStatsOverTime field.
func (r *queryResolver) GetNodeStatsOverTime(ctx context.Context, start float64, end float64, peerFilter *model.PeerFilter) ([]*model.NodeStatsOverTime, error) {
	data, err := r.historyStore.GetHistory(ctx, int64(start), int64(end), peerFilter)
	if err != nil {
		return nil, err
	}
	result := make([]*model.NodeStatsOverTime, 0)
	for _, v := range data {
		result = append(result, &model.NodeStatsOverTime{
			Time:          float64(v.Time),
			TotalNodes:    v.TotalNodes,
			SyncedNodes:   v.SyncedNodes,
			UnsyncedNodes: v.TotalNodes - v.SyncedNodes,
		})
	}
	return result, nil
}

// GetRegionalStats is the resolver for the getRegionalStats field.
func (r *queryResolver) GetRegionalStats(ctx context.Context, peerFilter *model.PeerFilter) (*model.RegionalStats, error) {
	countryAggrData, err := r.peerStore.AggregateByCountry(ctx, peerFilter)
	if err != nil {
		return nil, err
	}

	networkAggrData, err := r.peerStore.AggregateByNetworkType(ctx, peerFilter)
	if err != nil {
		return nil, err
	}

	var hostedCount, nonhostedCount, total int
	for i := range networkAggrData {
		total += networkAggrData[i].Count
		if networkAggrData[i].Name == string(svcModels.UsageTypeHosting) {
			hostedCount += networkAggrData[i].Count
		} else {
			nonhostedCount += networkAggrData[i].Count
		}
	}

	result := &model.RegionalStats{
		TotalParticipatingCountries: len(countryAggrData),
		HostedNodePercentage:        (float64(hostedCount) / float64(total)) * 100,
		NonhostedNodePercentage:     (float64(nonhostedCount) / float64(total)) * 100,
	}
	return result, nil
}

// GetAltairUpgradePercentage is the resolver for the getAltairUpgradePercentage field.
func (r *queryResolver) GetAltairUpgradePercentage(ctx context.Context, peerFilter *model.PeerFilter) (float64, error) {
	aggregateData, err := r.peerStore.AggregateByClientVersion(ctx, peerFilter)
	if err != nil {
		return 0, err
	}
	// check altair upgrade
	count := 0
	total := 0
	for _, client := range aggregateData {
		for _, v := range client.Versions {
			total += v.Count
			if model.SupportAltairUpgrade(client.Client, v.Name) {
				count += v.Count
			}
		}
	}
	percentage := float64(count) / float64(total) * 100
	return percentage, nil
}

// Query returns generated.QueryResolver implementation.
func (r *Resolver) Query() generated.QueryResolver { return &queryResolver{r} }

type queryResolver struct{ *Resolver }
