package read

import (
	"encoding/json"

	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/bson"
	"hermes/pkg/db/mongo"
	"hermes/pkg/drivers/nats"
)

type ReadSignal struct {
	UserID    string
	MessageID string
	ChannelID string
}

func Handle(sig *ReadSignal) error {
	bs, err := json.Marshal(sig)
	if err != nil {
		return errors.Wrap(err, "could not marshall read signal")
	}
	err = nats.PublishNewMessage("test-cluster", sig.UserID, "0.0.0.0:4222", sig.ChannelID, bs)
	if err != nil {
		return errors.Wrap(err, "error in publishing message read")
	}
	err = mongo.UpdateAllMatched("messages", bson.M{"message_id": sig.MessageID}, bson.M{"$set": bson.M{"read": true}})
	if err != nil {
		return errors.Wrap(err, "error in updating mongo db to put read")
	}
	return nil
}
