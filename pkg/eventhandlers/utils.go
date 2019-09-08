package eventhandlers

import (
	"encoding/json"
	"strings"
	"time"

	"github.com/amirrezaask/config"
	"github.com/mitchellh/mapstructure"
	"github.com/pkg/errors"
	uuid "github.com/satori/go.uuid"
	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
	"hermes/pkg/db"
	"hermes/pkg/discovery"
	"hermes/pkg/drivers/nats"
	"hermes/pkg/drivers/redis"
)

func hasRole(roles []string, role string) bool {
	for _, c := range roles {
		if string(c) == role {
			return true
		}
	}
	return false
}
func saveMessageToMongo(message *db.Message) error {
	m, err := message.ToMap()
	if err != nil {
		return errors.Wrap(err, "error in converting eventhandlers to map")
	}
	_, err = db.Instance().Messages.Add(m)
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

	_, err = db.Instance().Channels.Add(m)
	if err != nil {
		return errors.Wrap(err, "cannot save to mongo")
	}
	return nil
}

func getOrCreateExistingChannel(from string, to string) (*db.Channel, error) {
	logrus.Infof("\ncreating/getting new channel to send eventhandlers from %s to %s", from, to)
	channels, err := db.Instance().Channels.Get(bson.M{
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
		err = discovery.PublishUserDiscoveryEvent(discovery.UserDiscoveryEvent{ChannelID: channelID, UserID: userID})
		if err != nil {
			return errors.Wrap(err, "Error in publishing user discovery event")
		}
	}
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

//publishNewMessage is send function. Every eventhandlers should be published to a channel to
//be delivered to subscribers. In streaming, published Message is persistant.
func publishNewMessage(ChannelId string, msg *db.Message) error {
	logrus.Info("trying to connect to nats")
	conn, err := nats.NatsClient(config.Get("cluster_id"), config.Get("nats_host"), msg.From)
	if err != nil {
		return errors.Wrapf(err, "Can't connect: %v.\nMake sure a NATS Streaming Server is running at: %s", err, config.Get("nats_host"))
	}
	bs, err := json.Marshal(msg)
	if err != nil {
		return errors.Wrap(err, "error in marshaling new message event")
	}
	if err := (*conn).Publish(ChannelId, bs); err != nil {
		return errors.Wrap(err, "failed to publish new message event")
	}
	logrus.Infof("Published into %s a new message event as %s", ChannelId, msg.From)
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
	err = json.Unmarshal([]byte(res), &channels)
	if err != nil {

		return nil, errors.Wrap(err, "error while Unmarshalling redis output")
	}
	return channels, nil
}

func addSessionByUserID(userID string, channelID string) error {
	channels, err := getSubscribedChannels(userID)
	if err != nil {
		return errors.Wrap(err, "error while trying to get channels")
	}
	channels = append(channels, channelID)
	con, err := redis.ConnectRedis()
	if err != nil {
		return errors.Wrap(err, "error while trying to connect to redis")
	}
	defer con.Close()
	bs, err := json.Marshal(channels)
	if err != nil {
		return errors.Wrap(err, "error while trying to marshal channels")
	}
	err = con.Set(userID, string(bs), time.Hour*1).Err()
	if err != nil {
		return errors.Wrap(err, "error while adding new channels to redis")
	}
	return nil
}

func retryOp(name string, f func() error) {
	var err error
	for err == nil {
		err = f()
		if err != nil {
			logrus.Errorf("error in retrying operation: %v", name)
			continue
		} else {
			break
		}
	}
}
