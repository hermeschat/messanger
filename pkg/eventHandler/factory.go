package eventHandler

import (
	"fmt"

	"git.raad.cloud/cloud/hermes/pkg/api"
	"git.raad.cloud/cloud/hermes/pkg/drivers/nats"
	stan "github.com/nats-io/go-nats-streaming"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"golang.org/x/net/context"
)

//UserDiscoveryEventHandler handles user discovery
func UserDiscoveryEventHandler(userID string) func(msg *stan.Msg) {
	return func(msg *stan.Msg) {

		ude := &api.UserDiscoveryEvent{}
		err := ude.XXX_Unmarshal(msg.Data)
		if err != nil {
			logrus.Error(errors.Wrap(err, "cannot unmarshal UserDiscoveryEvent"))
			return
		}
		if ude.UserID == userID {
			ctx, _ := context.WithCancel(context.Background())
			sub := nats.MakeSubscriber(ctx, userID,"test-cluster", "0.0.0.0:4222", ude.ChannelID, NewMessageHandler(ude.ChannelID, ude.UserID))
			go sub()
		}
	}
}

//NewMessageHandler handles the message delivery from nats to user
func NewMessageHandler(channelID string, userID string) func(msg *stan.Msg) {
	return func(msg *stan.Msg) {
		// push kon be user
		fmt.Printf("New Message In %s", channelID)
	}
}

func subscribeChannel(channelID string, userID string) {
	ctx, _ := context.WithCancel(context.Background())
	sub := nats.MakeSubscriber(ctx, userID,"test-cluster", "0.0.0.0:4222", channelID, NewMessageHandler(channelID, userID))
	go sub()
}
