package grpcserver

import (
	"fmt"
	"net"
	"sync"

	"github.com/amirrezaask/config"
	"go.mongodb.org/mongo-driver/bson"
	"hermes/api"
	"hermes/db"
	"hermes/discovery"
	"hermes/eventhandlers"
	auth "hermes/paygearauth"
	"hermes/subscription"
	"hermes/subscription/nats"

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

func (h hermesServer) GetChannel(ctx context.Context, ch *api.ChannelID) (*api.Channel, error) {
	res := db.Channels().FindOne(ctx, bson.M{"_id": ch.Id})
	if err := res.Err(); err != nil {
		return nil, errors.Wrap(err, "error in find channel")
	}
	channel := new(api.Channel)
	err := res.Decode(channel)
	if err != nil {
		return nil, errors.Wrap(err, "error while decoding channel")
	}
	return channel, nil
}

var appContext = context.Background()

func (h hermesServer) Echo(ctx context.Context, a *api.Empty) (*api.Empty, error) {
	return &api.Empty{}, nil
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
		con, err := subscription.Redis()
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
	if !subscription.UserIsSubscribedTo(ident.ID, "user-discovery") {
		logrus.Infoln("User is not subscribed to user discovery")
		go subscription.NewSubsciption(a.Context(), ident.ID, "user-discovery", discovery.UserDiscoveryEventHandler(a.Context(), ident.ID, userSockets))
	}
	userSockets.Lock()
	userSockets.Us[ident.ID] = a
	userSockets.Unlock()
	for {
		e, err := a.Recv()
		if err != nil {
			logrus.Errorf("cannot receive event : %v", err)
			return errors.Wrap(err, "error in reading EventBuff")
		}
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
		case *api.Event_Dlv:
			logrus.Info("Event is Delivered")
		default:
			logrus.Infof("Type not matched : %+T", t)
		}
	}

}
