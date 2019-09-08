package nats

import (
	"context"
	"sync"

	"github.com/nats-io/go-nats-streaming/pb"
	"github.com/sirupsen/logrus"

	stan "github.com/nats-io/go-nats-streaming"
	"github.com/pkg/errors"
)

type Config struct {
	NatsSrvAddr string
	ClusterId   string
}

var State = struct {
	Mu sync.RWMutex
	Ss map[string]*stan.Conn
}{sync.RWMutex{}, map[string]*stan.Conn{}}

//NatsClient manage to keep nats connection open
func NatsClient(clusterID string, natsSrvAddr string, clientID string) (*stan.Conn, error) {
	State.Mu.Lock()
	conn, ok := State.Ss[clientID]
	if ok || conn != nil {
		State.Mu.Unlock()
		return conn, nil

	}
	logrus.Warnf("Connection was unusable: %v", conn)
	delete(State.Ss, clientID)
	logrus.Warnf("Trying to get a connection on behalf of %s", clientID)

	natsClient, err := stan.Connect(clusterID, clientID, stan.NatsURL(natsSrvAddr))
	//,stan.SetConnectionLostHandler(func(_ stan.Conn, reason error) {
	//		log.Fatalf("Connection lost, reason: %v", reason)
	//	})
	if err != nil {
		logrus.Warnf("Ok is %v conn is %v State is %+v", ok, conn, State.Ss)
		return nil, errors.Wrapf(err, "Can't connect %v: %v", clientID, err)
	}
	State.Ss[clientID] = &natsClient
	State.Mu.Unlock()
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
		return errors.Wrap(err, "failed to publish eventhandlers")
	}
	return nil
}

//t is type we need to pass to find our eventhandlers type
func MakeSubscriber(ctx context.Context, userID string, clusterID string, natsSrvAddr string, ChannelId string, handler func(msg *stan.Msg)) Subscriber {
	return func() {
		natscon, err := NatsClient(clusterID, natsSrvAddr, userID)

		if err != nil {
			logrus.Error(errors.Wrapf(err, "Can't connect: %v.\nMake sure a NATS Streaming Server is running at: %s", err, natsSrvAddr))
			return
		}
		logrus.Info("Connected to %s clusterID: [%s] clientID: [%s]\n", natsSrvAddr, clusterID, userID)

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
