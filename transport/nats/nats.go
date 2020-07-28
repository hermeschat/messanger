package nats

import (
	"github.com/hermeschat/engine/config"
	"github.com/hermeschat/engine/monitoring"
	nats_driver "github.com/nats-io/nats.go"
	"time"
)

func NatsClient() (*nats_driver.Conn, error) {
	return nats_driver.Connect(config.NatsURI())

}

func HealthCheck() bool {
	nc, err := nats_driver.Connect(nats_driver.DefaultURL, nats_driver.Name("API Ping Example"),
		nats_driver.PingInterval(20*time.Second), nats_driver.MaxPingsOutstanding(5))
	if err != nil {
		monitoring.Logger().Fatalf("can not connect to nats stream")

	}
	if nc.IsConnected() {
		monitoring.Logger().Infof("connected to nats!")
		return true
	}
	return false
}
