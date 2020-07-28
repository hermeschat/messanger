package mongo

import (
	"context"
	"github.com/hermeschat/engine/config"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"
)

func New(ctx context.Context) (*mongo.Database, error) {
	cli, err := mongo.NewClient(options.Client().ApplyURI(config.MongoURI()), options.Client().SetConnectTimeout(time.Second*1))
	if err != nil {
	    return nil, err
	}
	err = cli.Ping(ctx, nil)
	if err != nil {
		return nil, err
	}

	databaseName := config.C.GetString("database.name")
	if err != nil {
		return nil, err
	}
	return cli.Database(databaseName), nil
}
