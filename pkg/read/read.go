package read

import (
	"encoding/json"

	"git.raad.cloud/cloud/hermes/pkg/drivers/nats"
	"github.com/pkg/errors"
)

type ReadSignal struct {
	UserID    string
	MessageID string
	ChannelID string
}

func Handle(sig *ReadSignal) error {
	bs, err := json.Marshal(sig)
	if err != nil {
		return errors.Wrap(err, "error in marshaling signal")
	}
	err = nats.PublishNewMessage("test-cluster", sig.UserID, "0.0.0.0:4222", sig.ChannelID, bs)
	if err != nil {
		return errors.Wrap(err, "error in publishing message read")
	}
	return nil
}
