package delivered

import (
	"encoding/json"
	"github.com/pkg/errors"
	"hermes/pkg/drivers/nats"
)

//DeliveredSignal ...
type DeliverdSignal struct {
	UserID    string
	ChannelID string
	MessageID string
}

//Handle delivered signal will publish a new event of type delivered in given channel to notify other users that the given message is delivered
func Handle(sig *DeliverdSignal) error {
	bs, err := json.Marshal(sig)
	if err != nil {
		return errors.Wrap(err, "error in marshaling signal")
	}
	err = nats.PublishNewMessage("test-cluster", sig.UserID, "0.0.0.0:4222", sig.ChannelID, bs)
	if err != nil {
		return errors.Wrap(err, "error in publishing delivered signal")
	}
	return nil
}
