package nats

import (
	"context"
	"fmt"
	"sync"

	"github.com/amirrezaask/config"
	"github.com/nats-io/go-nats-streaming/pb"
	"github.com/sirupsen/logrus"

	stan "github.com/nats-io/go-nats-streaming"
	"github.com/pkg/errors"
)

type Config struct {
	NatsSrvAddr string
	ClusterId   string
}

type natsConnections struct {
	sync.RWMutex
	conns map[string]*stan.Conn
}

type Handler = stan.MsgHandler

func (n *natsConnections) CloseConnection(userID string) error {
	n.Lock()
	defer n.Unlock()
	natsCon, exists := n.conns[userID]
	if !exists {
		return errors.New(fmt.Sprintf("nats connection for user %s not found", userID))
	}
	err := (*natsCon).Close()
	if err != nil {
		return errors.Wrapf(err, "could not close nats connection of user %s due to %s", userID, err)
	}
	delete(n.conns, userID)
	return nil
}

var Connections = &natsConnections{
	sync.RWMutex{}, map[string]*stan.Conn{},
}

func Client(clientID string) (*stan.Conn, error) {
	Connections.Lock()
	defer Connections.Unlock()
	conn, ok := Connections.conns[clientID]
	if ok || conn != nil {
		return conn, nil

	}
	logrus.Warnf("Connection was unusable or does not exist for %s\n", clientID)
	delete(Connections.conns, clientID)
	logrus.Infof("Trying to get a connection on behalf of %s\n", clientID)

	//natsClient, err := stan.Connect(config.Get("cluster_id"), clientID, stan.NatsURL(string(config.Get("nats_uri"))))
	natsClient, err := stan.Connect("test-cluster", clientID, stan.NatsURL("0.0.0.0:4222"))
	stan.SetConnectionLostHandler(func(_ stan.Conn, reason error) { logrus.Errorf("Connection lost, reason: %v", reason) })
	if err != nil {
		return nil, errors.Wrapf(err, "Can't connect %v: %v", clientID, err)
	}
	Connections.conns[clientID] = &natsClient
	return &natsClient, nil
}

type Subscriber func()

func PublishNewMessage(userID, ChannelID string, bs []byte) error {
	natscon, err := Client(userID)
	if err != nil {
		return errors.Wrap(err, "failed to connect to nats")
	}
	if err := (*natscon).Publish(ChannelID, bs); err != nil {
		return errors.Wrap(err, "failed to publish eventhandlers")
	}
	return nil
}

//t is type we need to pass to find our eventhandlers type
func MakeSubscriber(ctx context.Context, userID string, ChannelId string, handler Handler) Subscriber {
	return func() {
		natscon, err := Client(userID)

		if err != nil {
			logrus.Error(errors.Wrapf(err, "Can't connect: %v.\nMake sure a NATS Streaming Server is running at: %s", err, config.Get("nats_cluster_id")))
			return
		}
		logrus.Info("Connected to %s clusterID: [%s] clientID: [%s]\n", config.Get("nats_uri"), config.Get("nats_cluster_id"), userID)
		sub, err := (*natscon).Subscribe(ChannelId, handler, stan.StartAt(pb.StartPosition_LastReceived))
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
