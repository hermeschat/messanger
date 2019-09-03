package message

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"hermes/pkg/db"
	"hermes/pkg/user_discovery"

	"github.com/mitchellh/mapstructure"
	uuid "github.com/satori/go.uuid"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
	"hermes/pkg/drivers/nats"
	"hermes/pkg/drivers/redis"
)

func Handle(message *db.Message) error {
	var err error
	logrus.Infof("######$$$$$$8==> %+v", message)
	if message.To == "" && message.ChannelID == "" {
		return errors.New("error in new message")
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
		ch, err := db.Instance().Repo("channels").Find(targetChannel.ChannelID)
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
		return errors.Wrap(err, "error in saving message to db")
	}
	for _, member := range targetChannel.Members {
		err := ensureChannel(targetChannel.ChannelID, member)
		if err != nil {
			logrus.Errorf("error in ensuring channel : %v", err)
			//go retryEnsure(message.Session, targetChannel.ChannelID, member, 0)()
		}
	}
	//}(targetChannel)

	// roles := targetChannel.Roles[message.From]
	// if checkRoles(roles[0]) { //TODO: fix roles to be array of string not single string in array
	// 	return errors.New("user doesn't have write permission in this channel")
	// }
	logrus.Infof("message is %+v", message)
	//save to db
	//err = saveMessageToMongo(message)
	//if err != nil {
	//	logrus.Errorf("cannot save message to mongodb :%v", err)
	//	return errors.Wrap(err, "error in saving message to mongo db")
	//}

	logrus.Info("Trying To publish")
	//Publish to nats
	err = publishNewMessage("test-cluster", "0.0.0.0:4222", targetChannel.ChannelID, message)
	if err != nil {
		return errors.Wrap(err, "error in publishing message")
	}
	return nil
}

func checkRoles(roles string) bool {
	for _, c := range roles {
		if string(c) == "W" {
			return true
		}
	}
	return false
}
func saveMessageToMongo(message *db.Message) error {
	m, err := message.ToMap()
	if err != nil {
		return errors.Wrap(err, "error in converting message to map")
	}
	_, err = db.Instance().Repo("messages").Add(m)
	if err != nil {
		return errors.Wrap(err, "cannot save to mongo")
	}
	return nil
}

func saveChannelToMongo(c *db.Channel) error {
	m, err := c.ToMap()
	if err != nil {
		return errors.Wrap(err, "error in converting channel to map")
	}

	_, err = db.Instance().Repo("channels").Add(m)
	if err != nil {
		return errors.Wrap(err, "cannot save to mongo")
	}
	return nil
}

func getOrCreateExistingChannel(from string, to string) (*db.Channel, error) {
	logrus.Infof("\ncreating/getting new channel to send message from %s to %s", from, to)
	channels, err := db.Instance().Repo("channels").Get(bson.M{
		"Members": bson.M{"$all": []string{from, to}, "$size": 2},
	})
	cts := []*db.Channel{}
	for _, c := range channels {
		ct := &db.Channel{}
		err := mapstructure.Decode(c, ct)
		if err != nil {
			return nil, errors.Wrap(err, "error in converting channle map to channel type")
		}
		cts = append(cts, ct)

	}
	if err != nil {
		return nil, errors.Wrap(err, "Cannot get channels")
	}
	uid := uuid.NewV4()

	var targetChannel *db.Channel
	if len(channels) < 1 {
		logrus.Info("no channel found")
		targetChannel = &db.Channel{
			Members:   []string{to, from},
			ChannelID: uid.String(),
			Roles: map[string][]string{
				to:   []string{"RWM"},
				from: []string{"RWM"},
			},
			Type: db.Private,
		}
		err := saveChannelToMongo(targetChannel)
		if err != nil {
			return nil, errors.Wrap(err, "Error in mongo save")
		}
		return targetChannel, nil
	} else {
		targetChannel = (cts)[0]
		return targetChannel, nil

	}
}

func ensureChannel(channelID string, userID string) error {
	channels, err := getSubscribedChannels(userID)
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
		err = user_discovery.PublishEvent(user_discovery.UserDiscoveryEvent{ChannelID: channelID, UserID: userID})
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
			err := ensureChannel(ChannelID, userID)
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
func publishNewMessage(clusterID string, natsSrvAddr string, ChannelId string, msg *db.Message) error {
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

func getSubscribedChannels(userID string) ([]string, error) {
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
