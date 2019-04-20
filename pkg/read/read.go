package read

import (
	"encoding/json"
	"fmt"
	"git.raad.cloud/cloud/hermes/pkg/api"
	"git.raad.cloud/cloud/hermes/pkg/drivers/nats"
	"github.com/mitchellh/mapstructure"
	"github.com/pkg/errors"
)

type ReadSignal struct {
	MessageID string
	ChannelID string
}


func Handle(sig *api.Signal) *api.Response {
	ds := &ReadSignal{}
	var j map[string]interface{}
	err := json.Unmarshal([]byte(sig.Payload), &j)
	if err != nil {
		msg := errors.Wrap(err, "error in unmarshalling payload")
		return &api.Response{
			Error : msg.Error(),
		}
	}
	err = mapstructure.Decode(j, ds)
	if err != nil {
		msg := errors.Wrap(err, "Error in unmarshaling signal pa")
		return &api.Response{
			Error: msg.Error(),
		}
	}
	nats.PublishNewMessage("test-cluster", "0.0.0.0:4222", ds.ChannelID, &api.InstantMessage{
		MessageType: "3",
		Body:        fmt.Sprintf(`{"message_id":%s}`, ds.MessageID),
	})
	return &api.Response{}
}
