package channel

import (
	"context"
	"log"

	"git.raad.cloud/cloud/hermes/pkg/drivers/mongo"
	"github.com/mitchellh/mapstructure"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	mgo "go.mongodb.org/mongo-driver/mongo"
)

//Type specifies which type of channel we are using
type Type int

const (
	//Secret chat with expire time
	Secret = iota
	//Private chat between two persons
	Private
	//TGChannel just Telegram channel
	TGChannel
	//Group is like telegram groups
	Group
)

//Channel ...
type Channel struct {
	ChannelID string              `json:"channelID" bson:"ChannelID"`
	Members   []string            `json:"Members" bson:"Members"`
	Creator   string              `json:"Creator" bson:"Creator"`
	Type      Type                `json:"type" bson:"type"`
	Roles     map[string][]string `json:"roles" bson:"roles"`
}

//ConstructFromMap ...
func ConstructFromMap(m map[string]interface{}) (*Channel, error) {
	channel := &Channel{}
	err := mapstructure.Decode(m, channel)
	if err != nil {
		return nil, errors.Wrap(err, "can't construct channel from given map")
	}
	return channel, nil
}

func Get(id string) (*Channel, error) {
	s, err := mongo.FindOneById("channels", id)
	if err != nil {
		return nil, errors.Wrap(err, "can't find channel with given id")
	}
	channel := &Channel{}
	err = mapstructure.Decode(s, channel)
	if err != nil {
		return nil, errors.Wrap(err, "can't construct channel from given map from mongo")
	}
	return channel, nil
}

func (s *Channel) ToMap() (map[string]interface{}, error) {
	m := map[string]interface{}{}
	err := mapstructure.Decode(s, &m)
	if err != nil {
		return nil, errors.Wrap(err, "can't create map from this channel")
	}
	return m, nil
}
func GetAll(query map[string]interface{}) ([]*Channel, error) {

	cur, err := mongo.FindAll("channels", query)
	if err == mgo.ErrNoDocuments {
		return nil, mgo.ErrNoDocuments
	}
	if err != nil {
		return nil, errors.Wrap(err, "can't find channel with given query")
	}
	var channels []*Channel
	//err = cur.Decode(channels)
	//if err != nil {
	//	return nil, errors.Wrap(err, "can't construct channel from given map from mongo")
	//}
	for cur.Next(context.Background()) {
		var elem Channel
		err := cur.Decode(&elem)
		if err != nil {
			log.Fatal(err)
		}
		logrus.Info(elem)
		channels = append(channels, &elem)
	}
	return channels, nil
}

//Add adds channel
func Add(channel *Channel) error {

	err := mongo.InsertOne("channels", channel)
	if err != nil {
		return errors.Wrap(err, "can't add this channel to mongodb")
	}
	return nil
}

func AddAll(channels []Channel) error {
	channelsMap := []interface{}{}
	for _, sess := range channels {
		m := map[string]interface{}{}
		err := mapstructure.Decode(sess, m)
		if err != nil {
			return errors.Wrap(err, "can't convert channels to map")
		}
		channelsMap = append(channelsMap, m)
	}
	err := mongo.InsertAll("channels", channelsMap)

	if err != nil {
		return errors.Wrap(err, "can't add this channel to mongodb")
	}
	return nil

}
func Delete(id string) error {
	err := mongo.DeleteById("channels", id)
	if err != nil {
		return errors.Wrap(err, "can't delete this channel from mongodb")
	}
	return nil
}

func DeleteAll(query map[string]interface{}) error {
	err := mongo.DeleteAllMatched("channels", query)
	if err != nil {
		return errors.Wrap(err, "can't delete with given query from mongo")
	}
	return nil
}

//TODO add count of updated docs
func UpdateOne(id string, query map[string]interface{}) error {
	err := mongo.UpdateOne("channels", id, query)
	if err != nil {
		return errors.Wrap(err, "can't update this id with given query")
	}
	return nil
}

func UpdateAll(selector map[string]interface{}, updator map[string]interface{}) error {
	err := mongo.UpdateAllMatched("channels", selector, updator)
	if err != nil {
		return errors.Wrap(err, "can't update channel with given query")
	}
	return nil
}
