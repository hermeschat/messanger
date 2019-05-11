package join

import (
	"git.raad.cloud/cloud/hermes/pkg/drivers/nats"
	"git.raad.cloud/cloud/hermes/pkg/eventHandler"
	"git.raad.cloud/cloud/hermes/pkg/repository/session"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"golang.org/x/net/context"
)
//JoinPayload ...
type JoinPayload struct {
	UserID string
	SessionId string
}

func Handle(sig *JoinPayload) error {

	s, err := session.Get(sig.SessionId)
	if err != nil {
		msg := errors.Wrap(err, "cannot get session").Error()
		logrus.Error(msg)
		return errors.Wrap(err, "error in joining")
	}
	//logic session validation
	_ = s
	// check jwt
	check := true
	if !check {
		msg := errors.New("jwt is shit")
		logrus.Error(msg.Error())
		return errors.Wrap(err, "error in authenticating")
	}
	//get user id from jwt
	userID := ""
	ctx, _ := context.WithCancel(context.Background())

	sub := nats.MakeSubscriber(ctx, sig.UserID,"test-cluster", "0.0.0.0:4222", "user-discovery", eventHandler.UserDiscoveryEventHandler(userID))
	go sub()
	return nil
}
