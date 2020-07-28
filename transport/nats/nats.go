package nats

import (
	"fmt"
	"github.com/hermeschat/engine/config"
	stan "github.com/nats-io/stan.go"
	"log"
)

func StreamClient(natsUrl string) (stan.Conn, error) {
	clusterID, clientID := config.StanURI()
	return stan.Connect(clusterID, clientID, stan.NatsURL(natsUrl))

}
func StreamSubscriber(sc stan.Conn, msg stan.Msg) stan.Subscription {
	sub, _ := sc.Subscribe("foo", func(m *stan.Msg) {
		fmt.Printf("Received a message: %s\n", string(m.Data))
	})
	return sub
}

func StreamPublisher(sc stan.Conn, subject string, msg []byte) error {
	return sc.Publish(subject, msg)
}

func StreamUnsubscriber(sub stan.Subscription) error {
	return sub.Unsubscribe()
}
func HealthCheck() bool {
	clusterID, clientID := config.StanURI()

	sc, err := stan.Connect(clusterID, clientID, stan.NatsConn(nc),
		stan.SetConnectionLostHandler(func(_ stan.Conn, reason error) {
			log.Fatalf("Connection lost, reason: %v", reason)
		}))
	if err != nil {
		log.Fatalf("Can't connect: %v.\nMake sure a NATS Streaming Server is running at: %s", err, URL)
	}
	log.Printf("Connected to %s clusterID: [%s] clientID: [%s]\n", URL, clusterID, clientID)

}
