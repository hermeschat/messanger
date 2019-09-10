package eventhandlers

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/amirrezaask/config"
	"github.com/hashicorp/go-uuid"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
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
func saveMessageToDB(message *db.Message) error {

	_, err := db.Channels().UpdateOne(context.Background(), bson.M{
		"_id": message.ChannelID,
	}, bson.M{
		"$push": bson.M{
			"messages": message,
		},
	})
	if err != nil {
		return errors.Wrap(err, "cannot save message to mongo")
	}
	return nil
}

func saveChannelToDB(c *db.Channel) error {
	_, err := db.Channels().InsertOne(context.Background(), c)
	if err != nil {
		return errors.Wrap(err, "cannot save to mongo")
	}
	return nil
}

func getOrCreateExistingChannel(from string, to string) (*db.Channel, error) {
	logrus.Infof("creating/getting new channel to send new message from %s to %s", from, to)
	res := db.Channels().FindOne(context.Background(), bson.M{
		"members": bson.M{"$all": []string{from, to}, "$size": 2},
	})
	err := res.Err()
	fmt.Printf("%T\n", err)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			id, err := uuid.GenerateUUID()
			if err != nil {
				return nil, errors.Wrap(err, "cannot create uuid")
			}
			targetChannel := &db.Channel{
				ChannelId: id,
				Creator:   from,
				Members:   []string{to, from},
				Messages:  []db.Message{},
				Roles: map[string][]string{
					to:   {"R", "W", "M"},
					from: {"R", "W", "M"},
				},
				Type: db.ChatTypePrivate,
			}
			err = saveChannelToDB(targetChannel)
			if err != nil {
				return nil, errors.Wrap(err, "Error in mongo save")
			}
			return targetChannel, nil
		}
		return nil, errors.Wrap(res.Err(), "error in getting channel from database")
	}
	targetChannel := new(db.Channel)
	err = res.Decode(targetChannel)
	if err != nil {
		return nil, errors.Wrap(err, "err in decoding data into channel")
	}
	return targetChannel, nil
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
		logrus.Infof("user is not subscribed to channel %s", userID)
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
	conn, err := nats.Client(msg.From)
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
	var channels []string
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
	err := errors.New("some shitty error")
	for err != nil {
		err = f()
		if err != nil {
			logrus.Errorf("error in retrying operation: %v due to %v", name, err)
			time.Sleep(time.Second * 1)
			continue
		} else {
			break
		}
	}
}
