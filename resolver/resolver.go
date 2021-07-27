// Copyright 2021 ChainSafe Systems
// SPDX-License-Identifier: LGPL-3.0-only

// Package resolver implements the ip resolver
package resolver

import (
	"context"

	"eth2-crawler/models"
)

type Provider interface {
	GetGeoLocation(ctx context.Context, ipAddr string) (*models.GeoLocation, error)
}
