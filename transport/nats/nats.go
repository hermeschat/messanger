package nats

import "C"
import (
	"fmt"
	"github.com/hermeschat/engine/config"
	stan "github.com/nats-io/stan.go"
)
type natsClient struct {
	stan.Conn
}

func Client(clusterID, clientID, natsUri string) (stan.Conn, error) {
	return stan.Connect(clusterID, clientID, stan.NatsURL(natsUri))

}
func (sc *natsClient) StreamSubscriber(msg stan.Msg) stan.Subscription {
	sub, _ := sc.Subscribe("foo", func(m *stan.Msg) {
		fmt.Printf("Received a message: %s\n", string(m.Data))
	})
	return sub
}

func (sc *natsClient) StreamPublisher(subject string, msg []byte) error {
	return sc.Publish(subject, msg)
}

func StreamUnsubscriber(sub stan.Subscription) error {
	return sub.Unsubscribe()
}
func HealthCheck() error {
	_, err := Client(config.C.GetString("nats.cluser"), config.C.GetString("nats.client"), config.NatsUrl())
	if err != nil {
		return err
	}
	return nil
}