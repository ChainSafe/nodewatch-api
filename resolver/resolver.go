// Package resolver implements the ip resolver
package resolver

import (
	"context"
	"encoding/json"
	"eth2-crawler/models"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

const (
	url = "https://api.ipgeolocation.io/ipgeo"
)

type Provider interface {
	GetGeoLocation(ctx context.Context, ipAddr string) (*models.GeoLocation, error)
}

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

func New(apiKey string, defaultTimeout time.Duration) Provider {
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
		return nil, fmt.Errorf("unable to unmarshal body. error::%s", err.Error())
	}

	geoLoc := &models.GeoLocation{
		ISP:          result.ISP,
		Organization: result.Organization,
		Country:      result.Country,
		State:        result.State,
		City:         result.City,
		Latitude:     result.Latitude,
		Longitude:    result.Longitude,
	}

	return geoLoc, nil
}
