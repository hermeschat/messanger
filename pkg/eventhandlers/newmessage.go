package eventhandlers

import (
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"hermes/pkg/db"
)

func HandleNewMessage(message *db.Message) error {
	var err error
	logrus.Infof("######$$$$$$8==> %+v", message)
	if message.To == "" && message.ChannelID == "" {
		return errors.New("error in new eventhandlers")
	}
	targetChannel := &db.Channel{}
	if message.ChannelID != "" {
		targetChannel.ChannelID = message.ChannelID
	} else {
		if message.To != "" {
			targetChannel, err = getOrCreateExistingChannel(message.From, message.To)
			if err != nil {
				logrus.Error(errors.Wrap(err, "failed to get or create channel"))
				return errors.Wrap(err, "error in getting channel")
			}
		} else {
			return errors.Wrap(err, "no valid receiver whether channel or userId found")
		}
	}

	logrus.Infof("target channel %+v", targetChannel)
	//func(targetChannel *channel.Channel) {
	if len(targetChannel.Members) < 1 || targetChannel.Members == nil {
		ch, err := db.Instance().Channels.Find(targetChannel.ChannelID)
		if err != nil {
			msg := errors.Wrap(err, "cannot get channel from db").Error()
			logrus.Error(msg)
		}
		err = targetChannel.FromMap(ch)
		if err != nil {
			logrus.Errorf("erorr while converting from map to channel:%v", err)
		}

	}
	message.ChannelID = targetChannel.ChannelID
	err = saveMessageToMongo(message)
	if err != nil {
		return errors.Wrap(err, "error in saving eventhandlers to db")
	}
	for _, member := range targetChannel.Members {
		err := ensureChannel(targetChannel.ChannelID, member)
		if err != nil {
			logrus.Errorf("error in ensuring channel : %v", err)
			//go retryEnsure(eventhandlers.Session, targetChannel.ChannelID, member, 0)()
		}
	}
	//}(targetChannel)

	// roles := targetChannel.Roles[eventhandlers.From]
	// if checkRoles(roles[0]) { //TODO: fix roles to be array of string not single string in array
	// 	return errors.New("user doesn't have write permission in this channel")
	// }
	logrus.Infof("eventhandlers is %+v", message)
	//save to db
	//err = saveMessageToMongo(eventhandlers)
	//if err != nil {
	//	logrus.Errorf("cannot save eventhandlers to mongodb :%v", err)
	//	return errors.Wrap(err, "error in saving eventhandlers to mongo db")
	//}

	logrus.Info("Trying To publish")
	//Publish to nats
	err = publishNewMessage("test-cluster", "0.0.0.0:4222", targetChannel.ChannelID, message)
	if err != nil {
		return errors.Wrap(err, "error in publishing eventhandlers")
	}
	return nil
}
