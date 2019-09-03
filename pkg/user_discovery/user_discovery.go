package user_discovery

import (
	"fmt"

	"github.com/gogo/protobuf/proto"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"hermes/api/pb"
	"hermes/pkg/drivers/nats"
)

func PublishEvent(ude UserDiscoveryEvent) error {

	u := &pb.UserDiscoveryEvent{ChannelID: ude.ChannelID, UserID: ude.UserID}
	fmt.Println("client id is ", ude.UserID)
	conn, err := nats.NatsClient("test-cluster", "0.0.0.0:4222", ude.UserID)
	if err != nil {
		return errors.Wrap(err, "cannot connect to nats")
	}
	bs, err := proto.Marshal(u)
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

//UserDiscoveryEvent is the message we send to discovery channel to tell a user
//to subscribe to a certain channel in async way
type UserDiscoveryEvent struct {
	ChannelID string
	UserID    string
}
