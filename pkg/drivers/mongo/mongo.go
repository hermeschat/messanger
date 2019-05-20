package mongo

import (
	"context"
	"time"

	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

//GetCollection gets collection that you gave us name of
func GetCollection(collectionName string) (*mongo.Collection, error) {
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	client, err := mongo.Connect(ctx, options.Client().ApplyURI("mongodb://localhost:27017"))
	if err != nil {
		return nil, errors.Wrap(err, "can't connect to mongodb FUCK")
	}
	//TODO use config
	coll := client.Database("dbname").Collection(collectionName)
	return coll, nil
}
