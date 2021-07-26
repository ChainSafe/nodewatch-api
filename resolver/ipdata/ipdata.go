// Copyright 2021 ChainSafe Systems
// SPDX-License-Identifier: LGPL-3.0-only

// Package ipdata implements the ip resolver using ipdata APIs
package ipdata

import (
	"context"
	"eth2-crawler/models"
	"eth2-crawler/resolver"
	"time"

	ipdata "github.com/ipdata/go"
)

type client struct {
	ipdataClient   ipdata.Client
	defaultTimeout time.Duration
}

func New(apiKey string, defaultTimeout time.Duration) (resolver.Provider, error) {
	ipdataClient, err := ipdata.NewClient(apiKey)
	if err != nil {
		return nil, err
	}

	return &client{
		ipdataClient:   ipdataClient,
		defaultTimeout: defaultTimeout,
	}, nil
}

func (c *client) GetGeoLocation(ctx context.Context, ipAddr string) (*models.GeoLocation, error) {
	ctx, cancel := context.WithTimeout(ctx, c.defaultTimeout)
	defer cancel()

	data, err := c.ipdataClient.LookupWithContext(ctx, ipAddr)
	if err != nil {
		return nil, err
	}

	var asnType models.UsageType
	switch data.ASN.Type {
	case "hosting":
		asnType = models.UsageTypeHosting
	case "isp":
		asnType = models.UsageTypeResidential
	case "business":
		asnType = models.UsageTypeBusiness
	case "edu":
		asnType = models.UsageTypeEducation
	case "gov":
		asnType = models.UsageTypeGovernment
	case "mil":
		asnType = models.UsageTypeMilitary
	default:
		asnType = models.UsageTypeNil
	}

	geoLoc := &models.GeoLocation{
		ASN: models.ASN{
			ID:     data.ASN.ASN,
			Name:   data.ASN.Name,
			Domain: data.ASN.Domain,
			Route:  data.ASN.Route,
			Type:   asnType,
		},
		Country:   data.CountryName,
		State:     data.Region,
		City:      data.City,
		Latitude:  data.Latitude,
		Longitude: data.Longitude,
	}
	return geoLoc, nil
}
