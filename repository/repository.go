package repository

import "time"


type Repository interface {
	Channels
	Messages
	Users
}



type MessageType uint8

const (
	MessageTypeText = iota + 1
)

type Message struct {
	ID          string      `bson:"_id" json:"message_id"`
	From        string      `bson:"from" json:"from"`
	To          string      `bson:"to" json:"to"`
	Time        time.Time   `bson:"time" json:"time"`
	Channel     string      `bson:"channel_id" json:"channel_id"`
	MessageType MessageType `bson:"message_type" json:"message_type"`
	Body        string      `bson:"body" json:"body"`
	Read        bool        `bson:"read" json:"read"`
}

type Messages interface {
	NewMessage(*Message) error
	GetMessages(query map[string]interface{}) ([]*Message, error)
	GetMessage(query map[string]interface{}) (*Message, error)
}


type Channel struct {
	ID        string            `bson:"_id" json:"channel_id"`
	Creator   string            `bson:"creator" json:"creator"`
	Type      string            `bson:"type" json:"type"`
	Members   []string          `bson:"members" json:"members"`
	Roles     map[string]string `bson:"roles" json:"roles"`
	Messages  []string          `bson:"messages" json:"messages"`
	CreatedAt time.Time         `bson:"created_at" json:"created_at"`
}

type Channels interface {
	GetChannelMessages(id string) ([]*Message, error)
	GetChannelsByCreator(id string) ([]*Channel, error)
	GetChannelsByMember(id string) ([]*Channel, error)
	GetChannelByMembers(ids []string) ([]*Channel, error)
	GetDirectChannelByMembers(ids []string) (*Channel, error)
}

type User struct {
	ID string
	Avatar string
	Handle string
}

type Users interface {
	NewUser(user *User) error
	GetUser(id string) (*User, error)
}
