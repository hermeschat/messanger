package mongo


import "go.mongodb.org/mongo-driver/mongo"

client, err := mongo.NewClient(options.Client().ApplyURI("mongodb://localhost:27017"))
