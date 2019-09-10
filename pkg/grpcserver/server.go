package grpcserver

import (
	"fmt"
	"net"
	"sync"

	"github.com/amirrezaask/config"
	"hermes/api"
	auth "hermes/paygearauth"
	"hermes/pkg/db"
	"hermes/pkg/discovery"
	"hermes/pkg/drivers/nats"
	"hermes/pkg/drivers/redis"
	"hermes/pkg/eventhandlers"

	grpcmiddleware "github.com/grpc-ecosystem/go-grpc-middleware"

	grpc_auth "github.com/grpc-ecosystem/go-grpc-middleware/auth"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

type hermesServer struct {
	Ctx context.Context
}

var userSockets = &struct {
	sync.RWMutex
	Us map[string]api.Hermes_EventBuffServer
}{
	sync.RWMutex{},
	map[string]api.Hermes_EventBuffServer{},
}

//CreateGRPCServer creates a new grpc server
func CreateGRPCServer(ctx context.Context) {
	logrus.Infof("hermes GRPC server server is on 0.0.0.0:%s", config.Get("port"))
	lis, err := net.Listen("tcp", fmt.Sprintf(":%s", config.Get("port")))
	if err != nil {
		logrus.Fatal("ERROR can't create a tcp listener ")
	}
	streamChain := grpcmiddleware.ChainStreamServer(grpc_auth.StreamServerInterceptor(unaryAuthJWTInterceptor))
	unaryChain := grpcmiddleware.ChainUnaryServer(grpc_auth.UnaryServerInterceptor(unaryAuthJWTInterceptor))
	logrus.Info("Interceptors Created")
	srv := grpc.NewServer(grpc.StreamInterceptor(streamChain), grpc.UnaryInterceptor(unaryChain))
	api.RegisterHermesServer(srv, hermesServer{ctx})
	logrus.Info("Registering Hermes GRPC")
	err = srv.Serve(lis)
	if err != nil {
		logrus.Fatal("ERROR in serving listener")
	}

	logrus.Info("GRPC is Live !!!")
}

func (h hermesServer) CreateSession(context.Context, *api.CreateSessionRequest) (*api.CreateSessionResponse, error) {
	panic("implement me")
}

func (h hermesServer) ListChannels(ctx context.Context, _ *api.Empty) (*api.Channels, error) {
	i := ctx.Value("identity")
	ident, ok := i.(*auth.Identity)
	_ = ident
	if !ok {
		return nil, errors.New("cannot get identity out of context")
	}
	chanCur, err := db.Channels().Find(ctx, map[string]interface{}{
		"members": map[string]interface{}{
			"$in": []string{ident.ID},
		},
	})
	if err != nil {
		return nil, errors.Wrap(err, "error while trying to get messages from db")
	}
	var chans []*api.Channel
	for chanCur.Next(ctx) {
		channel := new(api.Channel)
		err = chanCur.Decode(channel)
		if err != nil {
			return nil, errors.Wrap(err, "error while decoding channels into api.channel")
		}
		chans = append(chans, channel)
	}
	return &api.Channels{Msg: chans}, nil
}

func (h hermesServer) ListMessages(ctx context.Context, ch *api.ChannelID) (*api.Messages, error) {
	//msgCur, err := db.Messages().Find(ctx, map[string]interface{}{
	//	"channel_id": ch.Id,
	//})
	//if err != nil {
	//	return nil, errors.Wrap(err, "error while trying to get messages from db")
	//}
	//var messages []*api.Message
	//for msgCur.Next(ctx) {
	//	thisMessage := new(api.Message)
	//	err = msgCur.Decode(thisMessage)
	//	if err != nil {
	//		return nil, errors.Wrap(err, "error while decoding message")
	//	}
	//}
	//return &api.Messages{Msg: messages}, nil
	return nil, nil
}

func (h hermesServer) GetChannel(ctx context.Context, _ *api.ChannelID) (*api.Channel, error) {
	return nil, nil
}

var appContext = context.Background()

func (h hermesServer) Echo(ctx context.Context, a *api.Empty) (*api.Empty, error) {

	return &api.Empty{Status: "JWT is ok"}, nil
}

//EventBuff ..
func (h hermesServer) EventBuff(a api.Hermes_EventBuffServer) error {
	ctx := a.Context()
	i := ctx.Value("identity")
	ident, ok := i.(*auth.Identity)
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
		err = nats.Connections.CloseConnection(ident.ID)
		if err != nil {
			logrus.Errorf("erorr in closing nats connection: %v", err)
			return
		}
	}()
	if !eventhandlers.UserIsSubscribedTo(ident.ID, "user-discovery") {
		logrus.Infoln("User is not subscribed to user discovery")
		go nats.MakeSubscriber(a.Context(), ident.ID, "user-discovery", discovery.UserDiscoveryEventHandler(a.Context(), ident.ID, userSockets))()
	}
	//loop to continuously read messages from buffer
	for {
		e, err := a.Recv()
		if err != nil {
			logrus.Errorf("cannot receive event : %v", err)
			return errors.Wrap(err, "error in reading EventBuff")
		}
		userSockets.Lock()
		userSockets.Us[ident.ID] = a
		userSockets.Unlock()

		switch t := e.GetEvent().(type) {
		case *api.Event_Read:
			logrus.Info("Event is read")
			r := e.GetRead()
			rs := &eventhandlers.ReadSignal{
				UserID:    ident.ID,
				MessageID: r.MessageID,
				ChannelID: r.ChannelID,
			}
			err = eventhandlers.HandleRead(rs)
			if err != nil {
				logrus.Errorf("Error in handling read signal")
			}
		case *api.Event_Keep:
			logrus.Info("Event is keep")
			k := e.GetKeep()
			_ = k
		case *api.Event_NewMessage:
			logrus.Info("Event is New Message")
			m := e.GetNewMessage()
			if m != nil {
				logrus.Info("Event is Message")
				nm := &db.Message{
					Body:        m.Body,
					From:        ident.ID,
					To:          m.To,
					ChannelID:   m.Channel,
					MessageType: m.MessageType,
				}

				err = eventhandlers.HandleNewMessage(nm)
				if err != nil {
					logrus.Errorf("Error in Message Event : %v", err)
				}
			}
			//return nil
		case *api.Event_Join:
			j := e.GetJoin()
			logrus.Info(j)
			if j != nil {
				logrus.Info("Event is Join")
				jp := &eventhandlers.JoinPayload{
					UserID:    ident.ID, //should get from jwt
					SessionId: j.SessionId,
				}

				eventhandlers.HandleJoin(a.Context(), jp, userSockets)

			}
		default:
			logrus.Infof("Type not matched : %+T", t)
		}
	}

}
