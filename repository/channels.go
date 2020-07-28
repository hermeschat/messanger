package repository

import (
	"time"
)

type Channel struct {
	ID        string            `bson:"_id" json:"channel_id"`
	Creator   string            `bson:"creator" json:"creator"`
	Type      string            `bson:"type" json:"type"`
	Members   []string          `bson:"members" json:"members"`
	Roles     map[string]string `bson:"roles" json:"roles"`
	Messages  []string          `bson:"messages" json:"messages"`
	CreatedAt time.Time         `bson:"created_at" json:"created_at"`
}

type ChannelRepository interface {
	GetChannelMessages(id string) ([]*Message, error)
	GetChannelsByCreator(id string) ([]*Channel, error)
	GetChannelsByMember(id string) ([]*Channel, error)
	GetChannelByMembers(ids []string) ([]*Channel, error)
	GetDirectChannelByMembers(ids []string) (*Channel, error)
}
