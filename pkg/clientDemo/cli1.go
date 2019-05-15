package main

import (
	"context"
	"git.raad.cloud/cloud/hermes/pkg/api"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"time"
)

func main() {
	//con, err := grpc.Dial("https://chat.paygear.ir:443")

	con, err := grpc.Dial("localhost:9000", grpc.WithInsecure())
	if err != nil {
		logrus.Fatalf("error : %v", err)
	}
	ctx := context.Background()
	md := metadata.Pairs("Authorization", "bearer eyJ0eXAiOiJKV1QiLCJhbGciOiJIUzI1NiJ9.eyJhdWQiOiJiOWRjNzEyYzk1MmI0YWFmYjQ4MWFiZWRlMGZlYzRkOCIsImV4cCI6MTU1OTUzNTMxMSwibmJmIjoxNTU2OTQzMzExLCJpZCI6IjVjNGMyNjgzYmZkMDJhMmI5MjNhZjhiZSIsIm1lcmNoYW50X3JvbGVzIjp7IjViMmRmZTA0Y2YyNjU2MDAwYzk3YWFlNyI6WyJhZG1pbiJdfSwicm9sZSI6WyJ6ZXVzIl0sImFwcCI6IjU5YmVjM2ZhMGVjYTgxMDAwMWNlZWI4NiIsInN2YyI6eyJhY2NvdW50Ijp7InBlcm0iOjB9LCJjYXNoaWVyIjp7InBlcm0iOjB9LCJjYXNob3V0Ijp7InBlcm0iOjB9LCJjbHViIjp7InBlcm0iOjB9LCJjbHViX3NlcnZpY2UiOnsicGVybSI6MH0sImNvdXBvbiI6eyJwZXJtIjowfSwiY3JlZGl0Ijp7InBlcm0iOjB9LCJkZWxpdmVyeSI6eyJwZXJtIjowfSwiZXZlbnQiOnsicGVybSI6MH0sImZpbGUiOnsicGVybSI6MH0sImdhbWlmaWNhdGlvbiI6eyJwZXJtIjowfSwiZ2VvIjp7InBlcm0iOjB9LCJtZXNzYWdpbmciOnsicGVybSI6MH0sIm5vdGljZXMiOnsicGVybSI6MH0sInBheW1lbnQiOnsicGVybSI6MH0sInByb2R1Y3QiOnsicGVybSI6MH0sInB1c2giOnsicGVybSI6MH0sInFyIjp7InBlcm0iOjB9LCJzZWFyY2giOnsicGVybSI6MH0sInNvY2lhbCI6eyJwZXJtIjowfSwic3luYyI6eyJwZXJtIjowfSwidHJhbnNwb3J0Ijp7InBlcm0iOjB9fX0.kQ03S3UKh6JNbtEeAMg_e2Y2TnTJ6lPmFmc_5Tva6eY")
	ctx = metadata.NewOutgoingContext(ctx, md)
	cli := api.NewHermesClient(con)
	eventCli, err := cli.EventBuff(ctx)
	if err != nil {
		panic(err)
	}
	err = eventCli.Send(&api.Event{Event: &api.Event_Join{&api.JoinSignal{SessionId:"7c222aa1-68cd-4a84-b4d5-039941180323", UserID:"amir"}}})
	if err != nil {
		panic(err)
	}
	logrus.Info("Done")
	logrus.Info("Wait for any event")
	eventCli.
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
