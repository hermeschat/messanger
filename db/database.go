package db

import (
	"context"
	"fmt"
	"time"

	"github.com/amirrezaask/config"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type hermesDatabase struct {
	messages *mongo.Collection
	channels *mongo.Collection
}

//collection gets collection that you gave us name of
func collection(collectionName string) *mongo.Collection {
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(config.Get("mongo_uri")))
	if err != nil {
		panic(fmt.Errorf("cannot get %s collection from mongo due to :%v", err))
	}
	coll := client.Database(config.Get("database_name")).Collection(collectionName)
	return coll
}

func initMongoCollections(h *hermesDatabase) {
	h.channels = collection("channels")
}

var hermesDatabaseInstance *hermesDatabase

//Gets a Instance instance
func Channels() *mongo.Collection {
	return hermesDatabaseInstance.channels
}

func Init() {
	hermesDatabaseInstance = new(hermesDatabase)
	initMongoCollections(hermesDatabaseInstance)
}
