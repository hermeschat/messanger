package newMessage

import (
	"encoding/json"
	"fmt"
	"git.raad.cloud/cloud/hermes/pkg/repository"
	"git.raad.cloud/cloud/hermes/pkg/user_discovery"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"strings"
	"time"

	uuid "github.com/satori/go.uuid"

	"git.raad.cloud/cloud/hermes/pkg/drivers/nats"
	"git.raad.cloud/cloud/hermes/pkg/drivers/redis"
	"git.raad.cloud/cloud/hermes/pkg/repository/channel"
	message2 "git.raad.cloud/cloud/hermes/pkg/repository/message"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
)

type NewMessage struct {
	Session     string
	From        string
	To          string
	Channel     string
	MessageType string
	Body        string
}

func Handle(message *NewMessage) error {
	var err error
	if message.To == "" && message.Channel == "" {
		return errors.New("error in new message")
	}
	targetChannel := &channel.Channel{}
	if message.To != "" {
		targetChannel, err = getOrCreateExistingChannel(message.From, message.To)
		if err != nil {
			logrus.Error(errors.Wrap(err, "failed to get or create channel"))
			return errors.Wrap(err, "error in getting channel")
		}
	}
	if message.Channel != "" {
		targetChannel = &channel.Channel{ChannelID: message.Channel}
		// In case of dDos attack with lots of invalid channelid posted here, we should
		// check for channel existance in db or cache
	}
	logrus.Infof("target channel %+v", targetChannel)
	func(targetChannel *channel.Channel) {
		if len(targetChannel.Members) < 1 || targetChannel.Members == nil {
			targetChannel, err = channel.Get(targetChannel.ChannelID)
			if err != nil {
				msg := errors.Wrap(err, "cannot get channel from db").Error()
				logrus.Error(msg)
			}
		}

		for _, member := range targetChannel.Members {
			err := ensureChannel(message.Session, targetChannel.ChannelID, member)
			if err != nil {
				logrus.Errorf("error in ensuring channel : %v", err)
				//go retryEnsure(message.Session, targetChannel.ChannelID, member, 0)()
			}
		}
	}(targetChannel)
	message.Channel = targetChannel.ChannelID
	logrus.Infof("message is %+v", message)
	//save to db
	//err = saveMessageToMongo(message)
	//if err != nil {
	//	logrus.Errorf("cannot save message to mongodb :%v", err)
	//	return errors.Wrap(err, "error in saving message to mongo db")
	//}
	go saveMessageToMongo(message)
	logrus.Info("Trying To publish")
	//Publish to nats
	err = publishNewMessage("test-cluster", "0.0.0.0:4222", targetChannel.ChannelID, message)
	if err != nil {
		return errors.Wrap(err, "error in publishing message")
	}
	return nil
}

func saveMessageToMongo(message *NewMessage) error {
	//uid, err := uuid.NewV4()
	//if err != nil {
	//	logrus.Errorf("canot create uuid  : %v", err)
	//}
	err := message2.Add(&message2.Message{
		To:          message.To,
		From:        message.From,
		Time:        time.Now(),
		MessageType: message.MessageType,
		Body:        message.Body,
		MessageID:   primitive.NewObjectIDFromTimestamp(time.Now()),
		ChannelID:   message.Channel,
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

func getOrCreateExistingChannel(from string, to string) (*channel.Channel, error) {
	logrus.Infof("\ncreating/getting new channel to send message from %s to %s", from, to)
	channels, err := channel.GetAll(bson.M{
		"Members": bson.M{"$in": []string{from, to}, "$size": 2},
	})

	if err != nil {
		return nil, errors.Wrap(err, "Cannot get channels")
	}
	uid, err := uuid.NewV4()
	if err != nil {
		return nil, errors.Wrap(err, "can't create uuid")
	}
	var targetChannel *channel.Channel
	if len(channels) < 1 {
		logrus.Info("no channel found")
		targetChannel = &channel.Channel{
			Members:   []string{to, from},
			ChannelID: uid.String(),
			Roles: map[string][]string{
				to:   []string{"RWM"},
				from: []string{"RWM"},
			},
			Type: channel.Private,
		}
		err := saveChannelToMongo(targetChannel)
		if err != nil {
			return nil, errors.Wrap(err, "Error in mongo save")
		}
		return targetChannel, nil
	} else {
		targetChannel = (channels)[0]
		return targetChannel, nil

	}
}

func ensureChannel(sessionID string, channelID string, userID string) error {
	channels, err := getSessionsByUserID(sessionID)
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
		logrus.Infof("\nuser is not subscribed to channel %s", userID)
		//subscribeChannel(channelID, userID)
		//Send user discovery event
		//user discovery event publishes a userid and a chanellid
		//which this channels subscriber can listen to it and subscribe to given channel
		//its equivalent for subscribeChannel(channelID, userID) in async way
		err = user_discovery.PublishEvent(repository.UserDiscoveryEvent{ChannelID: channelID, UserID: userID})
		if err != nil {
			return errors.Wrap(err, "Error in publishing user discovery event")
		}
	}
	//err = user_discovery.PublishEvent(repository.UserDiscoveryEvent{ChannelID: channelID, UserID: userID})
	//if err != nil {
	//	return errors.Wrap(err, "Error in publishing user discovery event")
	//}
	return nil
}

func retryEnsure(sessionID string, ChannelID string, userID string, retryCount int) func() {
	return func() {
		maxRetry := 10 // load from config
		if retryCount < maxRetry {
			err := ensureChannel(sessionID, ChannelID, userID)
			if err != nil {
				retryCount++
				time.Sleep(time.Millisecond * time.Duration(100*retryCount))
				retryEnsure(sessionID, ChannelID, userID, retryCount)()
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

//publishNewMessage is send function. Every message should be published to a channel to
//be delivered to subscribers. In streaming, published Message is persistant.
func publishNewMessage(clusterID string, natsSrvAddr string, ChannelId string, msg *NewMessage) error {
	// Connect to NATS-Streaming
	logrus.Info("trying to connect to nats")
	//logrus.Info(msg.From)
	conn, err := nats.NatsClient(clusterID, natsSrvAddr, msg.From)
	if err != nil {
		return errors.Wrapf(err, "Can't connect: %v.\nMake sure a NATS Streaming Server is running at: %s", err, natsSrvAddr)
	}
	bs, err := json.Marshal(msg)
	if err != nil {
		return errors.Wrap(err, "error in marshaling message")
	}
	if err := (*conn).Publish(ChannelId, bs); err != nil {
		return errors.Wrap(err, "failed to publish message")
	}
	logrus.Info("trying to publish")
	logrus.Infof("Published into %s a new message as %s", ChannelId, msg.From)
	return nil
}

func getSessionsByUserID(userID string) ([]string, error) {
	con, err := redis.ConnectRedis()
	if err != nil {
		return nil, errors.Wrap(err, "error while trying to connect to redis")
	}
	defer con.Close()
	channels := []string{}
	res, err := con.Get(userID).Result()
	if err != nil {
		if err == redis.Nil {
			return []string{}, nil
		}
		return nil, errors.Wrap(err, "error while scanning redis output")
	}
	fmt.Println(">>>>>>>>>>>" + res)
	err = json.Unmarshal([]byte(res), &channels)
	if err != nil {

		return nil, errors.Wrap(err, "error while Unmarshalling redis output")
	}
	return channels, nil
}
