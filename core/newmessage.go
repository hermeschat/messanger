package core

import (
	"context"
	"github.com/hermeschat/engine/models"
	"github.com/volatiletech/sqlboiler/v4/boil"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"
)

func (c *chatService) NewMessage(ctx context.Context, message *models.Message) error {
	err := message.Insert(ctx, c.db, boil.Infer())
	if err != nil {
	 	return err
	}
	ch, err := models.Channels(qm.Where("id=$1", message.DSTID)).One(ctx, c.db)
	if err != nil {
		return err
	}
	for _, cm := range ch.R.ChannelMembers { // should be done concurrently
		err = pushToChannel(cm.R.User.Username)
		if err != nil {
		 	return err
		}
	}
	return nil
}
func NewMessageEventHandler() {}
func pushToChannel(channelId string) error {}
/*
Telegram:
	Each person has a channel, each group is a channel, each kanal is a channel
	Channel ( kanal )
	Private => MNIM, COMRADE => MNIM -> *COMRADE
	Group


*/