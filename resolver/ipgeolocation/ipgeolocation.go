// Copyright 2021 ChainSafe Systems
// SPDX-License-Identifier: LGPL-3.0-only

// Package ipgeolocation implements the ip resolver using ipgeolocation APIs
package ipgeolocation

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"

	"eth2-crawler/models"
	"eth2-crawler/resolver"
)

const (
	url = "https://api.ipgeolocation.io/ipgeo"
)

type client struct {
	httpClient *http.Client
	apiKey     string
}

type geoInformation struct {
	ISP          string `json:"isp"`
	Organization string `json:"organization"`
	Country      string `json:"country_name"`
	State        string `json:"state_prov"`
	City         string `json:"city"`
	Latitude     string `json:"latitude"`
	Longitude    string `json:"longitude"`
}

func New(apiKey string, defaultTimeout time.Duration) resolver.Provider {
	return &client{
		httpClient: &http.Client{Timeout: defaultTimeout},
		apiKey:     apiKey,
	}
}

func (c *client) GetGeoLocation(ctx context.Context, ipAddr string) (*models.GeoLocation, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	q := req.URL.Query()
	q.Add("apiKey", c.apiKey)
	q.Add("ip", ipAddr)
	req.URL.RawQuery = q.Encode()

	res, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}

	// nolint
	defer res.Body.Close()
	resBody, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, fmt.Errorf("unable to read response body")
	}

	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("invalid status code returned. statusCode::%d response::%s", res.StatusCode, string(resBody))
	}

	result := &geoInformation{}
	err = json.Unmarshal(resBody, result)
	if err != nil {
		return nil, fmt.Errorf("unable to unmarshal body. error::%w", err)
	}

	lat, _ := strconv.ParseFloat(result.Latitude, 64)
	long, _ := strconv.ParseFloat(result.Longitude, 64)

	geoLoc := &models.GeoLocation{
		Country:   result.Country,
		State:     result.State,
		City:      result.City,
		Latitude:  lat,
		Longitude: long,
	}
	return geoLoc, nil
}
