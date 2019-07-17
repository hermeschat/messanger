package pkg

import (
	"time"

	"git.raad.cloud/cloud/hermes/pkg/api"
	"git.raad.cloud/cloud/hermes/pkg/auth"
	"git.raad.cloud/cloud/hermes/pkg/eventHandler"
	"git.raad.cloud/cloud/hermes/pkg/join"
	"git.raad.cloud/cloud/hermes/pkg/newMessage"
	"git.raad.cloud/cloud/hermes/pkg/read"
	"git.raad.cloud/cloud/hermes/pkg/session"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"golang.org/x/net/context"
)

type HermesServer struct {
	Ctx context.Context
}

var AppContext = context.Background()

func (h HermesServer) Echo(ctx context.Context, a *api.Empty) (*api.Empty, error) {

	logrus.Infof("Identity is :\n %+v", ctx.Value("identity"))
	return &api.Empty{Status: "JWT is ok"}, nil
}

func (h HermesServer) ListChannels(context.Context, *api.Empty) (*api.Channels, error) {
	panic("implement")
}

func (h HermesServer) ListMessages(context.Context, *api.Empty) (*api.Messages, error) {
	return &api.Messages{
		Msg: []*api.Message{
			&api.Message{From: "5c4c2683bfd02a2b923af8be", To: "5c4c2683bfd02a2b923af8bf", Body: "salam aleyk"},
		},
	}, nil
}

func (h HermesServer) EventBuff(a api.Hermes_EventBuffServer) error {
	ctx := a.Context()
	i := ctx.Value("identity")
	ident, ok := i.(*auth.Identity)
	logrus.Info(">>>>>>> We Are in Event Buff ")
	if !ok {
		logrus.Errorf("Cannot get identity out of context")
	}

	time.Sleep(time.Second)
	for{
	e, err := a.Recv()
	if err != nil {
		logrus.Errorf("cannot receive event : %v", err)
		return errors.Wrap(err, "error in reading EventBuff")
	}
	eventHandler.UserSockets.Lock()
	eventHandler.UserSockets.Us[ident.ID] = a
	eventHandler.UserSockets.Unlock()
	logrus.Info("we have a new event")

	switch t := e.GetEvent().(type) {
	case *api.Event_Read:
		logrus.Info("Event is read")
		r := e.GetRead()
		rs := &read.ReadSignal{
			UserID:    ident.ID,
			MessageID: r.MessageID,
			ChannelID: r.ChannelID,
		}
		err = read.Handle(rs)
		if err != nil {
			logrus.Errorf("Error in handling read signal")
		}
	case *api.Event_Keep:
		logrus.Info("Event is keep")
		k := e.GetKeep()
		_ = k
		//find logic
	case *api.Event_NewMessage:
		logrus.Info("Event is New Message")
		m := e.GetNewMessage()
		if m != nil {
			logrus.Info("Event is NewMessage")
			nm := &newMessage.NewMessage{
				Body:        m.Body,
				From:        ident.ID,
				To:          m.To,
				Channel:     m.Channel,
				MessageType: m.MessageType,
				Session:     "",
			}

			err = newMessage.Handle(nm)
			if err != nil {
				logrus.Errorf("Error in NewMessage Event : %v", err)
			}
		}
		//return nil
	case *api.Event_Join:
		j := e.GetJoin()
		logrus.Info(j)
		if j != nil {
			logrus.Info("Event is Join")
			jp := &join.JoinPayload{
				UserID:    ident.ID, //should get from jwt
				SessionId: j.SessionId,
			}

			join.Handle(h.Ctx, jp)
			//if err != nil {
			//	logrus.Errorf("Error in Join event : %v", err)
			//}
		}
		//return nil
	default:
		logrus.Infof("Type not matched : %+T", t)
	}
	//goto loop
	}
	return nil
}

func (h HermesServer) CreateSession(ctx context.Context, req *api.CreateSessionRequest) (*api.CreateSessionResponse, error) {
	i := ctx.Value("identity")
	ident, ok := i.(*auth.Identity)
	if !ok {
		return &api.CreateSessionResponse{}, errors.New("could not get identity")
	}
	cs := &session.CreateSession{
		UserIP:        req.GetUserIP(),
		UserID:        ident.ID, //from jwt
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
		SessionID: s.SessionID,
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
