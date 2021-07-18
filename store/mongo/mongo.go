// Package mongo represent store driver for mongodb
package mongo

import (
	"context"
	"errors"
	"eth2-crawler/crawler/peer"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson"

	"eth2-crawler/utils/config"

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

// Create creates new peer entry
func (s *store) Create(ctx context.Context, peer *peer.Peer) error {
	var cancel context.CancelFunc
	ctx, cancel = context.WithTimeout(ctx, s.timeout)
	defer cancel()
	// check if already exists
	_, err := s.View(ctx, string(peer.ID))
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			_, err = s.coll.InsertOne(ctx, peer, options.InsertOne())
			return err
		}
		return err
	}
	return nil
}

// View returns Peer by ID
func (s *store) View(ctx context.Context, id string) (*peer.Peer, error) {
	ctx, cancel := context.WithTimeout(ctx, s.timeout)
	defer cancel()

	filter := bson.D{{Key: "_id", Value: id}}
	res := new(peer.Peer)
	err := s.coll.FindOne(ctx, filter).Decode(res)
	return res, err
}
