package eventHandler

import (
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
func UserDiscoveryEventHandler(ctx context.Context, userID string, currentSession string) func(msg *stan.Msg) {
	return func(msg *stan.Msg) {
		logrus.Info("!!!!!!!!!!!!!!!!discovery event handler called ")
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

		logrus.Warnf("%s is now subscribed to %s", ude.UserID, ude.ChannelID)
		sub := nats.MakeSubscriber(ctx, userID, "test-cluster", "0.0.0.0:4222", ude.ChannelID, NewMessageEventHandler(ude.ChannelID, ude.UserID))
		sub()

	}
}

var UserSockets = struct {
	sync.RWMutex
	Us map[string]api.Hermes_EventBuffServer
}{
	sync.RWMutex{},
	map[string]api.Hermes_EventBuffServer{},
}

//NewMessageEventHandler handles the message delivery from nats to user
func NewMessageEventHandler(channelID string, userID string) func(msg *stan.Msg) {
	return func(msg *stan.Msg) {
		logrus.Warnf("Message is %v", string(msg.Data))
		//m := &api.Message{}
		//err := json.Unmarshal(msg.Data, m)
		////_ ,err := m.XXX_Marshal(msg.Data, false)
		//if err != nil {
		//	logrus.Errorf("error in unmarshalling nats message in message handler")
		//}
		logrus.Info("In NewMessage Event Handler")
		logrus.Infof("Recieved a new message in %s", channelID)
		c, ok := BaseHub.ClientsMap[userID]
		if !ok {
			logrus.Error("no active connection found for user")
			return
		}

		c.send <- msg.Data
		//UserSockets.RLock()
		//userSocket, ok := UserSockets.Us[userID]
		//if !ok {
		//	logrus.Errorf("error: user socket not found ")
		//	return
		//}
		//err=userSocket.Send(&api.Event{Event:&api.Event_NewMessage{m}})
		//if err != nil {
		//	logrus.Errorf("error: cannot send event new message to user ")
		//	return
		//}
		//UserSockets.RUnlock()
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

	//s, err := session.GetSession(sig.SessionId)
	//if err != nil {
	//	msg := errors.Wrap(err, "cannot get session").Error()
	//	logrus.Error(msg)
	//	logrus.Error(errors.Wrap(err, "error in joining"))
	//}
	//logic session validation
	//_ = s
	// check jwt
	//check := true
	//if !check {
	//	msg := errors.New("jwt is shit")
	//	logrus.Error(msg.Error())
	//	logrus.Error(errors.Wrap(err, "error in authenticating"))
	//}
	//get user id from jwt
	//ctx, _ = context.WithTimeout(ctx, time.Hour*1)

	logrus.Infof("Subscribing to user-discovery as %s", sig.UserID)
	sub := nats.MakeSubscriber(ctx, sig.UserID, "test-cluster", "0.0.0.0:4222", "user-discovery", UserDiscoveryEventHandler(ctx, sig.UserID, ""))
	go sub()

}
