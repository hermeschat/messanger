package main

import (
	"context"
	"git.raad.cloud/cloud/hermes/pkg/api"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"time"
)

func main() {

	con, err := grpc.Dial("192.168.41.37:30041", grpc.WithInsecure())
	if err != nil {
		logrus.Fatalf("error : %v", err)
	}
	ctx := context.Background()
	cli := api.NewHermesClient(con)
	eventCli, err := cli.EventBuff(ctx)
	if err != nil {
		panic(err)
	}

	err = eventCli.SendMsg(&api.Event{
		Event: &api.Event_Read{
			Read: &api.ReadSignal{},
		},
	})
	if err != nil {
		panic(err)
	}
	time.Sleep(time.Second * 3)
	//_, err = cli.Echo(ctx, &api.Some{})
	//if err != nil {
	//	logrus.Fatalf("error : %v", err)
	//}
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
	//sid := "7c222aa1-68cd-4a84-b4d5-039941180323"
	//_, err = cli.Join(ctx, &api.JoinSignal{
	//	UserID: "amir",
	//	SessionId: sid,
	//})
	//if err != nil {
	//	panic(err)
	//}
	//msgCli, err := cli.NewMessage(context.Background())
	//if err != nil {
	//	panic(err)
	//}
	//cli.ListChannels()
	//err = msgCli.Send(&api.Message{
	//	MessageType: "1",
	//	From:"amir",
	//	To:"reza",
	//	Body: "hey",
	//})
	//if err != nil {
	//	panic(err)
	//}

	// cli.NewMessage(ctx)
}
