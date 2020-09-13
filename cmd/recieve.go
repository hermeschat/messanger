package cmd

import (
	"context"
	"fmt"
	"github.com/hermeschat/engine/config"
	"github.com/hermeschat/proto"
	"github.com/spf13/cobra"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"os"
	"os/signal"
	"syscall"
	"time"
)

var receiveCmd = &cobra.Command{
	Use:   "receiver",
	Short: "hermes-cli receiver mode",
	Long: `in receiver mode you can receive messages
	usage:
		hermes-cli receive`,
	Run: func(cmd *cobra.Command, args []string) {
		config.Init()
		sigs := make(chan os.Signal)
		signal.Notify(sigs, syscall.SIGTERM)
		con, err := grpc.Dial(fmt.Sprintf("%s:%s", config.C.Get("host"), config.C.Get("port")), grpc.WithInsecure())
		if err != nil {
			fmt.Fprintf(os.Stdout, "error in grpc dial: %v", err)
			os.Exit(1)
		}
		fmt.Println("Waiting for any message")
		cli := proto.NewHermesClient(con)
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()
		md := metadata.Pairs("Authorization", fmt.Sprintf("%s", config.C.Get("receiver_token")))
		ctx = metadata.NewOutgoingContext(ctx, md)
		buff, err := cli.EventBuff(ctx)
		if err != nil {
			fmt.Fprintf(os.Stdout, "error in calling event buff: %v", err)
			os.Exit(1)
		}
		go func() {
			for {
				e, err := buff.Recv()
				if err != nil {
					fmt.Fprintf(os.Stderr, "error in receiving event: %v", err)
					time.Sleep(time.Second * 3)
					continue
				}
				fmt.Println("event is")
				switch e.GetEvent().(type) {
				case *proto.Event_NewMessage:
					fmt.Println("New Message recieved")
					m := e.GetNewMessage()
					fmt.Printf("%+v\n", m)
				}
			}
		}()
		<-sigs

	},
}

func init() {

	rootCmd.AddCommand(receiveCmd)
	receiveCmd.Aliases = []string{"recv"}

}