package nats

import (
	"context"
	"github.com/nats-io/go-nats-streaming/pb"
	"github.com/sirupsen/logrus"
	"sync"

	stan "github.com/nats-io/go-nats-streaming"
	"github.com/pkg/errors"
)

type Config struct {
	NatsSrvAddr string
	ClusterId   string
}

var state = struct {
	mu sync.Mutex
	Ss map[string]*stan.Conn
}{sync.Mutex{}, map[string]*stan.Conn{}}

//NatsClient manage to keep nats connection open
func NatsClient(clusterID string, natsSrvAddr string, clientID string) (*stan.Conn, error) {
	state.mu.Lock()
	conn, ok := state.Ss[clientID]
	if ok || conn != nil {
		state.mu.Unlock()
		return conn, nil

	}
	logrus.Warnf("Connection was unusable: %v", conn)
	delete(state.Ss, clientID)
	logrus.Warnf("Trying to get a connection on behalf of %s", clientID)

	natsClient, err := stan.Connect(clusterID, clientID, stan.NatsURL(natsSrvAddr))
	//,stan.SetConnectionLostHandler(func(_ stan.Conn, reason error) {
	//		log.Fatalf("Connection lost, reason: %v", reason)
	//	})
	if err != nil {
		logrus.Warnf("Ok is %v conn is %v State is %+v", ok, conn, state.Ss)
		return nil, errors.Wrapf(err, "Can't connect %v: %v", clientID, err)
	}
	state.Ss[clientID] = &natsClient
	state.mu.Unlock()
	return &natsClient, nil
}

type Subscriber func()

func PublishNewMessage(clusterID string, userID, natsSrvAddr string, ChannelId string, bs []byte) error {
	// Connect to NATS-Streaming

	natscon, err := NatsClient(clusterID, natsSrvAddr, userID)
	if err != nil {
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
		natscon, err := NatsClient(clusterID, natsSrvAddr, userID)

		if err != nil {
			logrus.Error(errors.Wrapf(err, "Can't connect: %v.\nMake sure a NATS Streaming Server is running at: %s", err, natsSrvAddr))
			return
		}
		logrus.Info("Connected to %s clusterID: [%s] clientID: [%s]\n", natsSrvAddr, clusterID, userID)
		//defer sc.Close()
		// sample handler
		//i := 0
		//mcb := func(msg *stan.Msg) {
		//	i++
		//	fmt.Println(msg, i)
		//}

		//startOpt := stan.StartAt(pb.StartPosition_NewOnly)
		//without := stan.StartWithLastReceived()
		//wait := stan.AckWait(time.Second * 1)
		// sub, err := (*natscon).Subscribe(ChannelId, handler, startOpt, stan.DurableName(userID))
		sub, err := (*natscon).Subscribe(ChannelId, handler, stan.StartAt(pb.StartPosition_NewOnly))
		if err != nil {
			logrus.Error(err)
			return
		}
		_ = sub

		logrus.Infof("Listening on [%s], clientID=[%s], qgroup=[%s] durable=[%s]\n", ChannelId, userID, "qgroup", userID)
		durable := userID
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
