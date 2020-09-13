package core

import (
	"context"
	"fmt"
	"github.com/hermeschat/engine/models"
	"github.com/nats-io/stan.go"
	"github.com/volatiletech/null/v8"
	"github.com/volatiletech/sqlboiler/v4/boil"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"
)

const USERDISCOVERY = "user-discovery"

func (c *chatService) Join(userID int) error {
	cc, err := models.ChannelMembers(qm.Where("user_id=$1", userID)).Count(context.TODO(), c.db)
	if err != nil {
		return err
	}
	if cc <= 0 {
		// create channel in database
		prvChan := &models.Channel{
			CreatorID: null.Int{
				Int:   userID,
				Valid: true,
			}}
		err = prvChan.Insert(context.TODO(), c.db, boil.Infer())
		if err != nil {
			return err
		}
		cm := &models.ChannelMember{
			ChannelID: prvChan.ID,
			UserID:    userID,
		}
		err = cm.Insert(context.TODO(), c.db, boil.Infer())
		if err != nil {
			return err
		}
	}
	err = c.subscribeToChannel(USERDISCOVERY, userDiscoveryNatsHandler)
	if err != nil {
		return err
	}
	err = c.subscribeToChannel(fmt.Sprint(userID), c.newMessageNatsHandler)
	if err != nil {
		return err
	}
	cms, err := models.ChannelMembers(qm.Where("user_id=$1", userID)).All(context.TODO(), c.db)
	if err != nil {
		return err
	}
	for _, cm := range cms {
		err = c.subscribeToChannel(fmt.Sprint(cm.R.Channel.ID), c.newMessageNatsHandler)
		if err != nil {
			return err
		}
	}
	return nil
}

func (c *chatService) subscribeToChannel(channelID string, handler stan.MsgHandler) error {
	sub, err := c.nc.Subscribe(channelID, handler)
	if err != nil {
		return err
	}

	_ = sub // do something pls
	return nil
}
