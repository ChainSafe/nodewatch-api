package graph

import (
	svcModels "eth2-crawler/models"
	"github.com/hashicorp/go-version"
)

func supportAltairUpgrade(clientName, ver string) bool {
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

func isMergeReady(clientName, ver string) bool {
	if len(ver) != 0 && ver[0:1] != "v" {
		ver = "v" + ver
	}
	clientVersion, err := version.NewVersion(ver)
	if err != nil {
		return false
	}

	switch svcModels.ClientName(clientName) {
	case svcModels.PrysmClient:
		v, _ := version.NewVersion("v3.0.0")
		if clientVersion.GreaterThanOrEqual(v) {
			return true
		}
	case svcModels.TekuClient:
		v, _ := version.NewVersion("TODO")
		if clientVersion.GreaterThanOrEqual(v) {
			return true
		}
	case svcModels.LighthouseClient:
		v, _ := version.NewVersion("TODO")
		if clientVersion.GreaterThanOrEqual(v) {
			return true
		}
	case svcModels.NimbusClient:
		v, _ := version.NewVersion("TODO")
		if clientVersion.GreaterThanOrEqual(v) {
			return true
		}
	case svcModels.LodestarClient:
		v, _ := version.NewVersion("v1.0.0")
		if clientVersion.GreaterThanOrEqual(v) {
			return true
		}
	default:
		return false
	}
	return false
}
