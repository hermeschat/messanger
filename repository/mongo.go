package repository

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type Mongo struct {
	userColl *mongo.Collection
	msgColl  *mongo.Collection
	chanColl *mongo.Collection
}


func getOne(ctx context.Context, coll *mongo.Collection, query bson.M, dest interface{}) error {
	res := coll.FindOne(ctx, query)
	if err := res.Err(); err != nil {
		return err
	}
	err := res.Decode(dest)
	if err != nil {
		return err
	}
	return nil
}

func (m *Mongo) GetChannel(ctx context.Context, id string) (*Channel, error) {
	c := &Channel{}
	err := getOne(ctx, m.chanColl, bson.M{
		"_id": id,
	}, c)
	return c, err
}

func (m *Mongo) NewChannel(ctx context.Context, channel *Channel) error {
	_, err := m.userColl.InsertOne(ctx, channel)
	return err
}

func (m *Mongo) GetChannelMessages(ctx context.Context, id string) ([]*Message, error) {
	panic("implement me")
}

func (m *Mongo) GetChannelsByCreator(ctx context.Context, id string) ([]*Channel, error) {
	panic("implement me")
}

func (m *Mongo) GetChannelsByMember(ctx context.Context, id string) ([]*Channel, error) {
	panic("implement me")
}

func (m *Mongo) GetChannelByMembers(ctx context.Context, ids []string) ([]*Channel, error) {
	panic("implement me")
}

func (m *Mongo) GetDirectChannelByMembers(ctx context.Context, ids []string) (*Channel, error) {
	panic("implement me")
}

func (m *Mongo) NewMessage(ctx context.Context, message *Message) error {
	_, err := m.msgColl.InsertOne(ctx, message)
	return err
}

func (m *Mongo) GetMessages(ctx context.Context, query map[string]interface{}) ([]*Message, error) {
	panic("implement me")
}

func (m *Mongo) GetMessage(ctx context.Context, id string) (*Message, error) {
	msg := &Message{}
	err := getOne(ctx, m.msgColl, bson.M{"_id": id}, msg)
	return msg, err
}

func (m *Mongo) NewUser(ctx context.Context, user *User) error {
	_, err := m.userColl.InsertOne(ctx, user)
	return err
}

func (m *Mongo) GetUser(ctx context.Context, id string) (*User, error) {
	u := &User{}
	err := getOne(ctx, m.userColl, bson.M{"_id": id}, u)
	return u, err
}

func NewMongo(database *mongo.Database) *Mongo {
	return &Mongo{
		userColl: database.Collection("users"),
		msgColl:  database.Collection("messages"),
		chanColl: database.Collection("channels"),
	}
}
