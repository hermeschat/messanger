package cmd

import (
	"context"
	"crypto/tls"
	"git.raad.cloud/cloud/hermes/pkg"
	"git.raad.cloud/cloud/hermes/pkg/api"
	"git.raad.cloud/cloud/hermes/pkg/interceptor"
	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpc_auth "github.com/grpc-ecosystem/go-grpc-middleware/auth"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"net"

	"github.com/sirupsen/logrus"
)

var AppContext = context.Background()

func Launch(configPath string) {
	lis, err := net.Listen("tcp", "localhost:9044")
	if err != nil {
		logrus.Fatal("ERROR can't create a tcp listener ")
	}
	logrus.Infof(" Initializing Hermes ...")
	creds := credentials.NewTLS(&tls.Config{
		// TLS config values here
	})

	tlsServerOption := grpc.Creds(creds)
	_ = tlsServerOption
	streamChain := grpc_middleware.ChainStreamServer(grpc_auth.StreamServerInterceptor(interceptor.UnaryAuthJWTInterceptor))
	unaryChain := grpc_middleware.ChainUnaryServer(grpc_auth.UnaryServerInterceptor(interceptor.UnaryAuthJWTInterceptor))

	srv := grpc.NewServer(grpc.StreamInterceptor(streamChain), grpc.UnaryInterceptor(unaryChain), tlsServerOption)
	hermes := pkg.HermesServer{}
	api.RegisterHermesServer(srv, hermes)
	err = srv.Serve(lis)
	if err != nil {
		logrus.Fatal("ERROR in serving listener")
	}

}
