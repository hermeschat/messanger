package new_message

import (
	"strings"
	"time"

	"git.raad.cloud/cloud/credit-service/db"
	"git.raad.cloud/cloud/hermes/pkg/api"
	"git.raad.cloud/cloud/hermes/pkg/drivers/nats"
	"git.raad.cloud/cloud/hermes/pkg/repository/channel"
	message2 "git.raad.cloud/cloud/hermes/pkg/repository/message"
	"git.raad.cloud/cloud/hermes/pkg/user_discovery"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
)

func Handle(message *api.InstantMessage) *api.Response {
	if message.To == "" && message.Channel == "" {
		return &api.Response{Error: errors.New("Channel ID or To should be present in payload").Error()}
	}
	if message.Channel != "" {

		err := nats.PublishNewMessage("test-cluster", "0.0.0.0:4222", message.Channel, message)
		if err != nil {
			return &api.Response{Error: errors.Wrap(err, "Cannot publish message to nats").Error()}
		}
		targetChannel, err := channel.Get(message.Channel)
		if err != nil {
			msg := errors.Wrap(err, "cannot get channel from db").Error()
			logrus.Error(msg)
			return &api.Response{
				Error: msg,
			}

		}
		for _, member := range targetChannel.Members {
			err := user_discovery.PublishEvent(&api.UserDiscoveryEvent{
				UserID:    member,
				ChannelID: targetChannel.ChannelID,
			})
			if err != nil {
				logrus.Error(errors.Wrap(err, "cannot publish to user-discovery"))
			}
		}
		return &api.Response{
			Code: "200",
		}
	}
	if message.To != "" {
		channels, err := channel.GetAll(bson.M{
			"Members": bson.M{"$in": []string{message.From, message.To}},
		})
		if err != nil {
			return &api.Response{
				Error: errors.Wrap(err, "Cannot get channels").Error(),
			}
		}
		var targetChannel *channel.Channel
		if len(*channels) < 1 {
			targetChannel = &channel.Channel{
				Members: []string{message.To, message.From},
			}
			err := saveChannelToMongo(targetChannel)
			if err != nil {
				return &api.Response{
					Error: "Internal Service problem",
				}
			}
		} else {
			targetChannel = (*channels)[0]
		}
		err = nats.PublishNewMessage("test-cluster", "0.0.0.0", targetChannel.ChannelID, message)
		if err != nil {
			return &api.Response{
				Error: errors.Wrap(err, "Error while publishing to NATS").Error(),
			}
		}
		err = saveMessageToMongo(message)
		if err != nil {
			return &api.Response{
				Error: errors.Wrap(err, "Error in saving to mongo").Error(),
			}
		}
		for _, member := range targetChannel.Members {
			err := user_discovery.PublishEvent(&api.UserDiscoveryEvent{
				UserID:    member,
				ChannelID: targetChannel.ChannelID,
			})
			if err != nil {
				logrus.Error(errors.Wrap(err, "cannot publish to user-discovery"))
			}
		}
		return &api.Response{
			Code: "200",
		}

	}
	return &api.Response{
		Code: "Unknown",
	}
}

func saveMessageToMongo(message *api.InstantMessage) error {
	err := message2.Add(&message2.Message{
		To:          message.To,
		From:        message.From,
		Time:        time.Now(),
		MessageType: message.MessageType,
		Body:        message.Body,
	})
	if err != nil {
		return errors.Wrap(err, "cannot save to mongo")
	}
	return nil
}

func saveChannelToMongo(c *channel.Channel) error {
	err := channel.Add(c)
	if err != nil {
		return errors.Wrap(err, "cannot save to mongo")
	}
	return nil
}

func getExistingChannel(from string, to string) {
	channels, err := channel.GetAll(bson.M{
		"Members": bson.M{"$in": []string{message.From, message.To}, {"$size": 2}},
	})
	if err != nil {
		return &api.Response{
			Error: errors.Wrap(err, "Cannot get channels").Error(),
		}
	}
	var targetChannel *channel.Channel
	if len(*channels) < 1 {
		targetChannel = &channel.Channel{
			Members: []string{message.To, message.From},
		}
		err := saveChannelToMongo(targetChannel)
		if err != nil {
			return &api.Response{
				Error: "Internal Service problem",
			}
		}
	} else {
		targetChannel = (*channels)
	}
}

func ensureChannel(sessionID string, channelID string) {
	channels, err := getSession(sessionID)
	for _, c := range channel {

	}
}

func getSession(sessionID string) ([]string, error) {
	redisCon, err := db.ConnectRedis()
	if err != nil {
		return nil, errors.Wrap(err, "Fail to connect to redis")
	}
	channels, err := redisCon.Get("session-" + sessionID).Result()
	if err != nil {
		return nil, errors.Wrap(err, "Fail get from redis")
	}
	return strings.Split(channels, ","), nil
}
