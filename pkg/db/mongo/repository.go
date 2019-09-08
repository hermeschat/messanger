package mongo

import (
	"context"
	"github.com/amirrezaask/config"
	"time"

	"github.com/mitchellh/mapstructure"
	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

//collection gets collection that you gave us name of
func collection(collectionName string) (*mongo.Collection, error) {
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(config.Get("mongo_uri")))
	if err != nil {
		return nil, errors.Wrap(err, "can't connect to mongodb FUCK")
	}
	coll := client.Database(config.Get("database_name")).Collection(collectionName)
	return coll, nil
}

type repository struct {
	name string
}

func (r *repository) Name() string {
	return r.name
}

func (r *repository) Find(id string) (map[string]interface{}, error) {
	s, err := FindOneById("messages", id)
	if err != nil {
		return nil, errors.Wrap(err, "can't find message with given id")
	}
	message := map[string]interface{}{}
	err = mapstructure.Decode(s, &message)
	if err != nil {
		return nil, errors.Wrap(err, "can't construct message from given map from mongo")
	}
	return message, nil
}

func (r *repository) Get(query map[string]interface{}) ([]map[string]interface{}, error) {
	s, err := FindAll("messages", query)
	if err != nil {
		return nil, errors.Wrap(err, "can't find message with given query")
	}
	//messages := []*Message{}
	messages := []map[string]interface{}{}
	for s.Next(context.Background()) {
		msg := &map[string]interface{}{}
		err = s.Decode(msg)
		if err != nil {
			return nil, errors.Wrap(err, "error in decoding")
		}
		messages = append(messages, *msg)
	}

	return messages, nil
}

func (r *repository) Update(selector map[string]interface{}, update map[string]interface{}) (int, error) {
	panic("implement me")
}

func (r *repository) Add(object map[string]interface{}) (string, error) {

	err := InsertOne("messages", object)
	if err != nil {
		return "", errors.Wrap(err, "can't add this message to mongodb")
	}
	//TODO: return objectid of inserted document
	return "", nil
}

func (r *repository) Delete(query map[string]interface{}) error {
	err := DeleteAllMatched(r.Name(), query)
	if err != nil {
		return errors.Wrap(err, "can't delete this message from mongodb")
	}
	return nil
}
