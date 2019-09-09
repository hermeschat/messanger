package eventhandlers

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

//HandleDeliver delivered signal will publish a new event of type delivered in given channel to notify other users that the given eventhandlers is delivered
func HandleDeliver(sig *DeliverdSignal) error {
	bs, err := json.Marshal(sig)
	if err != nil {
		return errors.Wrap(err, "error in marshaling signal")
	}
	err = nats.PublishNewMessage(sig.UserID, sig.ChannelID, bs)
	if err != nil {
		return errors.Wrap(err, "error in publishing delivered signal")
	}
	return nil
}
