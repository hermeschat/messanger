package db

import (
	"github.com/mitchellh/mapstructure"
)

type Message struct {
	MessageId   string `bson:"_id" json:"message_id"`
	From        string `bson:"from" json:"from"`
	To          string `bson:"to" json:"to"`
	Time        string `bson:"time" json:"time"`
	ChannelID   string `bson:"channel_id" json:"channel_id"`
	MessageType string `bson:"message_type" json:"message_type"`
	Body        string `bson:"body" json:"body"`
	Read        bool   `bson:"read" json:"read"`
}

func (m *Message) ToMap() (map[string]interface{}, error) {
	o := map[string]interface{}{}
	err := mapstructure.Decode(m, &o)
	if err != nil {
		return nil, err
	}
	return o, nil
}

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

type Channel struct {
	ChannelId string              `bson:"_id" json:"channel_id"`
	Members   []string            `bson:"members" json:"members"`
	Creator   string              `bson:"creator" json:"creator"`
	Type      int                 `bson:"type" json:"type"`
	Roles     map[string][]string `bson:"roles" json:"roles"`
}

func (c *Channel) ToMap() (map[string]interface{}, error) {
	m := map[string]interface{}{}
	err := mapstructure.Decode(c, &m)
	if err != nil {
		return nil, err
	}
	return m, nil
}
func (c *Channel) FromMap(m map[string]interface{}) error {
	err := mapstructure.Decode(m, c)
	if err != nil {
		return err
	}
	return nil
}
