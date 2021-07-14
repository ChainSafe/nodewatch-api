// Package mongo represent store driver for mongodb
package mongo

import (
	"context"
	"eth2-crawler/config"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

type store struct {
	client  *mongo.Client
	coll    *mongo.Collection
	timeout time.Duration
}

// New creates new instance of Entry Store based on MongoDB
func New(cfg *config.Database) (*store, error) {
	timeout := time.Duration(cfg.Timeout) * time.Second
	opts := options.Client()

	opts.ApplyURI(cfg.URI)
	client, err := mongo.NewClient(opts)
	if err != nil {
		return nil, fmt.Errorf("connecton error [%s]: %w", opts.GetURI(), err)
	}

	// connect to the mongoDB cluster
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	err = client.Connect(ctx)
	if err != nil {
		return nil, err
	}

	// test the connection
	err = client.Ping(ctx, readpref.Primary())
	if err != nil {
		return nil, err
	}

	return &store{
		client:  client,
		coll:    client.Database(cfg.Database).Collection(cfg.Collection),
		timeout: timeout,
	}, nil
}
