package discovery

import (
	"encoding/json"
	"sync"

	stan "github.com/nats-io/go-nats-streaming"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"golang.org/x/net/context"
	"hermes/api"
	"hermes/subscription"
	"hermes/subscription/nats"
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
		m := &api.Message{}
		err := json.Unmarshal(msg.Data, m)
		if err != nil {
			logrus.Errorf("error in unmarshalling nats eventhandlers in eventhandlers handler")
			return
		}
		logrus.Warnf("in new Message event handler and message is %v", string(msg.Data))
		logrus.Infof("In NewMessage Event Handler as %s", userID)

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
		userSockets.RUnlock()
	}
}

//UserDiscoveryEventHandler handles user discovery
func UserDiscoveryEventHandler(ctx context.Context, userID string, userSockets *struct {
	sync.RWMutex
	Us map[string]api.Hermes_EventBuffServer
}) func(msg *stan.Msg) {
	return func(msg *stan.Msg) {
		logrus.Infoln(">>>>>>>>In UserDiscoveryEventHandler")
		ude := &UserDiscoveryEvent{}
		err := json.Unmarshal(msg.Data, ude)
		if err != nil {
			logrus.Error(errors.Wrap(err, "cannot unmarshal UserDiscoveryEvent"))
			return
		}
		if ude.UserID == userID {
			channels, err := subscription.GetSubscribedChannels(ude.UserID)
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
			if !channelExist {
				logrus.Warnf("%s is getting subscribed to %s", ude.UserID, ude.ChannelID)
				go subscription.NewSubsciption(ctx, userID, ude.ChannelID, NewMessageEventHandler(ude.ChannelID, ude.UserID, userSockets))
			}
		}
	}
}
func PublishUserDiscoveryEvent(ude UserDiscoveryEvent) error {

	u := &UserDiscoveryEvent{ChannelID: ude.ChannelID, UserID: ude.UserID}
	conn, err := nats.Client(ude.UserID)
	if err != nil {
		return errors.Wrap(err, "cannot connect to nats")
	}
	bs, err := json.Marshal(u)
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
