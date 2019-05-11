package read

import (
	"encoding/json"
	stan "github.com/nats-io/go-nats-streaming"
	"github.com/pkg/errors"
)

type ReadSignal struct {
	UserID string
	MessageID string
	ChannelID string
}

func Handle(sig *ReadSignal) error {
	err := publishNewMessage("test-cluster",sig.UserID, "0.0.0.0:4222", sig.ChannelID, &ReadSignal{
		MessageID: sig.MessageID,
		ChannelID:        sig.ChannelID,
	})
	if err != nil {
		return errors.Wrap(err, "error in publishing message read")
	}
	return nil
}

//publishNewMessage is send function. Every message should be published to a channel to
//be delivered to subscribers. In streaming, published Message is persistant.
func publishNewMessage(clusterID string,userID, natsSrvAddr string, ChannelId string, msg *ReadSignal) error {
	// Connect to NATS-Streaming
	natsClient, err := stan.Connect(clusterID, userID, stan.NatsURL(natsSrvAddr))
	if err != nil {
		return errors.Wrapf(err, "Can't connect: %v.\nMake sure a NATS Streaming Server is running at: %s", err, natsSrvAddr)
	}
	defer natsClient.Close()
	bs, err := json.Marshal(msg)
	if err != nil {
		return errors.Wrap(err, "error in marshaling message")
	}
	if err := natsClient.Publish(ChannelId, bs); err != nil {
		return errors.Wrap(err, "failed to publish message")
	}
	return nil
}