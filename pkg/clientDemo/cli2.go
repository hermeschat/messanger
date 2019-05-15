package main

import (
	"context"
	"git.raad.cloud/cloud/hermes/pkg/api"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
)

func main2() {

	con, err := grpc.Dial(":9044", grpc.WithInsecure())
	if err != nil {
		logrus.Fatalf("error : %v", err)
	}
	ctx := context.Background()
	cli := api.NewHermesClient(con)
	_ = ctx
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
	//userid:= "reza"
	//sid := "86ada2a3-f70f-40bb-a702-f2a6598000b4"
	//_, err = cli.Join(ctx, &api.JoinSignal{
	//	UserID: userid,
	//	SessionId: sid,
	//})
	if err != nil {
		panic(err)
	}
	msgCli, err := cli.EventBuff(context.Background())
	m := &api.Message{}
	err = msgCli.RecvMsg(m)
	if err != nil {
		panic(err)
	}
	logrus.Infof("Message recieved : %v", m)
	logrus.Info("Done")
	//cli.NewMessage(ctx, &api.Message{
	//	MessageType: "1",
	//	From:"amir",
	//	To:"reza",
	//	Body: "hey",
	//})
	// cli.NewMessage(ctx)
}


// join ro bokon event
// bad vaghti join befreste socketesh ro to factroy negah dar
// tamoom