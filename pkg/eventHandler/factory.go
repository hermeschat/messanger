package eventHandler

import (
	"fmt"
	"git.raad.cloud/cloud/hermes/pkg/api"
	"git.raad.cloud/cloud/hermes/pkg/drivers/nats"
	"git.raad.cloud/cloud/hermes/pkg/drivers/redis"
	stan "github.com/nats-io/go-nats-streaming"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"golang.org/x/net/context"
	"strings"
	"sync"
)

//UserDiscoveryEventHandler handles user discovery
func UserDiscoveryEventHandler(ctx context.Context,userID string, currentSession string) func(msg *stan.Msg) {
	return func(msg *stan.Msg) {
		ude := &api.UserDiscoveryEvent{}
		err := ude.XXX_Unmarshal(msg.Data)
		if err != nil {
			logrus.Error(errors.Wrap(err, "cannot unmarshal UserDiscoveryEvent"))
			return
		}
		//if ude.UserID == userID {
		//	channels, err := getSession(currentSession)
		//	if err != nil {
		//		logrus.Error(errors.Wrap(err, "Error in get session from redis"))
		//	}
		//	channelExist := false
		//	for _, c := range channels {
		//		if c == ude.ChannelID {
		//			channelExist = true
		//		}
		//	}
		//	channelExist := false

		logrus.Infof("%s is now subscribed to %s", ude.UserID, ude.ChannelID)
		sub := nats.MakeSubscriber(ctx, userID, "test-cluster", "0.0.0.0:4222", ude.ChannelID, NewMessageEventHandler(ude.ChannelID, ude.UserID))
		go sub()
		if err != nil {
			logrus.Errorf("Error in subscribing to channel", err)
		}
	}
	}
var UserSockets = struct {
	sync.RWMutex
	Us map[string]*api.Hermes_EventBuffServer
}{
	sync.RWMutex{},
	map[string]*api.Hermes_EventBuffServer{},
}


//NewMessageEventHandler handles the message delivery from nats to user
func NewMessageEventHandler(channelID string, userID string) func(msg *stan.Msg) {
	return func(msg *stan.Msg) {
		logrus.Info("In NewMessage Event Handler")

		UserSockets.RLock()
		err := (*UserSockets.Us[userID]).Send(&api.Event{Event:&api.Event_NewMessage{&api.Message{
			From: "server",
		}}})
		UserSockets.RUnlock()
		logrus.Errorf("cannot push to client ", err)
		fmt.Printf("Recieved New Message In %s", channelID)
	}
}

func subscribeChannel(ctx context.Context, channelID string, userID string) {
	//ctx, _ := context.WithTimeout(context.Background(), time.Hour * 1)
	sub := nats.MakeSubscriber(ctx,userID,"test-cluster", "0.0.0.0:4222", channelID, NewMessageEventHandler(channelID, userID))
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
