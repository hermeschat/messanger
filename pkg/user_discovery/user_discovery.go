package user_discovery

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/gogo/protobuf/proto"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"hermes/api/pb"
	"hermes/pkg/drivers/nats"
	"hermes/pkg/drivers/redis"
)

func PublishEvent(ude UserDiscoveryEvent) error {

	u := &pb.UserDiscoveryEvent{ChannelID: ude.ChannelID, UserID: ude.UserID}
	fmt.Println("client id is ", ude.UserID)
	conn, err := nats.NatsClient("test-cluster", "0.0.0.0:4222", ude.UserID)
	if err != nil {
		return errors.Wrap(err, "cannot connect to nats")
	}
	bs, err := proto.Marshal(u)
	if err != nil {
		return errors.Wrap(err, "cannot marshal UserDiscoveryEvent")
	}
	err = (*conn).Publish("user-discovery", bs)
	if err != nil {
		return errors.Wrap(err, "cannot publish UserDiscoveryEvent")
	}
	logrus.Infof("Published User Discovery event %+v", u)
	return nil

}

//UserDiscoveryEvent is the message we send to discovery channel to tell a user
//to subscribe to a certain channel in async way
type UserDiscoveryEvent struct {
	ChannelID string
	UserID    string
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

func addSessionByUserID(userID string, channelID string) error {
	channels, err := getSessionsByUserID(userID)
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
