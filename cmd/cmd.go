package cmd

import (
	"context"
	"git.raad.cloud/cloud/hermes/pkg"
	"git.raad.cloud/cloud/hermes/pkg/api"
	"google.golang.org/grpc"
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
	srv := grpc.NewServer()
	hermes := pkg.HermesServer{}
	api.RegisterHermesServer(srv, hermes)
	err = srv.Serve(lis)
	if err != nil {
		logrus.Fatal("ERROR in serving listener")
	}

}
