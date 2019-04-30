package newMessage

import (
	"context"
	"strings"
	"time"

	"git.raad.cloud/cloud/hermes/pkg/api"
	"git.raad.cloud/cloud/hermes/pkg/drivers/nats"
	"git.raad.cloud/cloud/hermes/pkg/drivers/redis"
	"git.raad.cloud/cloud/hermes/pkg/eventHandler"
	"git.raad.cloud/cloud/hermes/pkg/repository/channel"
	message2 "git.raad.cloud/cloud/hermes/pkg/repository/message"
	"git.raad.cloud/cloud/hermes/pkg/user_discovery"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
	"git.raad.cloud/cloud/hermes/pkg/repository"
)

func Handle(message *api.Message, sessionID string) *api.Response {
	var err error
	if message.To == "" && message.Channel == "" {
		return &api.Response{Error: errors.New("Channel ID or To should be present in payload").Error()}
	}
	targetChannel := &channel.Channel{}
	if message.To != "" {
		targetChannel, err = getOrCreateExistingChannel(message.From, message.To)
		if err != nil {
			logrus.Error(errors.Wrap(err, "failed to get or create channel"))
			return &api.Response{
				Error: "Internal service error",
			}
		}
	}
	if message.Channel != "" {
		targetChannel = &channel.Channel{ChannelID:message.Channel}
		// In case of dDos attack with lots of invalid channelid posted here, we should
		// check for channel existance in db or cache
	}
	go func (targetChannel *channel.Channel){
		if len(targetChannel.Members) < 1 || targetChannel.Members == nil {
			targetChannel, err = channel.Get(message.Channel)
			if err != nil {
				msg := errors.Wrap(err, "cannot get channel from db").Error()
				logrus.Error(msg)
			}
		}

		for _, member := range targetChannel.Members {
			err := ensureChannel(sessionID,targetChannel.ChannelID, member)
			if err != nil {
				go  retryEnsure(sessionID,targetChannel.ChannelID, member, 0)()
			}
		}
	}(targetChannel)

	//save to db
	go saveMessageToMongo(message)
	//Publish to nats
	err = nats.PublishNewMessage("test-cluster", "0.0.0.0:4222", targetChannel.ChannelID, message)
	if err != nil {
		return &api.Response{Error: errors.Wrap(err, "Cannot publish message to nats").Error()}
	}
	return &api.Response{
		Code: "200",
	}
}

func saveMessageToMongo(message *api.Message) error {
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

func getOrCreateExistingChannel(from string, to string)  (*channel.Channel,error) {
	channels, err := channel.GetAll(bson.M{
		"Members": bson.M{"$in": []string{from, to},"$size": 2},
	})
	if err != nil {
		return nil,errors.Wrap(err, "Cannot get channels")
	}
	var targetChannel *channel.Channel
	if len(*channels) < 1 {
		targetChannel = &channel.Channel{
			Members: []string{to, from},
		}
		err := saveChannelToMongo(targetChannel)
		if err != nil {
			return nil, errors.Wrap(err, "Error in mongo save")
		}
		return  targetChannel, nil
	} else {
		targetChannel = (*channels)[0]
		return targetChannel, nil
	}
}

func ensureChannel(sessionID string, channelID string, userID string) error {
	channels, err := getSession(sessionID)
	if err != nil {
		return errors.Wrap(err, "Error in get session from redis")
	}
	channelExist := false
	for _, c := range channels {
		if c == channelID {
			channelExist = true
		}
	}
	if !channelExist {
		//subscribeChannel(channelID, userID)
		//Send user discovery event
		//user discovery event publishes a userid and a chanellid
		//which this channels subscriber can listen to it and subscribe to given channel
		//its equivalent for subscribeChannel(channelID, userID) in async way
		err = user_discovery.PublishEvent(repository.UserDiscoveryEvent{ChannelID:channelID,UserID:userID})
		if err != nil {
			return errors.Wrap(err, "Error in publishing user discovery event")
		}
	}
	return nil
}

func retryEnsure(sessionID string,ChannelID string, userID string, retryCount int) func() {
	return func() {
		maxRetry := 10 // load from config
		if retryCount < maxRetry {
			err := ensureChannel(sessionID,ChannelID,userID)
			if err != nil {
				retryCount++
				time.Sleep(time.Millisecond * time.Duration( 100 * retryCount))
				retryEnsure(sessionID,ChannelID,userID, retryCount)()
			}
		}
	}
}

func getSession(sessionID string) ([]string, error) {
	redisCon, err := redis.ConnectRedis()
	if err != nil {
		return nil, errors.Wrap(err, "Fail to connect to redis")
	}
	channels, err := redisCon.Get("session-" + sessionID).Result()
	if err != nil {
		return nil, errors.Wrap(err, "Fail get from redis")
	}
	return strings.Split(channels, ","), nil
}

func subscribeChannel(channelID string, userID string) {
	ctx, _ := context.WithCancel(context.Background())
	sub := nats.MakeSubscriber(ctx, "test-cluster", "0.0.0.0:4222", channelID, eventHandler.NewMessageHandler(channelID, userID))
	go sub()
}
