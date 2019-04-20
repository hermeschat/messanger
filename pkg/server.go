package pkg

import (
	"git.raad.cloud/cloud/hermes/pkg/api"
	"git.raad.cloud/cloud/hermes/pkg/delivered"
	"git.raad.cloud/cloud/hermes/pkg/join"
	"git.raad.cloud/cloud/hermes/pkg/new_message"
	"git.raad.cloud/cloud/hermes/pkg/read"
	"golang.org/x/net/context"
)

type HermesServer struct{
	ClusterID string
	NatsSrvAddr string
}

func (h *HermesServer) KeepAlive(context.Context, *api.Signal) (*api.Response, error) {
	panic("implement me")
}

func (h *HermesServer) NewMessage(ctx context.Context, message *api.InstantMessage) (*api.Response, error) {
	return new_message.Handle(message), nil
}

func (h *HermesServer) Join(ctx context.Context, message *api.Signal) (*api.Response, error){
	return join.Handle(message), nil
}

func (h *HermesServer) Deliverd(ctx context.Context, message *api.Signal) (*api.Response, error) {
	return delivered.Handle(message), nil
}
func (h *HermesServer) Read(ctx context.Context, message *api.Signal) (*api.Response, error) {
	return read.Handle(message),nil
}
func (h *HermesServer) CreateSession(context.Context, *api.CreateSessionRequest) (*api.CreateSessionResponse, error) {
	panic("implement me")
}

func (h *HermesServer) DestroySession(context.Context, *api.DestroySessionRequest) (*api.DestroySessionResponse, error) {
	panic("implement me")
}


