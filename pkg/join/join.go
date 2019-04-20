package join

import (
	"encoding/json"
	"git.raad.cloud/cloud/hermes/pkg/api"
	"git.raad.cloud/cloud/hermes/pkg/drivers/nats"
	"git.raad.cloud/cloud/hermes/pkg/new_message"
	"git.raad.cloud/cloud/hermes/pkg/session"
	stan "github.com/nats-io/go-nats-streaming"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"golang.org/x/net/context"
)




type JoinPayload struct {
	SessionId string
}
func Handle(sig *api.Signal) *api.Response {
	payload := &JoinPayload{}
	err := json.Unmarshal([]byte(sig.Payload), payload)
	if err != nil {
		msg:= errors.Wrap(err, "cannot unmarshal payload").Error()
		logrus.Error(msg)
		return &api.Response{
			Code: "500",
			Error: msg,
		}
	}
	if payload.SessionId == "" {
		msg:= "SessionId not exists"
		logrus.Info(msg)
		// create new session
	}

	s, err := session.Get(payload.SessionId)
	if err != nil {
		msg:= errors.Wrap(err, "cannot get session").Error()
		logrus.Error(msg)
		return &api.Response{
			Code: "500",
			Error: msg,
		}
	}
	//logic session validation
	_ = s
	// check jwt
	check := true
	if !check {
		msg := errors.New( "jwt is shit")
		logrus.Error(msg.Error())
		return &api.Response{
			Code: "500",
			Error: msg.Error(),
		}
	}
	//get user id from jwt
	userID := ""
	ctx, _ := context.WithCancel(context.Background())

	sub := nats.SubscriberFactory(ctx, "test-cluster", "0.0.0.0:4222", "user-discovery", UserDiscoveryEventHandlerFactory(userID))
	go sub()
	return &api.Response{}
}


func UserDiscoveryEventHandlerFactory(userId string) func (msg *stan.Msg) {
	return func(msg *stan.Msg) {

		ude := &api.UserDiscoveryEvent{}
		err := ude.XXX_Unmarshal(msg.Data)
		if err != nil {
			logrus.Error(errors.Wrap(err, "cannot unmarshal UserDiscoveryEvent"))
			return
		}
		if ude.UserID == userId {
			ctx, _ := context.WithCancel(context.Background())
			sub := nats.SubscriberFactory(ctx, "test-cluster", "0.0.0.0:4222", ude.ChannelID, new_message.NewMessageHandlerFactory(ude.ChannelID))
			go sub()
		}
	}
}
