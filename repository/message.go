package repository

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type MessageType uint8

const (
	MessageTypeText = iota + 1
)

type Message struct {
	ID          string             `bson:"_id" json:"message_id"`
	From        primitive.ObjectID `bson:"from" json:"from"`
	To          primitive.ObjectID `bson:"to" json:"to"`
	Time        time.Time          `bson:"time" json:"time"`
	Channel     primitive.ObjectID `bson:"channel_id" json:"channel_id"`
	MessageType MessageType        `bson:"message_type" json:"message_type"`
	Body        string             `bson:"body" json:"body"`
	Read        bool               `bson:"read" json:"read"`
}

type MessageRepository interface {
	NewMessage(*Message) error
	GetMessages(query map[string]interface{}) ([]*Message, error)
	GetMessage(query map[string]interface{}) (*Message, error)
}
