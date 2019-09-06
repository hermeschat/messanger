package eventHandler

import (
	"encoding/json"
	"sync"
	"time"

	stan "github.com/nats-io/go-nats-streaming"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"golang.org/x/net/context"
	"hermes/api/pb"
	"hermes/pkg/drivers/nats"
	"hermes/pkg/drivers/redis"
)

//JoinPayload ...
type JoinPayload struct {
	UserID    string
	SessionId string
}

//NewMessageEventHandler handles the message delivery from nats to user
func NewMessageEventHandler(channelID string, userID string, userSockets struct {
	sync.RWMutex
	Us map[string]pb.Hermes_EventBuffServer
}) func(msg *stan.Msg) {
	return func(msg *stan.Msg) {
		logrus.Warnf("22Message is %v", string(msg.Data))
		m := &pb.Message{}
		err := json.Unmarshal(msg.Data, m)
		if err != nil {
			logrus.Errorf("error in unmarshalling nats message in message handler")
		}
		logrus.Infof("In NewMessage Event Handler as %s", userID)
		logrus.Infof("Recieved a new message in %s", channelID)

		userSockets.RLock()
		userSocket, ok := userSockets.Us[userID]
		if !ok {
			logrus.Errorf("error: user socket not found ")
			return
		}
		err = userSocket.Send(&pb.Event{Event: &pb.Event_NewMessage{m}})
		if err != nil {
			logrus.Errorf("error: cannot send event new message to user ")
			return
		}
		userSockets.RUnlock() //TODO: defer
	}
}

func Handle(ctx context.Context, sig *JoinPayload, userSockets struct {
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
		sub := nats.MakeSubscriber(ctx, sig.UserID, "test-cluster", "0.0.0.0:4222", "user-discovery", UserDiscoveryEventHandler(ctx, sig.UserID, userSockets))
		go sub()
		return
	}
	logrus.Info("Already subscribing to user-discovery")
}

//UserDiscoveryEventHandler handles user discovery
func UserDiscoveryEventHandler(ctx context.Context, userID string, userSockets struct {
	sync.RWMutex
	Us map[string]pb.Hermes_EventBuffServer
}) func(msg *stan.Msg) {
	return func(msg *stan.Msg) {
		logrus.Info("!!!!!!!!!!!!!!!!discovery event handler called ")
		ude := &pb.UserDiscoveryEvent{}
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
				sub := nats.MakeSubscriber(ctx, userID, "test-cluster", "0.0.0.0:4222", ude.ChannelID, NewMessageEventHandler(ude.ChannelID, ude.UserID, userSockets))
				go sub()
				return
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
