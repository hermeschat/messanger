package eventhandlers

import (
	"context"
	"time"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"hermes/pkg/db"
)

func HandleNewMessage(message *db.Message) error {
	logrus.Infof("%+v\n", message)
	if message.To == "" && message.ChannelID == "" {
		return errors.New("error in new new message event handler")
	}

	tc, err := targetChannel(message)
	if err != nil {
		return errors.Wrap(err, "error in finding target channel")
	}
	tc, err = loadChannelData(tc)
	if err != nil {
		return errors.Wrap(err, "error in loading members")
	}
	message.ChannelID = tc.ChannelId
	message.Time = time.Now()
	message.MessageType = db.MessageTypeText
	go retryFunc("saving message to mongodb", func() error { return saveMessageToDB(message) })
	for _, member := range tc.Members {
		logrus.Infof("ensuring that %s is subscribed to %s", member, tc.ChannelId)
		//go retryFunc("ensuring every one of the members are subscribed to channel", func() error { return ensureChannel(tc.ChannelId, member) })
		err := ensureChannel(tc.ChannelId, member)
		if err != nil {
			logrus.Errorf("%s error in ensureChannel: %v", member, err)
		}
	}
	if !hasWriteRole(message.From, tc) {
		return errors.Wrap(err, "error, access denied")
	}
	//time.Sleep(time.Second * 3)
	go retryFunc("publish new message", func() error {
		return publishNewMessage(tc.ChannelId, message)
	})
	return nil
}

func targetChannel(message *db.Message) (*db.Channel, error) {
	targetChannel := &db.Channel{}
	var err error
	if message.ChannelID != "" {
		targetChannel.ChannelId = message.ChannelID
	} else {
		if message.To != "" {
			targetChannel, err = getOrCreateExistingChannel(message.From, message.To)
			if err != nil {
				return nil, errors.Wrap(err, "error in getting channel")
			}
		} else {
			return nil, errors.Wrap(err, "no valid receiver whether channel or userId found")
		}
	}
	return targetChannel, nil
}
func loadChannelData(tc *db.Channel) (*db.Channel, error) {
	if len(tc.Roles) < 1 || tc.Roles == nil {
		result := db.Channels().FindOne(context.Background(), map[string]string{
			"channel_id": tc.ChannelId,
		})
		if result.Err() != nil {
			return nil, errors.Wrap(result.Err(), "error in finding channel by channel id")
		}
		tc := new(db.Channel)
		err := result.Decode(tc)
		if err != nil {
			logrus.Errorf("erorr while converting from map to channel:%v", err)
			return nil, errors.Wrap(err, "error in decoding map into target channel")
		}
	} else {
		return tc, nil
	}
	return nil, nil
}

func hasWriteRole(userID string, channel *db.Channel) bool {

	return hasRole(channel.Roles[userID], "W")
}
