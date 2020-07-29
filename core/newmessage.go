package core

import (
	"context"
	"fmt"
	gproto "github.com/golang/protobuf/proto"

	"github.com/hermeschat/engine/models"
	"github.com/hermeschat/engine/monitoring"
	"github.com/hermeschat/proto"
	"github.com/volatiletech/null/v8"
	"github.com/volatiletech/sqlboiler/v4/boil"
)
const (
	MessageStateInServer = iota + 1
	MessageStateSent
	MessageStateDelivered
	MessageStateRead
)
func (c *chatService) NewMessage(ctx context.Context, msg *proto.Message) error {
	message := &models.Message{
		OriginID:  null.Int{int(msg.From), true},
		DSTID:     null.Int{int(msg.To), true},
		ParentID:  null.Int{int(msg.Parent), true},
		Body:      null.String{msg.Body, true},
		State:     null.Int{int(MessageStateInServer), true},
	}
	err := message.Insert(ctx, c.db, boil.Infer())
	if err != nil {
	 	return err
	}
	msg.MessageID = fmt.Sprint(message.ID)
	go func() {
		err = c.pushMessageToChannel(fmt.Sprint(message.DSTID), msg)
		if err != nil {
		 	monitoring.Logger().Errorf("error in sending into nats channel: %s", err)
		}
	}()
	if err != nil {
		return err
	}
	return nil
}
func (c *chatService) pushMessageToChannel(channelId string, data *proto.Message) error {
	//TODO: check for authorization
	bs, err := gproto.Marshal(data)
	if err != nil {
	 	return err
	}
	guid, err := c.nc.PublishAsync(channelId, bs, nil)
	_ = guid // TODO:
	if err != nil {
		return err
	}
	return nil
}
/*
Telegram:
	Each person has a channel, each group is a channel, each kanal is a channel
	Channel ( kanal )
	Private => MNIM, COMRADE => MNIM -> *COMRADE
	{
		UserID
		ChannelID
	}

*/