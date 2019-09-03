package eventHandler

import (
	"encoding/json"
	"fmt"
	"strings"
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

//UserDiscoveryEventHandler handles user discovery
func UserDiscoveryEventHandler(ctx context.Context, userID string, currentSession string) func(msg *stan.Msg) {
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
				sub := nats.MakeSubscriber(ctx, userID, "test-cluster", "0.0.0.0:4222", ude.ChannelID, NewMessageEventHandler(ude.ChannelID, ude.UserID))
				go sub()
				return
			}
		}
	}
}

var UserSockets = struct {
	sync.RWMutex
	Us map[string]pb.Hermes_EventBuffServer
}{
	sync.RWMutex{},
	map[string]pb.Hermes_EventBuffServer{},
}

//NewMessageEventHandler handles the message delivery from nats to user
func NewMessageEventHandler(channelID string, userID string) func(msg *stan.Msg) {
	return func(msg *stan.Msg) {
		logrus.Warnf("22Message is %v", string(msg.Data))
		m := &pb.Message{}
		err := json.Unmarshal(msg.Data, m)
		//_ ,err := m.XXX_Marshal(msg.Data, false)
		if err != nil {
			logrus.Errorf("error in unmarshalling nats message in message handler")
		}
		logrus.Infof("In NewMessage Event Handler as %s", userID)
		logrus.Infof("Recieved a new message in %s", channelID)
		//c, ok := BaseHub.ClientsMap[userID]
		//if !ok {
		//	logrus.Error("no active connection found for user")
		//	return
		//}

		//c.send <- msg.Data
		UserSockets.RLock()
		userSocket, ok := UserSockets.Us[userID]
		if !ok {
			logrus.Errorf("error: user socket not found ")
			return
		}
		err = userSocket.Send(&pb.Event{Event: &pb.Event_NewMessage{m}})
		if err != nil {
			logrus.Errorf("error: cannot send event new message to user ")
			return
		}
		UserSockets.RUnlock()
	}
}

func subscribeChannel(ctx context.Context, channelID string, userID string) {
	//ctx, _ := context.WithTimeout(context.Background(), time.Hour * 1)
	sub := nats.MakeSubscriber(ctx, userID, "test-cluster", "0.0.0.0:4222", channelID, NewMessageEventHandler(channelID, userID))
	go sub()
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

//JoinPayload ...
type JoinPayload struct {
	UserID    string
	SessionId string
}

func Handle(ctx context.Context, sig *JoinPayload) {

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
		sub := nats.MakeSubscriber(ctx, sig.UserID, "test-cluster", "0.0.0.0:4222", "user-discovery", UserDiscoveryEventHandler(ctx, sig.UserID, ""))
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
