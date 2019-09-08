package eventhandlers

import (
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"hermes/pkg/db"
)

func HandleNewMessage(message *db.Message) error {
	if message.To == "" && message.ChannelID == "" {
		return errors.New("error in new new message event handler")
	}

	tc, err := targetChannel(message)
	if err != nil {
		return errors.Wrap(err, "error in finding target channel")
	}
	logrus.Infof("target channel %+v", tc)
	tc, err = loadMembers(tc)
	if err != nil {
		return errors.Wrap(err, "error in loading members")
	}
	message.ChannelID = tc.ChannelID
	go retryOp("saving message to mongodb", func() error { return saveMessageToMongo(message) })
	for _, member := range tc.Members {
		err := ensureChannel(tc.ChannelID, member)
		if err != nil {
			logrus.Errorf("error in ensuring channel : %v", err)
			//go retryEnsure(eventhandlers.Session, targetChannel.ChannelID, member, 0)()
		}
	}
	// roles := targetChannel.Roles[eventhandlers.From]
	// if checkRoles(roles[0]) { //TODO: fix roles to be array of string not single string in array
	// 	return errors.New("user doesn't have write permission in this channel")
	// }
	logrus.Info("Trying To publish")
	err = publishNewMessage("test-cluster", "0.0.0.0:4222", tc.ChannelID, message)
	if err != nil {
		return errors.Wrap(err, "error in publishing eventhandlers")
	}
	return nil
}

func targetChannel(message *db.Message) (*db.Channel, error) {
	targetChannel := &db.Channel{}
	var err error
	if message.ChannelID != "" {
		targetChannel.ChannelID = message.ChannelID
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
	return nil, nil
}
func loadMembers(tc *db.Channel) (*db.Channel, error) {
	if len(tc.Members) < 1 || tc.Members == nil {
		ch, err := db.Instance().Channels.Find(tc.ChannelID)
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
