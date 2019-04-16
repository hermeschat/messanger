package message

import "time"

type Message struct {
	MessageID   string `bson:"_id" json:"_id"`
	From        string
	Time        time.Time
	ChannelID   string
	MessageType string
	Body        string
}
