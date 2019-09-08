package eventhandlers

import (
	"context"
	"sync"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"hermes/api"
	"hermes/pkg/discovery"
	"hermes/pkg/drivers/nats"
)

//JoinPayload ...
type JoinPayload struct {
	UserID    string
	SessionId string
}

//HandleJoin handles join event
func HandleJoin(ctx context.Context, sig *JoinPayload, userSockets *struct {
	sync.RWMutex
	Us map[string]api.Hermes_EventBuffServer
}) {
	channels, err := getSubscribedChannels(sig.UserID)
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
		sub := nats.MakeSubscriber(ctx, sig.UserID, discovery.UserDiscoveryChannel, discovery.UserDiscoveryEventHandler(ctx, sig.UserID, userSockets))
		go sub()
		return
	}
	logrus.Info("Already subscribing to user-discovery")
}
