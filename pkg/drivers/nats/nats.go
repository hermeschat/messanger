package nats

import stan "github.com/nats-io/go-nats-streaming"

sc, _ := stan.Connect(clusterID, clientID)

// Simple Synchronous Publisher
sc.Publish("foo", []byte("Hello World")) // does not return until an ack has been received from NATS Streaming

// Simple Async Subscriber
sub, _ := sc.Subscribe("foo", func(m *stan.Msg) {
    fmt.Printf("Received a message: %s\n", string(m.Data))
})

// Unsubscribe
sub.Unsubscribe()

// Close connection
sc.Close()