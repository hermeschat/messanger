package pkg

import (
	"git.raad.cloud/cloud/hermes/pkg/api"
	"git.raad.cloud/cloud/hermes/pkg/delivered"
	"git.raad.cloud/cloud/hermes/pkg/join"
	"git.raad.cloud/cloud/hermes/pkg/newMessage"
	"git.raad.cloud/cloud/hermes/pkg/read"
	"git.raad.cloud/cloud/hermes/pkg/session"
	"github.com/pkg/errors"
	"golang.org/x/net/context"
	"time"
)

type HermesServer struct {
	ClusterID   string
	NatsSrvAddr string
}

func (h HermesServer) KeepAlive(context.Context, *api.Message) (*api.Empty, error) {
	panic("implement me")
}

func (h HermesServer) NewMessage(ctx context.Context, message *api.Message) (*api.Empty, error) {
	nm := &newMessage.NewMessage{
		Body:message.Body,
		From:message.From,
		To:message.To,
		Channel:message.Channel,
		MessageType:message.MessageType,
		Session:"",
	}
	err := newMessage.Handle(nm)
	if err != nil {
		return &api.Empty{
			Status:"500",
		}, errors.Wrap(err, "error in new message")
	}
	return &api.Empty{Status:"200"}, nil
}

func (h HermesServer) Join(ctx context.Context, message *api.JoinSignal) (*api.Empty, error) {
	jp := &join.JoinPayload{
		UserID:"", //should get from jwt
		SessionId: message.SessionId,
	}
	err := join.Handle(jp)
	if err != nil {
		return &api.Empty{Status:"500"}, errors.Wrap(err, "error in joining")
	}
	return &api.Empty{Status:"200"}, nil
}

func (h HermesServer) Deliverd(ctx context.Context, message *api.DeliveredSignal) (*api.Empty, error) {
	ds := &delivered.DeliverdSignal{
		MessageID:message.MessageID,
		ChannelID:message.ChannelID,
	}
	err := delivered.Handle(ds)
	if err != nil {
		return &api.Empty{Status:"500"}, errors.Wrap(err, "error in delivering message")
	}
	return &api.Empty{Status:"200"}, nil
}
func (h HermesServer) Read(ctx context.Context, message *api.ReadSignal) (*api.Empty, error) {
	rs := &read.ReadSignal{
		MessageID: message.MessageID,
		ChannelID:message.ChannelID,
	}
	err := read.Handle(rs)
	if err != nil {
		return &api.Empty{Status:"500"}, errors.Wrap(err, "error in reading")
	}
	return &api.Empty{Status:"200"}, nil
}
func (h HermesServer) CreateSession(ctx context.Context, req *api.CreateSessionRequest) (*api.Empty, error) {
	cs := &session.CreateSession{
		UserIP: "", // az ye jayi
		UserID:"", //from jwt
		ClientVersion:req.ClientVersion,
		Node:req.Node,
	}
	_, err := session.Create(cs)
	if err != nil {
		return &api.Empty{Status:"error"}, errors.Wrap(err, "error in creating session")
	}
	return &api.Empty{Status:"OK"}, nil
}

func (h HermesServer) DestroySession(context.Context, *api.DestroySessionRequest) (*api.Empty, error) {
	panic("implement me")
}

type HermesContext struct {
	DeafultCtx context.Context
	User       struct {
		ID         string
		AppID      string
		ExpireDate time.Time
		Roles      []string
	}
}
