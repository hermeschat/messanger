package delivered

import (
	"encoding/json"
	stan "github.com/nats-io/go-nats-streaming"
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
	err := publishNewMessage("test-cluster",sig.UserID, "0.0.0.0:4222", sig.ChannelID, sig)
	if err != nil {
		return errors.Wrap(err, "error in publishing delivered signal")
	}
	return nil
}
//publishNewMessage is send function. Every message should be published to a channel to
//be delivered to subscribers. In streaming, published Message is persistant.
func publishNewMessage(clusterID string,userID, natsSrvAddr string, ChannelId string, msg *DeliverdSignal) error {
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