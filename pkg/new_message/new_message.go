package new_message

import (
	"git.raad.cloud/cloud/hermes/pkg/api"
	"git.raad.cloud/cloud/hermes/pkg/channel"
	"git.raad.cloud/cloud/hermes/pkg/drivers/nats"
	message2 "git.raad.cloud/cloud/hermes/pkg/message"
	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/bson"
	"time"
)

func Handle(message *api.InstantMessage) (api.Response)  {
	if message.To == ""  && message.Channel == ""{
		return api.Response{ Error: errors.New("Channel ID or To should be present in payload").Error()}
	}
	if message.Channel != "" {
		err := nats.Publish("test-cluster", "0.0.0.0:4222",message.Channel, message)
		if err != nil {
			return api.Response{Error:errors.Wrap(err, "Cannot publish message to nats").Error()}
		}
		return api.Response{
			Code:"200",
		}
	}
	if message.To != "" {
		channels, err := channel.GetAll(bson.M{
			"Members": bson.M{"$in" : []string{message.From,message.To}},
		})
		if err != nil {
			return api.Response{
				Error: errors.Wrap(err, "Cannot get channels").Error(),
			}
		}
		var targetChannel *channel.Channel
		if len(*channels) < 1 {
			targetChannel = &channel.Channel{
				Members: []string{message.To, message.From},
			}
			err := saveChannelToMongo(targetChannel)
			if err != nil {
				return api.Response{
					Error:"Internal Service problem",
				}
			}
		} else {
			targetChannel = (*channels)[0]
		}
		err = nats.Publish("test-cluster", "0.0.0.0", targetChannel.ChannelID, message )
		if err != nil {
			return api.Response{
				Error: errors.Wrap(err, "Error while publishing to NATS").Error(),
			}
		}
		err = saveMessageToMongo(message)
		if err != nil {
			return api.Response{
				Error : errors.Wrap(err, "Error in saving to mongo").Error(),
			}
		}
		return api.Response{
			Code:"200",
		}

	}
	return api.Response{
		Code: "Unknown",
	}
}
func saveChannelToMongo(c *channel.Channel) error {
	err := channel.Add(c)
	if err != nil {
		return errors.Wrap(err, "cannot save to mongo")
	}
	return nil
}

func saveMessageToMongo(message *api.InstantMessage) error {
	err := message2.Add(&message2.Message{
		To : message.To,
		From: message.From,
		Time: time.Now(),
		MessageType: message.MessageType,
		Body: message.Body,
	})
	if err != nil {
		return errors.Wrap(err, "cannot save to mongo")
	}
	return nil
}