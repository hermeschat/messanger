package subscription

import (
	"context"
	"encoding/json"
	"time"

	"github.com/go-redis/redis"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"hermes/pkg/subscription/nats"
)

func NewSubsciption(ctx context.Context, userID string, channelID string, handler nats.Handler) error {
	if UserIsSubscribedTo(userID, channelID) {
		logrus.Infof("user %s already subscribed to channel %s", userID, channelID)
		return nil
	}
	go AddSubscriptionToUserID(userID, channelID)
	go nats.MakeSubscriber(ctx, userID, channelID, handler)()
	logrus.Infof("User %s is subscribed to %s\n", userID, channelID)
	return nil
}

func UserIsSubscribedTo(userID, channelID string) bool {
	chans, err := GetSubscribedChannels(userID)
	if err != nil {
		logrus.Errorf("error in getting subscribed channels: %v")
		return false
	}
	for _, ch := range chans {
		if ch == channelID {
			return true
		}
	}
	return false
}
func GetSubscribedChannels(userID string) ([]string, error) {
	con, err := Redis()
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

func AddSubscriptionToUserID(userID string, channelID string) error {
	channels, err := GetSubscribedChannels(userID)
	if err != nil {
		return errors.Wrap(err, "error while trying to get channels")
	}
	channels = append(channels, channelID)
	con, err := Redis()
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

func Clean() {
	con, err := Redis()
	if err != nil {
		logrus.Errorf("error in connecting to redis: %v", err)
		return
	}
	con.FlushDB()
}
