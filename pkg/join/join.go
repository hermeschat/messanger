package join

import (
	"context"
	"encoding/json"
	"hermes/api/pb"
	"hermes/pkg/discovery"
	"hermes/pkg/drivers/nats"
	"hermes/pkg/drivers/redis"
	"sync"
	"time"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

//JoinPayload ...
type JoinPayload struct {
	UserID    string
	SessionId string
}

//Handle handles join event
func Handle(ctx context.Context, sig *JoinPayload, userSockets *struct {
	sync.RWMutex
	Us map[string]pb.Hermes_EventBuffServer
}) {
	channels, err := getSessionsByUserID(sig.UserID)
	if err != nil {
		logrus.Error(errors.Wrap(err, "error while trying to get channels from redis").Error())
		return
	}
	udSub := false
	for _, c := range channels {
		if c == "user-discovery" {
			udSub = true
		}
	}
	if !udSub {
		logrus.Info("Not subscribed to user-discovery")

		err = addSessionByUserID(sig.UserID, "user-discovery")
		if err != nil {
			logrus.Error("error while trying to add session to redis")
			return
		}
		logrus.Infof("Subscribing to user-discovery as %s", sig.UserID)
		sub := nats.MakeSubscriber(ctx, sig.UserID, "test-cluster", "0.0.0.0:4222", "user-discovery", discovery.UserDiscoveryEventHandler(ctx, sig.UserID, userSockets))
		go sub()
		return
	}
	logrus.Info("Already subscribing to user-discovery")
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
