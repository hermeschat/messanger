package nats

import (
	"context"
	"github.com/nats-io/go-nats-streaming/pb"
	"log"

	"github.com/sirupsen/logrus"

	"github.com/nats-io/go-nats-streaming"
	"github.com/pkg/errors"
	"git.raad.cloud/cloud/hermes/pkg/session"
)

type Config struct {
	NatsSrvAddr string
	ClusterId   string
}
//NatsClient manage to keep nats connection open
func NatsClient(clusterID string, natsSrvAddr string, clientID string,) (*stan.Conn, error) {

	session.State.Lock()
	conn, ok := (session.State).Ss[clientID]
	session.State.Unlock()
	if ok {
		logrus.Warn("e khali nist")
		logrus.Warn(clientID)
		return conn, nil
	}
	logrus.Warn("__________")
	for k, _:= range(session.State).Ss{
		logrus.Warn(k)
	}
	logrus.Warn("-------------------")
	logrus.Warn("e khalie!!!!!!!!!!!!!!!")
	logrus.Warn(clientID)
	natsClient, err := stan.Connect(clusterID, clientID, stan.NatsURL(natsSrvAddr) )
	//,stan.SetConnectionLostHandler(func(_ stan.Conn, reason error) {
	//		log.Fatalf("Connection lost, reason: %v", reason)
	//	})
	if err != nil {
		return nil,errors.Wrapf(err, "Can't connect: %v.\nMake sure a NATS Streaming Server is running at: %s", err, natsSrvAddr)
	}
	session.State.Lock()
	logrus.Warn("0000000000")
	for k, _:= range(session.State).Ss{
		logrus.Warn(k)
	}
	logrus.Warn("0000000000")
	(session.State).Ss[clientID] = &natsClient
	logrus.Warn("11111111111")
	for k, _:= range(session.State).Ss{
		logrus.Warn(k)
	}
	logrus.Warn("111111111111")
	session.State.Unlock()
	return &natsClient, nil
}

type Subscriber func()

func PublishNewMessage(clusterID string, userID, natsSrvAddr string, ChannelId string, bs []byte) error {
	// Connect to NATS-Streaming

	natscon, err := NatsClient(clusterID, natsSrvAddr, userID)
	if err != nil{
		return errors.Wrap(err, "failed to connect to nats")
	}
	if err := (*natscon).Publish(ChannelId, bs); err != nil {
		return errors.Wrap(err, "failed to publish message")
	}
	return nil
}

//t is type we need to pass to find our message type
func MakeSubscriber(ctx context.Context, userID string, clusterID string, natsSrvAddr string, ChannelId string, handler func(msg *stan.Msg)) Subscriber {
	return func() {
		durable := "durable"
		natscon, err := NatsClient(clusterID, natsSrvAddr, userID)

		if err != nil {
			logrus.Error(errors.Wrapf(err, "Can't connect: %v.\nMake sure a NATS Streaming Server is running at: %s", err, natsSrvAddr))
		}
		logrus.Info("Connected to %s clusterID: [%s] clientID: [%s]\n", natsSrvAddr, clusterID, userID)
		//defer sc.Close()
		// sample handler
		//i := 0
		//mcb := func(msg *stan.Msg) {
		//	i++
		//	fmt.Println(msg, i)
		//}

		startOpt := stan.DeliverAllAvailable()
		startOpt = stan.StartAt(pb.StartPosition_NewOnly)
		sub, err := (*natscon).Subscribe(ChannelId, handler, startOpt, stan.DurableName(userID))
		if err != nil {
			log.Fatal(err)
		}
		_ = sub

		logrus.Infof("Listening on [%s], clientID=[%s], qgroup=[%s] durable=[%s]\n", ChannelId, userID, "qgroup", userID)

		// Wait for a SIGINT (perhaps triggered by user with CTRL-C)
		// Run cleanup when signal is received
		go func() {
			select {
			case <-ctx.Done():

				err := ctx.Err()
				if err != nil {
					logrus.Infof("Closing Subscription <-ctx.Done(): ", err)
				}
				logrus.Infof("\n unsubscribing and closing connection...\n\n")
				// Do not unsubscribe a durable on exit, except if asked to.
				if durable == "" {
					sub.Unsubscribe()
				}
				(*natscon).Close()
				break
			}

		}()
		//return nil
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
