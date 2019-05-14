package main

import (
	"context"
	"os"

	"git.raad.cloud/cloud/hermes/pkg/api"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
)

func main() {
	if len(os.Args) < 3 {
		logrus.Fatal("Need a userid and sid")
	}
	con, err := grpc.Dial(":9044", grpc.WithInsecure())
	if err != nil {
		logrus.Fatalf("error : %v", err)
	}
	ctx := context.Background()
	cli := api.NewHermesClient(con)
	_, err = cli.Echo(ctx, &api.Some{})
	if err != nil {
		logrus.Fatalf("error : %v", err)
	}
	//resp, err := cli.CreateSession(ctx, &api.CreateSessionRequest{
	//	ClientType: "Proudly Windows",
	//	UserID:     os.Args[1],
	//})
	//if err != nil {
	//	logrus.Fatalf("error: %v", err)
	//}
	//
	//logrus.Println(resp.SessionID)
	//sid := resp.SessionID
	sid := os.Args[2]
	_, err = cli.Join(ctx, &api.JoinSignal{
		SessionId: sid,
	})
	if err != nil {
		panic(err)
	}
	cli.NewMessage(ctx, &api.Message{
		MessageType: "1",
		From:"amir",
		To:"reza",
		Body: "hey",
	})
	// cli.NewMessage(ctx)
}
