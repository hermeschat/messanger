package pkg

import "git.raad.cloud/cloud/hermes/pkg/api"

type HermesServer struct{
	ClusterID string
	NatsSrvAddr string
}

func (h *HermesServer) KeepAlive(context.Context, *api.Signal) (*api.Response, error) {
	panic("implement me")
}

func (h *HermesServer) NewMessage(context.Context, *api.InstantMessage) (*api.Response, error) {

}

func (h *HermesServer) Deliverd(context.Context, *api.Signal) (*api.Response, error) {
	panic("implement me")
}

func (h *HermesServer) CreateSession(context.Context, *api.CreateSessionRequest) (*api.CreateSessionResponse, error) {
	panic("implement me")
}

func (h *HermesServer) DestroySession(context.Context, *api.DestroySessionRequest) (*api.DestroySessionResponse, error) {
	panic("implement me")
}


