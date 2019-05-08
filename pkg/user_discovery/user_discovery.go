package user_discovery

import (
	"git.raad.cloud/cloud/hermes/pkg/api"
	"git.raad.cloud/cloud/hermes/pkg/drivers/nats"
	"git.raad.cloud/cloud/hermes/pkg/repository"
	"github.com/gogo/protobuf/proto"
	"github.com/pkg/errors"
)

func PublishEvent(ude repository.UserDiscoveryEvent) error {
	u := &api.UserDiscoveryEvent{ChannelID: ude.ChannelID, UserID: ude.UserID}
	conn, err := nats.NatsClient("test-cluster", "0.0.0.0:4222")
	if err != nil {
		return errors.Wrap(err, "cannot connect to nats")
	}
	bs, err := proto.Marshal(u)
	if err != nil {
		return errors.Wrap(err, "cannot marshal UserDiscoveryEvent")
	}
	err = conn.Publish("user-discovery", bs)
	if err != nil {
		return errors.Wrap(err, "cannot publish UserDiscoveryEvent")
	}
	return nil
}
