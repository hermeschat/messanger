package grpc

import (
	"context"
	"fmt"
	"github.com/hermeschat/proto"
	"google.golang.org/grpc"
	"github.com/hermeschat/engine/config"
	"github.com/hermeschat/engine/monitoring"
	"net"
)

//CreateGRPCServer creates a new grpc server
func CreateGRPCServer(ctx context.Context) {
	serverURL, err := config.C.GetString("app.url")
	monitoring.Logger().Infof("hermes GRPC server server is on %s", serverURL)
	lis, err := net.Listen("tcp", fmt.Sprintf("%s", serverURL))
	if err != nil {
		monitoring.Logger().Fatal("ERROR can't create a tcp listener ")
	}


	srv := grpc.NewServer()
	monitoring.Logger().Info("Registering Hermes GRPC")
	err = srv.Serve(lis)
	if err != nil {
		monitoring.Logger().Fatal("ERROR in serving listener")
	}

	monitoring.Logger().Info("GRPC is Live !!!")
}
