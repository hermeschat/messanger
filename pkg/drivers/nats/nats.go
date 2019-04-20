package nats

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/sirupsen/logrus"

	"git.raad.cloud/cloud/hermes/pkg/api"
	"github.com/gogo/protobuf/proto"
	stan "github.com/nats-io/go-nats-streaming"
	"github.com/pkg/errors"
	uuid "github.com/satori/go.uuid"
)


type Config struct {
	NatsSrvAddr string
	ClusterId string
}

//Publish is send function. Every message should be published to a channel to
//be delivered to subscribers. In streaming, published Message is persistant.
func Publish(clusterID string, natsSrvAddr string, ChannelId string,msg *api.InstantMessage) error {
	// Connect to NATS-Streaming
	id, err := uuid.NewV4()
	if err != nil {
		return errors.Wrap(err, "Can't generate UUID?!")
	}
	natsClient, err := stan.Connect(clusterID, id.String(), stan.NatsURL(natsSrvAddr))
	if err != nil {
		return errors.Wrapf(err, "Can't connect: %v.\nMake sure a NATS Streaming Server is running at: %s", err, natsSrvAddr)
	}
	defer natsClient.Close()

	bs, err := proto.Marshal(msg)
	if err != nil {
		return errors.Wrap(err, "failed to marshal proto message")
	}

	if err := natsClient.Publish(ChannelId, bs); err != nil {
		return errors.Wrap(err, "failed to publish message")
	}
	return nil
}

//Subscribe used when a reviever wants to get messages.
func Subscribe(ctx context.Context, clusterID string, natsSrvAddr string, msg *api.InstantMessage) error {
	durable := ""
	id, err := uuid.NewV4()
	sc, err := stan.Connect(clusterID, id.String(), stan.NatsURL(natsSrvAddr),
		stan.SetConnectionLostHandler(func(_ stan.Conn, reason error) {
			log.Fatalf("Connection lost, reason: %v", reason)
		}))
	if err != nil {
		return errors.Wrapf(err, "Can't connect: %v.\nMake sure a NATS Streaming Server is running at: %s", err, natsSrvAddr)
	}
	logrus.Info("Connected to %s clusterID: [%s] clientID: [%s]\n", natsSrvAddr, clusterID, id.String())
	defer sc.Close()

	i := 0
	mcb := func(msg *stan.Msg) {
		i++
		printMsg(msg, i)
	}

	startOpt := stan.DeliverAllAvailable()

	sub, err := sc.QueueSubscribe(msg.Channel, "qgroup", mcb, startOpt, stan.DurableName("durable"))
	if err != nil {
		sc.Close()
		log.Fatal(err)
	}

	logrus.Infof("Listening on [%s], clientID=[%s], qgroup=[%s] durable=[%s]\n", msg.Channel, id.String(), "qgroup", "durable")

	// Wait for a SIGINT (perhaps triggered by user with CTRL-C)
	// Run cleanup when signal is received
	go func() {
		select {
		case <-ctx.Done():

			err := ctx.Err()
			if err != nil {
				fmt.Println("in <-ctx.Done(): ", err)
			}
			logrus.Infof("\n unsubscribing and closing connection...\n\n")
			// Do not unsubscribe a durable on exit, except if asked to.
			if durable == "" {
				sub.Unsubscribe()
			}
			sc.Close()
			break
		}

	}()
	return nil
}

func Run() {
	timeout := 10 * time.Second
	timeOutContext, _ := context.WithTimeout(context.Background(), timeout)
	Publish("test-cluster", "0.0.0.0:4222", "ChannelMan",&api.InstantMessage{
		MessageType: "1",
		Channel:     "ChannelMan",
	})
	go Subscribe(timeOutContext, "test-cluster", "0.0.0.0:4222", &api.InstantMessage{
		MessageType: "1",
		Channel:     "ChannelMan",
	})

	// Wait for the timeout to expire
	<-timeOutContext.Done()

}

func printMsg(m *stan.Msg, i int) {
	log.Printf("[#%d] Received on [%s]: '%s'\n", i, m.Subject, m)
}
