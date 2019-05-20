package join

import (
	"git.raad.cloud/cloud/hermes/pkg/drivers/nats"
	"git.raad.cloud/cloud/hermes/pkg/eventHandler"
	"git.raad.cloud/cloud/hermes/pkg/session"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"golang.org/x/net/context"
)

//JoinPayload ...
type JoinPayload struct {
	UserID    string
	SessionId string
}

func Handle(ctx context.Context, sig *JoinPayload) {

	s, err := session.GetSession(sig.SessionId)
	if err != nil {
		msg := errors.Wrap(err, "cannot get session").Error()
		logrus.Error(msg)
		logrus.Error(errors.Wrap(err, "error in joining"))
	}
	//logic session validation
	_ = s
	// check jwt
	check := true
	if !check {
		msg := errors.New("jwt is shit")
		logrus.Error(msg.Error())
		logrus.Error(errors.Wrap(err, "error in authenticating"))
	}
	//get user id from jwt
	//ctx, _ = context.WithTimeout(ctx, time.Hour*1)

	logrus.Infof("Subscribing to user-discovery as %s", sig.UserID)
	sub := nats.MakeSubscriber(ctx, sig.UserID, "test-cluster", "0.0.0.0:4222", "user-discovery", eventHandler.UserDiscoveryEventHandler(ctx, sig.UserID, s.SessionID))
	go sub()
}
