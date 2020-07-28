package grpc

import (
	"context"
	"fmt"
	"github.com/hermeschat/engine/core"
	"github.com/hermeschat/proto"

	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpc_log "github.com/grpc-ecosystem/go-grpc-middleware/logging/zap"
	grpc_recovery "github.com/grpc-ecosystem/go-grpc-middleware/recovery"
	grpc_prom "github.com/grpc-ecosystem/go-grpc-prometheus"
	"github.com/hermeschat/engine/config"
	"github.com/hermeschat/engine/monitoring"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"

	"net"
)

//CreateGRPCServer creates a new grpc server
func CreateGRPCServer(ctx context.Context) {
	serverURL := config.C.GetString("app.url")
	monitoring.Logger().Infof("hermes GRPC server server is on %s", serverURL)
	lis, err := net.Listen("tcp", fmt.Sprintf("%s", serverURL))
	if err != nil {
		monitoring.Logger().Fatal("ERROR can't create a tcp listener ")
	}
	srv := grpc.NewServer(grpc.UnaryInterceptor(grpc_middleware.ChainUnaryServer(
		grpc_prom.UnaryServerInterceptor,
		grpc_log.UnaryServerInterceptor(monitoring.LoggerInstance, grpc_log.WithLevels(func(c codes.Code) zapcore.Level {
			if c > 0 {
				return zap.ErrorLevel
			}
			return zap.DebugLevel
		})),
		grpc_recovery.UnaryServerInterceptor(),
	)))
	hermesSrv, err := NewHermesServer()
	proto.RegisterHermesServer(srv, hermesSrv)
	monitoring.Logger().Info("Registering Hermes GRPC")
	err = srv.Serve(lis)
	if err != nil {
		monitoring.Logger().Fatal("ERROR in serving listener")
	}
}

type HermesServer struct {
	ChatService core.ChatService
}

func NewHermesServer() (*HermesServer, error) {
	chatService, err := core.NewChatService()
	if err != nil {
		return nil, err
	}
	return &HermesServer{ChatService: chatService}, nil
}
