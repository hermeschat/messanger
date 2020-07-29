package core

import (
	"fmt"
	hermesproto "github.com/hermeschat/proto"
	"github.com/golang/protobuf/proto"
	"github.com/hermeschat/engine/monitoring"
	"github.com/nats-io/stan.go"
)

func userDiscoveryNatsHandler(msg *stan.Msg) {
	// subsc
}

func (c *chatService) newMessageNatsHandler(msg *stan.Msg) {
	var err error
	for _, p := range c.ps {
		pmessage := &hermesproto.Message{}
		err = proto.Unmarshal(msg.Data, pmessage)
		err = p.Push(fmt.Sprint(pmessage.To), pmessage)
	}
	if err != nil {
		monitoring.Logger().Warnf("all push tries failed :%s", err)
	}
}
