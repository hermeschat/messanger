package discovery

import (
	"encoding/json"
	"fmt"
	"sync"
	"time"

	"hermes/api"
	"hermes/pkg/drivers/nats"
	"hermes/pkg/drivers/redis"

	"github.com/gogo/protobuf/proto"
	stan "github.com/nats-io/go-nats-streaming"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"golang.org/x/net/context"
)

const (
	UserDiscoveryChannel = "user-discovery"
)

//UserDiscoveryEvent is the eventhandlers we send to discovery channel to tell a user
//to subscribe to a certain channel in async way
type UserDiscoveryEvent struct {
	ChannelID string
	UserID    string
}

//NewMessageEventHandler handles the eventhandlers delivery from nats to user
func NewMessageEventHandler(channelID string, userID string, userSockets *struct {
	sync.RWMutex
	Us map[string]api.Hermes_EventBuffServer
}) func(msg *stan.Msg) {
	return func(msg *stan.Msg) {
		logrus.Warnf("22Message is %v", string(msg.Data))
		m := &api.Message{}
		err := json.Unmarshal(msg.Data, m)
		if err != nil {
			logrus.Errorf("error in unmarshalling nats eventhandlers in eventhandlers handler")
		}
		logrus.Infof("In NewMessage Event Handler as %s", userID)
		logrus.Infof("Recieved a new eventhandlers in %s", channelID)

		userSockets.RLock()
		userSocket, ok := userSockets.Us[userID]
		if !ok {
			logrus.Errorf("error: user socket not found ")
			return
		}
		err = userSocket.Send(&api.Event{Event: &api.Event_NewMessage{m}})
		if err != nil {
			logrus.Errorf("error: cannot send event new eventhandlers to user ")
			return
		}
		userSockets.RUnlock() //TODO: defer
	}
}

//UserDiscoveryEventHandler handles user discovery
func UserDiscoveryEventHandler(ctx context.Context, userID string, userSockets *struct {
	sync.RWMutex
	Us map[string]api.Hermes_EventBuffServer
}) func(msg *stan.Msg) {
	return func(msg *stan.Msg) {
		ude := &api.UserDiscoveryEvent{}
		err := ude.XXX_Unmarshal(msg.Data)
		if err != nil {
			logrus.Error(errors.Wrap(err, "cannot unmarshal UserDiscoveryEvent"))
			return
		}

		if ude.UserID == userID {
			channels, err := getSessionsByUserID(ude.UserID)
			if err != nil {
				logrus.Error(errors.Wrap(err, "Error in get session from redis"))
				return
			}
			channelExist := false
			for _, c := range channels {
				if c == ude.ChannelID {
					channelExist = true
				}
			}
			logrus.Warnf("%s is now subscribed to %s", ude.UserID, ude.ChannelID)
			if !channelExist {
				sub := nats.MakeSubscriber(ctx, userID, ude.ChannelID, NewMessageEventHandler(ude.ChannelID, ude.UserID, userSockets))
				go sub()
				go addSessionByUserID(ude.UserID, ude.ChannelID)
			}
		}
	}
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

func PublishUserDiscoveryEvent(ude UserDiscoveryEvent) error {

	u := &api.UserDiscoveryEvent{ChannelID: ude.ChannelID, UserID: ude.UserID}
	fmt.Println("client id is ", ude.UserID)
	conn, err := nats.NatsClient(ude.UserID)
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
