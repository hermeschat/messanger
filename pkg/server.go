package pkg

import (
	"fmt"
	"git.raad.cloud/cloud/hermes/pkg/api"
	"git.raad.cloud/cloud/hermes/pkg/eventHandler"
	"git.raad.cloud/cloud/hermes/pkg/join"
	"git.raad.cloud/cloud/hermes/pkg/newMessage"
	"git.raad.cloud/cloud/hermes/pkg/session"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"golang.org/x/net/context"
)

type HermesServer struct {

}



func (h HermesServer) Echo(ctx context.Context, s *api.Some) (*api.Empty, error) {
	fmt.Println("Echo")
	return &api.Empty{}, nil
}

func (h HermesServer) NewMessage(msg api.Hermes_NewMessageServer) error {
	m,err := msg.Recv()
	if err != nil {
		logrus.Infof("canot recieve : %v", err)
	}
	nm := &newMessage.NewMessage{
		Body:        m.Body,
		From:        m.From,
		To:          m.To,
		Channel:     m.Channel,
		MessageType: m.MessageType,
		Session:     "",
	}

	err = newMessage.Handle(nm)
	if err != nil {
		return errors.Wrap(err, "error in new message")
	}

	return nil
}
func (h HermesServer) EventBuff(a api.Hermes_EventBuffServer) error {
	eventHandler.UserSockets.Lock()
	eventHandler.UserSockets.Us[m.SessionID] = &msg
	eventHandler.UserSockets.Unlock()
	panic("implement me")
}

func (h HermesServer) Join(ctx context.Context, message *api.JoinSignal) (*api.Empty, error) {
	jp := &join.JoinPayload{
		UserID:    message.UserID, //should get from jwt
		SessionId: message.SessionId,
	}
	logrus.Infof(message.SessionId)

	err := join.Handle(jp)
	if err != nil {
		return &api.Empty{Status: "500"}, errors.Wrap(err, "error in joining")
	}
	return &api.Empty{Status: "200"}, nil
}

func (h HermesServer) CreateSession(ctx context.Context, req *api.CreateSessionRequest) (*api.CreateSessionResponse, error) {
	cs := &session.CreateSession{
		UserIP:        req.GetUserIP(),
		UserID:        req.UserID, //from jwt
		ClientVersion: req.ClientVersion,
		Node:          req.Node,
	}
	logrus.Println("created session model")
	s, err := session.Create(cs)
	if err != nil {
		return &api.CreateSessionResponse{}, errors.Wrap(err, "error in creating session")
	}
	logrus.Println("done")
	return &api.CreateSessionResponse{
		SessionID:s.SessionID,
	}, nil
}
//
//func (h HermesServer) Deliverd(ctx context.Context, message *api.DeliveredSignal) (*api.Empty, error) {
//	ds := &delivered.DeliverdSignal{
//		MessageID: message.MessageID,
//		ChannelID: message.ChannelID,
//	}
//	err := delivered.Handle(ds)
//	if err != nil {
//		return &api.Empty{Status: "500"}, errors.Wrap(err, "error in delivering message")
//	}
//	return &api.Empty{Status: "200"}, nil
//}
//func (h HermesServer) Read(ctx context.Context, message *api.ReadSignal) (*api.Empty, error) {
//	rs := &read.ReadSignal{
//		MessageID: message.MessageID,
//		ChannelID: message.ChannelID,
//	}
//	err := read.Handle(rs)
//	if err != nil {
//		return &api.Empty{Status: "500"}, errors.Wrap(err, "error in reading")
//	}
//	return &api.Empty{Status: "200"}, nil
//}


//func (h HermesServer) DestroySession(context.Context, *api.DestroySessionRequest) (*api.Empty, error) {
//	panic("implement me")
//}

//
//func (h HermesServer) KeepAlive(context.Context, *api.Message) (*api.Empty, error) {
//	panic("implement me")
//}