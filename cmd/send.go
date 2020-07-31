package cmd

import "C"
import (
	"context"
	"fmt"
	"github.com/hermeschat/engine/config"
	"github.com/hermeschat/proto"
	"github.com/spf13/cobra"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"os"
	"strconv"
	"time"
)

var sendCmd = &cobra.Command{
	Use:   "send",
	Short: "hermes-cli send mode",
	Long: `in send mode you can send messages
	usage:
		hermes-cli send [receiver] [body]`,
	Run: func(cmd *cobra.Command, args []string) {
		config.Init()
		if len(args) < 2 {
			fmt.Fprintf(os.Stderr, "send needs exactly two arguments")
			os.Exit(1)
		}
		argOne := args[0]
		receiverID, err := strconv.Atoi(argOne)
		if err != nil {
			fmt.Printf("reciever id must be integer , error : %v", err)
			os.Exit(1)
		}
		msgBody := args[1]
		con, err := grpc.Dial(fmt.Sprintf("%s:%s", config.C.Get("host"), config.C.Get("port")), grpc.WithInsecure())
		if err != nil {
			fmt.Fprintf(os.Stderr, "error in grpc dial: %v", err)
			os.Exit(1)
		}

		cli := proto.NewHermesClient(con)
		ctx, cancel := context.WithCancel(context.Background())
		md := metadata.Pairs("Authorization", fmt.Sprintf("%s", config.C.Get("sender_token")))
		ctx = metadata.NewOutgoingContext(ctx, md)
		defer cancel()
		buff, err := cli.EventBuff(ctx)
		if err != nil {
			fmt.Fprintf(os.Stderr, "error in calling event buff: %v", err)
			os.Exit(1)
		}

		err = buff.Send(&proto.Event{Event: &proto.Event_NewMessage{NewMessage: &proto.Message{
			To:   int64(receiverID),
			Body: msgBody,
		}}})
		if err != nil {
			fmt.Fprintf(os.Stderr, "error in sending message:%v", err)
			os.Exit(1)

		}
		time.Sleep(time.Hour * 2)
		fmt.Println("message sent")
	},
}

func init() {
	rootCmd.AddCommand(sendCmd)

}