package read

import (
	"fmt"
	"git.raad.cloud/cloud/hermes/pkg/api"
	"git.raad.cloud/cloud/hermes/pkg/drivers/nats"
	"github.com/pkg/errors"
)

type ReadSignal struct {
	MessageID string
	ChannelID string
}

func Handle(sig *ReadSignal) error {
	err := nats.PublishNewMessage("test-cluster", "0.0.0.0:4222", sig.ChannelID, &api.Message{
		MessageType: "3",
		Body:        fmt.Sprintf(`{"message_id":%s}`, sig.MessageID),
	})
	if err != nil {
		return errors.Wrap(err, "error in publishing message read")
	}
	return nil
}
