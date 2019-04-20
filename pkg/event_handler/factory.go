package event_handler

import (
	"fmt"

	"git.raad.cloud/cloud/hermes/pkg/api"

	"git.raad.cloud/cloud/hermes/pkg/drivers/nats"
	stan "github.com/nats-io/go-nats-streaming"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"golang.org/x/net/context"
)

func UserDiscoveryEventHandler(userId string) func(msg *stan.Msg) {
	return func(msg *stan.Msg) {

		ude := &api.UserDiscoveryEvent{}
		err := ude.XXX_Unmarshal(msg.Data)
		if err != nil {
			logrus.Error(errors.Wrap(err, "cannot unmarshal UserDiscoveryEvent"))
			return
		}
		if ude.UserID == userId {
			ctx, _ := context.WithCancel(context.Background())
			sub := nats.MakeSubscriber(ctx, "test-cluster", "0.0.0.0:4222", ude.ChannelID, NewMessageHandler(ude.ChannelID))
			go sub()
		}
	}
}

func NewMessageHandler(channelId string) func(msg *stan.Msg) {
	return func(msg *stan.Msg) {
		fmt.Printf("New Message In %s", channelId)
	}
}
