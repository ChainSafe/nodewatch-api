// Package mongo represent store driver for mongodb
package mongo

import (
	"context"
	"errors"
	"fmt"
	"time"

	"eth2-crawler/models"
	"eth2-crawler/utils/config"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

type store struct {
	client  *mongo.Client
	coll    *mongo.Collection
	timeout time.Duration
}

func (s *store) Upsert(ctx context.Context, peer *models.Peer) error {
	_, err := s.View(ctx, peer.PeerID)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return s.create(ctx, peer)
		}
		return err
	}

	return s.update(ctx, peer)
}

func (s *store) create(ctx context.Context, peer *models.Peer) error {
	_, err := s.coll.InsertOne(ctx, peer, options.InsertOne())
	return err
}

func (s *store) update(ctx context.Context, peer *models.Peer) error {
	filter := bson.D{
		{Key: "id", Value: peer.ID},
	}
	_, err := s.coll.UpdateOne(ctx, filter, bson.D{{Key: "$set", Value: peer}})
	if err != nil {
		return err
	}
	return nil
}

func (s *store) View(ctx context.Context, peerID string) (*models.Peer, error) {
	filter := bson.D{
		{Key: "peer_id", Value: peerID},
	}
	res := new(models.Peer)
	err := s.coll.FindOne(ctx, filter).Decode(res)
	if err != nil {
		return nil, err
	}

	return res, nil
}

type aggregateData struct {
	ID    string `json:"_id" bson:"_id"`
	Count int    `json:"count" bson:"count"`
}

func (s *store) AggregateByAgentName(ctx context.Context) ([]*models.AggregateData, error) {
	query := mongo.Pipeline{
		bson.D{
			{Key: "$group", Value: bson.D{
				{Key: "_id", Value: "$useragent.name"},
				{Key: "count", Value: bson.D{{Key: "$sum", Value: 1}}},
			}},
		},
	}

	cursor, err := s.coll.Aggregate(ctx, query)
	if err != nil {
		return nil, err
	}

	result := []*models.AggregateData{}
	for cursor.Next(ctx) {
		// create a value into which the single document can be decoded
		data := new(aggregateData)
		err := cursor.Decode(data)
		if err != nil {
			return nil, err
		}

		result = append(result, &models.AggregateData{Name: data.ID, Count: data.Count})
	}
	return result, nil
}

func (s *store) AggregateByOperatingSystem(ctx context.Context) ([]*models.AggregateData, error) {
	query := mongo.Pipeline{
		bson.D{
			{Key: "$group", Value: bson.D{
				{Key: "_id", Value: "$useragent.os"},
				{Key: "count", Value: bson.D{{Key: "$sum", Value: 1}}},
			}},
		},
	}
	cursor, err := s.coll.Aggregate(ctx, query)
	if err != nil {
		return nil, err
	}

	result := []*models.AggregateData{}
	for cursor.Next(ctx) {
		// create a value into which the single document can be decoded
		data := new(aggregateData)
		err := cursor.Decode(data)
		if err != nil {
			return nil, err
		}

		result = append(result, &models.AggregateData{Name: data.ID, Count: data.Count})
	}
	return result, nil
}

func (s *store) AggregateByCountry(ctx context.Context) ([]*models.AggregateData, error) {
	query := mongo.Pipeline{
		bson.D{
			{Key: "$group", Value: bson.D{
				{Key: "_id", Value: "$geolocation.country"},
				{Key: "count", Value: bson.D{{Key: "$sum", Value: 1}}},
			}},
		},
	}
	cursor, err := s.coll.Aggregate(ctx, query)
	if err != nil {
		return nil, err
	}

	result := []*models.AggregateData{}
	for cursor.Next(ctx) {
		// create a value into which the single document can be decoded
		data := new(aggregateData)
		err := cursor.Decode(data)
		if err != nil {
			return nil, err
		}

		result = append(result, &models.AggregateData{Name: data.ID, Count: data.Count})
	}
	return result, nil
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
