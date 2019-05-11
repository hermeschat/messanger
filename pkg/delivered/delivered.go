package delivered

import (
	"encoding/json"
	"git.raad.cloud/cloud/hermes/pkg/drivers/nats"
	"github.com/pkg/errors"
)

//DeliveredSignal ...
type DeliverdSignal struct {
	UserID string
	ChannelID string
	MessageID string
}

//Handle ...
func Handle(sig *DeliverdSignal) error {
	bs, err := json.Marshal(sig)
	if err != nil {
		return errors.Wrap(err, "error in marshaling signal")
	}
	err = nats.PublishNewMessage("test-cluster",sig.UserID, "0.0.0.0:4222", sig.ChannelID, bs)
	if err != nil {
		return errors.Wrap(err, "error in publishing delivered signal")
	}
	return nil
}
