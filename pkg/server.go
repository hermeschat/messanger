package pkg

import (
	"fmt"

	"git.raad.cloud/cloud/hermes/pkg/api"
	"git.raad.cloud/cloud/hermes/pkg/auth"
	"git.raad.cloud/cloud/hermes/pkg/drivers/nats"
	"git.raad.cloud/cloud/hermes/pkg/drivers/redis"
	"git.raad.cloud/cloud/hermes/pkg/eventHandler"
	"git.raad.cloud/cloud/hermes/pkg/newMessage"
	"git.raad.cloud/cloud/hermes/pkg/read"
	"git.raad.cloud/cloud/hermes/pkg/repository/channel"
	"git.raad.cloud/cloud/hermes/pkg/repository/message"
	"git.raad.cloud/cloud/hermes/pkg/session"
	"github.com/mitchellh/mapstructure"
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

func (h HermesServer) ListChannels(ctx context.Context, _ *api.Empty) (*api.Channels, error) {
	fmt.Println("hereinlistmessages")
	i := ctx.Value("identity")
	ident, ok := i.(*auth.Identity)
	_ = ident
	if !ok {
		return nil, errors.New("cannot get identity out of context")
	}
	msgs, err := channel.GetAll(map[string]interface{}{
		"Members": map[string]interface{}{
			"$in": []string{ident.ID},
		}, //TODO fix query
	})
	if err != nil {
		return nil, errors.Wrap(err, "error while trying to get messages from database")
	}
	output := []*api.Channel{}
	for _, m := range msgs {
		amsg := &api.Channel{}
		err = mapstructure.Decode(m, amsg)
		if err != nil {
			return nil, errors.Wrap(err, "error while converting from repository message to api message")
		}
		output = append(output, amsg)
	}
	return &api.Channels{Msg: output}, nil
}

func (h HermesServer) ListMessages(ctx context.Context, _ *api.Empty) (*api.Messages, error) {
	i := ctx.Value("identity")
	ident, ok := i.(*auth.Identity)
	if !ok {
		return nil, errors.New("cannot get identity out of context")
	}
	chns, err := channel.GetAll(map[string]interface{}{
		"Members": map[string]interface{}{
			"$in": []string{ident.ID},
		},
	})
	var chnIds []string
	for _, chn := range chns {
		chnIds = append(chnIds, chn.ChannelID)
	}
	if err != nil {
		return nil, errors.Wrap(err, "error in getting channels that user is member of")
	}
	msgs, err := message.GetAll(map[string]interface{}{
		"$or": []map[string]interface{}{{"To": ident.ID}, {"From": ident.ID}, {"ChannelID": map[string]interface{}{
			"$in": chnIds,
		}}},
	})
	if err != nil {
		return nil, errors.Wrap(err, "error while trying to get messages from database")
	}
	output := []*api.Message{}
	for _, m := range msgs {
		fmt.Println(m)
		amsg := &api.Message{}
		err = mapstructure.Decode(m, amsg)
		if err != nil {
			return nil, errors.Wrap(err, "error while converting from repository message to api message")
		}
		output = append(output, amsg)
	}
	fmt.Printf("\n%+v", output)
	return &api.Messages{Msg: output}, nil
}

func (h HermesServer) EventBuff(a api.Hermes_EventBuffServer) error {
	ctx := a.Context()
	i := ctx.Value("identity")
	ident, ok := i.(*auth.Identity)
	logrus.Info(">>>>>>> We Are in Event Buff ")
	if !ok {
		logrus.Errorf("Cannot get identity out of context")
	}
	defer func() {
		con, err := redis.ConnectRedis()
		if err != nil {
			logrus.Errorf("error while trying to clear redis cache of subscribed channels : %v", err)
			return
		}
		_, err = con.Del(ident.ID).Result()
		if err != nil {
			logrus.Errorf("error while trying to clear redis cache of subscribed channels : %v", err)
			return
		}
		nats.State.Mu.Lock()
		natsCon, ok := nats.State.Ss[ident.ID]
		if !ok {
			logrus.Errorf("user nats connection not found")
			return
		}
		err = (*natsCon).Close()
		if err != nil {
			logrus.Errorf("error while trying to close user nats connection")
			return
		}
		delete(nats.State.Ss, ident.ID)
		nats.State.Mu.Unlock()
	}() //loop to continuously read messages from buffer
	for {

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
				jp := &eventHandler.JoinPayload{
					UserID:    ident.ID, //should get from jwt
					SessionId: j.SessionId,
				}

				eventHandler.Handle(h.Ctx, jp)
				//if err != nil {
				//	logrus.Errorf("Error in Join event : %v", err)
				//}
			}
			//return nil
		default:
			logrus.Infof("Type not matched : %+T", t)
		}
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
