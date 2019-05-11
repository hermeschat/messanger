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

func (h *HermesServer) KeepAlive(context.Context, *api.Message) (*api.Response, error) {
	panic("implement me")
}

func (h *HermesServer) NewMessage(ctx context.Context, message *api.Message) (*api.Empty, error) {
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

func (h *HermesServer) Join(ctx context.Context, message *api.JoinSignal) (*api.Empty, error) {
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

func (h *HermesServer) Deliverd(ctx context.Context, message *api.Message) (*api.Response, error) {
	return delivered.Handle(message), nil
}
func (h *HermesServer) Read(ctx context.Context, message *api.Message) (*api.Response, error) {
	return read.Handle(message), nil
}
func (h *HermesServer) CreateSession(ctx context.Context, req *api.CreateSessionRequest) (*api.Response, error) {
	return session.Create(req), nil
}

func (h *HermesServer) DestroySession(context.Context, *api.DestroySessionRequest) (*api.Message, error) {
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
