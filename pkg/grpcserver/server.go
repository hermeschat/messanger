package grpcserver

import (
	"fmt"
	"net"
	"sync"

	"github.com/amirrezaask/config"
	"hermes/api"
	auth "hermes/paygearauth"
	"hermes/pkg/db"
	"hermes/pkg/drivers/nats"
	"hermes/pkg/drivers/redis"
	"hermes/pkg/join"
	"hermes/pkg/message"
	"hermes/pkg/read"

	grpcmiddleware "github.com/grpc-ecosystem/go-grpc-middleware"

	grpc_auth "github.com/grpc-ecosystem/go-grpc-middleware/auth"
	"github.com/mitchellh/mapstructure"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson/primitive"
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
	msgs, err := db.Instance().Channels.Get(map[string]interface{}{
		"Members": map[string]interface{}{
			"$in": []string{ident.ID},
		}, //TODO fix query
	})
	if err != nil {
		return nil, errors.Wrap(err, "error while trying to get messages from db")
	}
	output := []*api.Channel{}
	for _, m := range msgs {
		amsg := &api.Channel{}
		amsg.ChannelId = fmt.Sprint(m["ChannelID"])
		members := []string{}
		for _, mem := range m["Members"].(primitive.A) {
			members = append(members, fmt.Sprint(mem))
		}
		amsg.Members = members
		roles := map[string]string{}
		for member, role := range m["Roles"].(map[string]interface{}) {
			roles[member] = fmt.Sprint(role)
		}
		amsg.Roles = roles
		amsg.Type = fmt.Sprint(m["Type"].(int32))
		output = append(output, amsg)
	}
	return &api.Channels{Msg: output}, nil
}

func (h hermesServer) ListMessages(ctx context.Context, ch *api.ChannelID) (*api.Messages, error) {
	msgs, err := db.Instance().Channels.Get(map[string]interface{}{
		"ChannelID": ch.Id,
	})
	if err != nil {
		return nil, errors.Wrap(err, "error while trying to get messages from db")
	}
	output := []*api.Message{}
	for _, msg := range msgs {
		msg["MessageID"] = fmt.Sprint(msg["MessageID"].(primitive.ObjectID).Hex())
	}
	for _, m := range msgs {
		fmt.Println(m)
		amsg := &api.Message{}
		err = mapstructure.Decode(m, amsg)
		if err != nil {
			return nil, errors.Wrap(err, "error while converting from repository message to message")
		}
		output = append(output, amsg)
	}
	fmt.Printf("\n%+v", output)
	return &api.Messages{Msg: output}, nil
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

		userSockets.Lock()
		userSockets.Us[ident.ID] = a
		userSockets.Unlock()
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
				logrus.Info("Event is Message")
				nm := &db.Message{
					Body:        m.Body,
					From:        ident.ID,
					To:          m.To,
					ChannelID:   m.Channel,
					MessageType: m.MessageType,
				}

				err = message.Handle(nm)
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
				jp := &join.JoinPayload{
					UserID:    ident.ID, //should get from jwt
					SessionId: j.SessionId,
				}

				join.Handle(a.Context(), jp, userSockets)

			}
		default:
			logrus.Infof("Type not matched : %+T", t)
		}
	}

}
