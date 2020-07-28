package cmd

import (
	"context"
	"github.com/hermeschat/engine/config"
	"github.com/hermeschat/engine/monitoring"
	nats "github.com/hermeschat/engine/transport/nats"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"time"
)

// Health Check
var areyouokCmd = &cobra.Command{
	Use:   "areyouok",
	Short: "areyouok runs set of checks to make sure hermes is ok",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		config.Init()
		isConnectd := nats.HealthCheck()
		if isConnectd != false {
			monitoring.Logger().Fatalf("error in connecting to nats")
		}

		ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
		client, err := mongo.Connect(ctx, options.Client().ApplyURI(config.MongoURI()))
		if err != nil {
			monitoring.Logger().Fatalf(errors.Wrap(err, "can't connect to mongodb FUCK").Error())
		}
		err = client.Ping(ctx, readpref.Primary())
		if err != nil {
			monitoring.Logger().Fatalf("could not ping db")
		}
		//con, err := subscription.Redis()
		//if err != nil {
		//	monitoring.Logger().Fatalf("could not connect redis:%v", err)
		//}
		//_, err = con.Ping().Result()
		//if err != nil {
		//	monitoring.Logger().Fatalf("could not ping redis:%v", err)
		//}
		return
	},
}

func init() {
	rootCmd.AddCommand(areyouokCmd)
	areyouokCmd.Aliases = []string{"youok", "ok"}
}
