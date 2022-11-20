package db

import (
	"context"

	"github.com/akwanmaroso/users-api/config"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

func NewMongoDBConnection(ctx context.Context, cfg *config.Config) (*mongo.Client, error) {
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(cfg.MongoDB.URI))
	if err != nil {
		// log.Fatal(ctx, "failed to connect to mongodb", err)
		return nil, err
	}

	err = client.Ping(ctx, readpref.Primary())
	if err != nil {
		// log.Fatal(ctx, "failed to ping connection", err)
		return nil, err
	}

	return client, nil
}
