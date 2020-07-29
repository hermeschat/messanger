package core

import (
	"github.com/hermeschat/engine/monitoring"
	"github.com/nats-io/stan.go"
)

func userDiscoveryNatsHandler(msg *stan.Msg) {
	// subsc
}

func (c *chatService) newMessageNatsHandler(msg *stan.Msg) {
	var err error
	for _, p := range c.ps {
		err = p.Push(msg.Data)
	}
	if err != nil {
		monitoring.Logger().Warnf("all push tries failed :%s", err)
	}
}
