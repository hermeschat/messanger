package eventhandlers

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/bson"
	"hermes/db"
	"hermes/subscription/nats"
)

type ReadSignal struct {
	UserID    string
	MessageID string
	ChannelID string
}

func HandleRead(sig *ReadSignal) error {
	bs, err := json.Marshal(sig)
	if err != nil {
		return errors.Wrap(err, "could not marshall read signal")
	}
	err = nats.PublishNewMessage(sig.UserID, sig.ChannelID, bs)
	if err != nil {
		return errors.Wrap(err, "error in publishing eventhandlers read")
	}
	_, err = db.Channels().UpdateOne(context.Background(), bson.M{fmt.Sprintf("%s.messages._id", sig.ChannelID): sig.MessageID}, bson.M{"$set": bson.M{fmt.Sprintf("%s.messages.%s.read", sig.ChannelID, sig.MessageID): true}})
	if err != nil {
		return errors.Wrap(err, "error in updating mongo db to put read")
	}

	return nil
}
