package repository

import "go.mongodb.org/mongo-driver/mongo"

type Mongo struct{
	coll *mongo.Collection
}
func NewMongo(connectionURI string) (*Mongo, error) {

}
