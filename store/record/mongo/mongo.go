// Copyright 2021 ChainSafe Systems
// SPDX-License-Identifier: LGPL-3.0-only

// Package mongo implements all the store methods
package mongo

import (
	"context"
	"eth2-crawler/models"
	"eth2-crawler/store/record"
	"eth2-crawler/utils/config"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

type mongoStore struct {
	client  *mongo.Client
	coll    *mongo.Collection
	timeout time.Duration
}

// New creates new instance of History Store based on MongoDB
func New(cfg *config.Database) (record.Provider, error) {
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

	return &mongoStore{
		client:  client,
		coll:    client.Database(cfg.Database).Collection(cfg.HistoryCollection),
		timeout: timeout,
	}, nil
}

func (s mongoStore) Create(ctx context.Context, history *models.History) error {
	_, err := s.coll.InsertOne(ctx, history, options.InsertOne())
	return err
}

func (s mongoStore) GetHistory(ctx context.Context, start int64, end int64) ([]*models.HistoryCount, error) {
	filter := bson.D{{Key: "$and", Value: bson.A{
		bson.D{{Key: "time", Value: bson.D{{Key: "$gt", Value: start}}}},
		bson.D{{Key: "time", Value: bson.D{{Key: "$lt", Value: end}}}},
	}}}
	cursor, err := s.coll.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	var result []*models.History
	for cursor.Next(ctx) {
		// create a value into which the single document can be decoded
		data := new(models.History)
		err := cursor.Decode(data)
		if err != nil {
			return nil, err
		}
		result = append(result, data)
	}

	count := make([]*models.HistoryCount, 0)
	for _, v := range result {
		count = append(count, &models.HistoryCount{
			Time:        v.Time,
			TotalNodes:  v.Eth2Nodes,
			SyncedNodes: v.SyncNodes,
		})
	}
	return count, nil
}
