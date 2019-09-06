package db

import (
	"time"

	"github.com/mitchellh/mapstructure"
)

type Message struct {
	MessageID   string
	From        string
	To          string
	Time        string
	ChannelID   string
	MessageType string
	Body        string
	Read        bool
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

type User struct {
	AccountID      string `bson:"_id" json:"_id"`
	UserName       string
	Name           string
	ProfilePicture string
	MobilePhone    string
	Active         bool
	CreatedAt      time.Time
	UpdatedAt      time.Time
}

type Channel struct {
	ChannelID string              `json:"channelID" bson:"ChannelID"`
	Members   []string            `json:"Members" bson:"Members"`
	Creator   string              `json:"Creator" bson:"Creator"`
	Type      int                 `json:"type" bson:"type"`
	Roles     map[string][]string `json:"roles" bson:"roles"`
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
