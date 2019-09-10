package eventhandlers

import (
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
	logrus.Infof("target channel %+v", tc)
	tc, err = loadChannelData(tc)
	if err != nil {
		return errors.Wrap(err, "error in loading members")
	}
	message.ChannelID = tc.ChannelId
	go retryOp("saving message to mongodb", func() error { return saveMessageToMongo(message) })
	for _, member := range tc.Members {
		go retryOp("esuring every one of the members are subscribed to channel", func() error { return ensureChannel(tc.ChannelId, member) })
	}
	if !hasWriteRole(message.From, tc) {
		return errors.Wrap(err, "error, access denied")
	}
	go retryOp("publish new message", func() error {
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
				logrus.Error(errors.Wrap(err, "failed to get or create channel"))
				return nil, errors.Wrap(err, "error in getting channel")
			}
		} else {
			return nil, errors.Wrap(err, "no valid receiver whether channel or userId found")
		}
	}
	return targetChannel, nil
}
func loadChannelData(tc *db.Channel) (*db.Channel, error) {
	if len(tc.Members) < 1 || tc.Members == nil {
		ch, err := db.Instance().Channels.Find(tc.ChannelId)
		if err != nil {
			return nil, errors.Wrap(err, "error in finding channel by channel id")
		}
		err = tc.FromMap(ch)
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
