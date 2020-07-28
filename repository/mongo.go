package repository

import "go.mongodb.org/mongo-driver/mongo"

type Mongo struct {
	userColl *mongo.Collection
	msgColl  *mongo.Collection
	chanColl *mongo.Collection
}

func (m *Mongo) GetChannelMessages(id string) ([]*Message, error) {
	panic("implement me")
}

func (m *Mongo) GetChannelsByCreator(id string) ([]*Channel, error) {
	panic("implement me")
}

func (m *Mongo) GetChannelsByMember(id string) ([]*Channel, error) {
	panic("implement me")
}

func (m *Mongo) GetChannelByMembers(ids []string) ([]*Channel, error) {
	panic("implement me")
}

func (m *Mongo) GetDirectChannelByMembers(ids []string) (*Channel, error) {
	panic("implement me")
}

func (m *Mongo) NewMessage(message *Message) error {
	panic("implement me")
}

func (m *Mongo) GetMessages(query map[string]interface{}) ([]*Message, error) {
	panic("implement me")
}

func (m *Mongo) GetMessage(query map[string]interface{}) (*Message, error) {
	panic("implement me")
}

func (m *Mongo) NewUser(user *User) error {
	panic("implement me")
}

func (m *Mongo) GetUser(id string) (*User, error) {
	panic("implement me")
}

func NewMongo(database *mongo.Database) *Mongo {
	return &Mongo{
		userColl: database.Collection("users"),
		msgColl:  database.Collection("messages"),
		chanColl: database.Collection("channels"),
	}
}
