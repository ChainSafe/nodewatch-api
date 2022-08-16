// Copyright 2022 ChainSafe Systems
// SPDX-License-Identifier: LGPL-3.0-only

package model

import (
	svcModels "eth2-crawler/models"
	"fmt"

	"github.com/hashicorp/go-version"
)

func SortByCount(data []*NextHardforkAggregation) []*NextHardforkAggregation {
	for i := 0; i < len(data); i++ {
		for j := i + 1; j < len(data); j++ {
			if data[i].Count < data[j].Count {
				data[i], data[j] = data[j], data[i]
			}
		}
	}
	return data
}

func GroupByHardforkSchedule(allPeers []*svcModels.Peer) map[string]*NextHardforkAggregation {
	result := map[string]*NextHardforkAggregation{}
	for _, peer := range allPeers {
		key := fmt.Sprintf("%s-%s", peer.NextForkVersion.String(), peer.NextForkEpoch.String())
		if _, ok := result[key]; !ok {
			result[key] = &NextHardforkAggregation{
				Epoch:   peer.NextForkEpoch.String(),
				Version: peer.NextForkVersion.String(),
				Count:   1,
			}
		} else {
			result[key].Count++
		}
	}
	return result
}

func SupportAltairUpgrade(clientName, ver string) bool {
	if len(ver) != 0 && ver[0:1] != "v" {
		ver = "v" + ver
	}
	clientVersion, err := version.NewVersion(ver)
	if err != nil {
		return false
	}

	switch svcModels.ClientName(clientName) {
	case svcModels.PrysmClient:
		v, _ := version.NewVersion("v2.0.0")
		if clientVersion.GreaterThanOrEqual(v) {
			return true
		}
	case svcModels.TekuClient:
		v, _ := version.NewVersion("v21.9.2")
		if clientVersion.GreaterThanOrEqual(v) {
			return true
		}
	case svcModels.LighthouseClient:
		v, _ := version.NewVersion("v2.0.0")
		if clientVersion.GreaterThanOrEqual(v) {
			return true
		}
	case svcModels.NimbusClient:
		v, _ := version.NewVersion("v1.5.0")
		if clientVersion.GreaterThanOrEqual(v) {
			return true
		}
	case svcModels.LodestarClient:
		v, _ := version.NewVersion("v0.31.0")
		if clientVersion.GreaterThanOrEqual(v) {
			return true
		}
	default:
		return false
	}
	return false
}
