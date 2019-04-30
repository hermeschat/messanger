package pkg

import (
	"git.raad.cloud/cloud/hermes/pkg/api"
	"git.raad.cloud/cloud/hermes/pkg/delivered"
	"git.raad.cloud/cloud/hermes/pkg/join"
	"git.raad.cloud/cloud/hermes/pkg/newMessage"
	"git.raad.cloud/cloud/hermes/pkg/read"
	"git.raad.cloud/cloud/hermes/pkg/session"
	"golang.org/x/net/context"
)

type HermesServer struct {
	ClusterID   string
	NatsSrvAddr string
}

func (h *HermesServer) KeepAlive(context.Context, *api.Message) (*api.Response, error) {
	panic("implement me")
}

func (h *HermesServer) NewMessage(ctx context.Context, message *api.Message) (*api.Response, error) {
	return newMessage.Handle(message,""), nil
}

func (h *HermesServer) Join(ctx context.Context, message *api.Message) (*api.Response, error) {
	return join.Handle(message), nil
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
