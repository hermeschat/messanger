package main

import (
	"context"
	"fmt"
	"time"

	"git.raad.cloud/cloud/hermes/pkg/api"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

func main() {
	//con, err := grpc.Dial("https://chat.paygear.ir:443")

	con, err := grpc.Dial("localhost:9000", grpc.WithInsecure())
	//con, err := grpc.Dial("192.168.41.221:30041", grpc.WithInsecure())

	if err != nil {
		logrus.Fatalf("error : %v", err)
	}
	ctx, _ := context.WithTimeout(context.Background(), time.Hour*2)
	md := metadata.Pairs("Authorization", "bearer eyJ0eXAiOiJKV1QiLCJhbGciOiJIUzI1NiJ9.eyJhdWQiOiJiOWRjNzEyYzk1MmI0YWFmYjQ4MWFiZWRlMGZlYzRkOCIsImV4cCI6MTU5OTAwMDU5MywibmJmIjoxNTYzNDA4NTkzLCJpZCI6IjVjMGZhYmI2YmZkMDJhMmI5MjBjZDRjNiIsInVzZXJuYW1lIjoiMDkzODU3MDA3OTIiLCJhcHAiOiI1OWJlYzNmYTBlY2E4MTAwMDFjZWViODYiLCJzdmMiOnsiYWNjb3VudCI6eyJwZXJtIjowfSwiY2FzaGllciI6eyJwZXJtIjowfSwiY2FzaG91dCI6eyJwZXJtIjowfSwiY2x1YiI6eyJwZXJtIjowfSwiY2x1Yl9zZXJ2aWNlIjp7InBlcm0iOjB9LCJjb3Vwb24iOnsicGVybSI6MH0sImNyZWRpdCI6eyJwZXJtIjowfSwiZGVsaXZlcnkiOnsicGVybSI6MH0sImV2ZW50Ijp7InBlcm0iOjB9LCJmaWxlIjp7InBlcm0iOjB9LCJnYW1pZmljYXRpb24iOnsicGVybSI6MH0sImdlbyI6eyJwZXJtIjowfSwibWVzc2FnaW5nIjp7InBlcm0iOjB9LCJub3RpY2VzIjp7InBlcm0iOjB9LCJwYXltZW50Ijp7InBlcm0iOjB9LCJwcm9kdWN0Ijp7InBlcm0iOjB9LCJwdXNoIjp7InBlcm0iOjB9LCJxciI6eyJwZXJtIjowfSwic2VhcmNoIjp7InBlcm0iOjB9LCJzZXR0bGVtZW50Ijp7InBlcm0iOjB9LCJzb2NpYWwiOnsicGVybSI6MH0sInN5bmMiOnsicGVybSI6MH0sInRyYW5zcG9ydCI6eyJwZXJtIjowfSwid2FyZyI6eyJwZXJtIjowfX19.GKwwYO9Q-Q2cZLTe_Rr84r2qrt2y_VkhDdTCB-OXzmo")
	ctx = metadata.NewOutgoingContext(ctx, md)
	cli := api.NewHermesClient(con)
	//resp, err := cli.CreateSession(ctx, &api.CreateSessionRequest{
	//	ClientType: "Ubuntu",
	//	UserAgent:  "Terminal",
	//})
	//if err != nil {
	//	panic(err)
	//}
	//sid := resp.SessionID
	//logrus.Info(sid)
	msgs, err := cli.ListMessages(ctx, &api.Empty{})
	if err != nil {
		panic(err)
	}
	fmt.Println(msgs)
	return
	eventCli, err := cli.EventBuff(ctx)
	if err != nil {
		panic(err)
	}
	err = eventCli.Send(&api.Event{Event: &api.Event_Join{&api.JoinSignal{SessionId: "7c222aa1-68cd-4a84-b4d5-039941180323"}}})
	if err != nil {
		panic(err)
	}

	logrus.Info("Wait for any event")
	for {
		ev, err := eventCli.Recv()
		if err != nil {
			continue
		}
		switch ev.GetEvent().(type) {
		case *api.Event_Read:
			logrus.Info("Message has been read")
		case *api.Event_NewMessage:
			logrus.Info("New Message recieved")
			m := ev.GetNewMessage()
			logrus.Infof("%+v", m)
		case *api.Event_Dlv:
			logrus.Info("Message delivered")
		}
	}
	//emp, err := cli.Echo(ctx, &api.Empty{})
	//if err != nil {
	//	panic(err)
	//}
	//logrus.Infof("status is %v", emp.Status)
	//err = eventCli.SendMsg(&api.Event{
	//	Event: &api.Event_Read{
	//		Read: &api.ReadSignal{},
	//	},
	//})
	//if err != nil {
	//	panic(err)
	//}
	time.Sleep(time.Second * 100)
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
