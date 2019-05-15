package cmd

import (
	"context"
	"net"

	"git.raad.cloud/cloud/hermes/pkg"
	"git.raad.cloud/cloud/hermes/pkg/api"
	"google.golang.org/grpc"

	"github.com/sirupsen/logrus"
)

var AppContext = context.Background()

func Launch(configPath string) {
	lis, err := net.Listen("tcp", ":9000")
	if err != nil {
		logrus.Fatal("ERROR can't create a tcp listener ")
	}
	logrus.Info("Initializing Hermes")

	//streamChain := grpc_middleware.ChainStreamServer(grpc_auth.StreamServerInterceptor(interceptor.UnaryAuthJWTInterceptor))
	//unaryChain := grpc_middleware.ChainUnaryServer(grpc_auth.UnaryServerInterceptor(interceptor.UnaryAuthJWTInterceptor))
	logrus.Info("Interceptors Created")
	//srv := grpc.NewServer(grpc.StreamInterceptor(streamChain), grpc.UnaryInterceptor(unaryChain))
	srv := grpc.NewServer()
	logrus.Info("Created New GRPC Server")
	hermes := pkg.HermesServer{}
	api.RegisterHermesServer(srv, hermes)
	logrus.Info("Registering Hermes RPCs")
	err = srv.Serve(lis)
	if err != nil {
		logrus.Fatal("ERROR in serving listener")
	}
	logrus.Info("We Are Live !!!!")
}
