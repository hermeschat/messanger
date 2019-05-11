package nats

import (
	"context"
	"fmt"
	"log"

	"github.com/sirupsen/logrus"

	stan "github.com/nats-io/go-nats-streaming"
	"github.com/pkg/errors"
	uuid "github.com/satori/go.uuid"
)

type Config struct {
	NatsSrvAddr string
	ClusterId   string
}

func NatsClient(clusterID string, natsSrvAddr string) (stan.Conn, error) {
	id, err := uuid.NewV4()
	if err != nil {
		return nil, errors.Wrap(err, "Can't generate UUID?!")
	}
	natsClient, err := stan.Connect(clusterID, id.String(), stan.NatsURL(natsSrvAddr))
	if err != nil {
		return nil, errors.Wrapf(err, "Can't connect: %v.\nMake sure a NATS Streaming Server is running at: %s", err, natsSrvAddr)
	}
	return natsClient, nil
}

type Subscriber func() error

func PublishNewMessage(clusterID string, userID, natsSrvAddr string, ChannelId string, bs []byte) error {
	// Connect to NATS-Streaming
	natsClient, err := stan.Connect(clusterID, userID, stan.NatsURL(natsSrvAddr))
	if err != nil {
		return errors.Wrapf(err, "Can't connect: %v.\nMake sure a NATS Streaming Server is running at: %s", err, natsSrvAddr)
	}
	defer natsClient.Close()

	if err := natsClient.Publish(ChannelId, bs); err != nil {
		return errors.Wrap(err, "failed to publish message")
	}
	return nil
}

//t is type we need to pass to find our message type
func MakeSubscriber(ctx context.Context, userID string, clusterID string, natsSrvAddr string, ChannelId string, handler func(msg *stan.Msg)) Subscriber {
	return func() error {
		durable := ""
		sc, err := stan.Connect(clusterID, userID, stan.NatsURL(natsSrvAddr),
			stan.SetConnectionLostHandler(func(_ stan.Conn, reason error) {
				log.Fatalf("Connection lost, reason: %v", reason)
			}))
		if err != nil {
			return errors.Wrapf(err, "Can't connect: %v.\nMake sure a NATS Streaming Server is running at: %s", err, natsSrvAddr)
		}
		logrus.Info("Connected to %s clusterID: [%s] clientID: [%s]\n", natsSrvAddr, clusterID, userID)
		defer sc.Close()

		// sample handler
		//i := 0
		//mcb := func(msg *stan.Msg) {
		//	i++
		//	fmt.Println(msg, i)
		//}

		startOpt := stan.DeliverAllAvailable()

		sub, err := sc.QueueSubscribe(ChannelId, "qgroup", handler, startOpt, stan.DurableName("durable"))
		if err != nil {
			sc.Close()
			log.Fatal(err)
		}

		logrus.Infof("Listening on [%s], clientID=[%s], qgroup=[%s] durable=[%s]\n", ChannelId, userID, "qgroup", "durable")

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
}

////Subscribe used when a reviever wants to get messages.
//func Subscribe(ctx context.Context, clusterID string, natsSrvAddr string, msg *api.InstantMessage) error {
//	durable := ""
//	id, err := uuid.NewV4()
//	sc, err := stan.Connect(clusterID, id.String(), stan.NatsURL(natsSrvAddr),
//		stan.SetConnectionLostHandler(func(_ stan.Conn, reason error) {
//			log.Fatalf("Connection lost, reason: %v", reason)
//		}))
//	if err != nil {
//		return errors.Wrapf(err, "Can't connect: %v.\nMake sure a NATS Streaming Server is running at: %s", err, natsSrvAddr)
//	}
//	logrus.Info("Connected to %s clusterID: [%s] clientID: [%s]\n", natsSrvAddr, clusterID, id.String())
//	defer sc.Close()
//
//	i := 0
//	mcb := func(msg *stan.Msg) {
//		i++
//		fmt.Println(msg, i)
//	}
//
//	startOpt := stan.DeliverAllAvailable()
//
//	sub, err := sc.QueueSubscribe(msg.Channel, "qgroup", mcb, startOpt, stan.DurableName("durable"))
//	if err != nil {
//		sc.Close()
//		log.Fatal(err)
//	}
//
//	logrus.Infof("Listening on [%s], clientID=[%s], qgroup=[%s] durable=[%s]\n", msg.Channel, id.String(), "qgroup", "durable")
//
//	// Wait for a SIGINT (perhaps triggered by user with CTRL-C)
//	// Run cleanup when signal is received
//	go func() {
//		select {
//		case <-ctx.Done():
//
//			err := ctx.Err()
//			if err != nil {
//				fmt.Println("in <-ctx.Done(): ", err)
//			}
//			logrus.Infof("\n unsubscribing and closing connection...\n\n")
//			// Do not unsubscribe a durable on exit, except if asked to.
//			if durable == "" {
//				sub.Unsubscribe()
//			}
//			sc.Close()
//			break
//		}
//
//	}()
//	return nil
//}

// func Run() {
// 	timeout := 10 * time.Second
// 	timeOutContext, _ := context.WithTimeout(context.Background(), timeout)
// 	PublishNewMessage("test-cluster", "0.0.0.0:4222", "ChannelMan",&api.InstantMessage{
// 		MessageType: "1",
// 		Channel:     "ChannelMan",
// 	})
// 	go Subscribe(timeOutContext, "test-cluster", "0.0.0.0:4222", &api.InstantMessage{
// 		MessageType: "1",
// 		Channel:     "ChannelMan",
// 	})

// 	// Wait for the timeout to expire
// 	<-timeOutContext.Done()

// }

// func printMsg(m *stan.Msg, i int) {
// 	log.Printf("[#%d] Received on [%s]: '%s'\n", i, m.Subject, m)
// }
