package repository

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type Channel struct {
	ID        primitive.ObjectID   `bson:"_id" json:"channel_id"`
	Creator   primitive.ObjectID   `bson:"creator" json:"creator"`
	Type      string               `bson:"type" json:"type"`
	Members   []primitive.ObjectID `bson:"members" json:"members"`
	Roles     map[string]string    `bson:"roles" json:"roles"`
	Messages  []primitive.ObjectID `bson:"messages" json:"messages"`
	CreatedAt time.Time            `bson:"created_at" json:"created_at"`
}

type ChannelRepository interface {
	GetChannelMessages(id primitive.ObjectID) ([]*Message, error)
	GetChannelsByCreator(id primitive.ObjectID) ([]*Channel, error)
	GetChannelsByMember(id primitive.ObjectID) ([]*Channel, error)
	GetChannelByMembers(ids []primitive.ObjectID) ([]*Channel, error)
	GetDirectChannelByMembers(ids []primitive.ObjectID) (*Channel, error)
}
