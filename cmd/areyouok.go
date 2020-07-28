package cmd

import (
	"context"
	"github.com/hermeschat/engine/config"
	"time"

	"github.com/hermeschat/engine/subscription"
	"github.com/hermeschat/engine/subscription/nats"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

// areyouokCmd represents the areyouok command
var areyouokCmd = &cobra.Command{
	Use:   "areyouok",
	Short: "areyouok runs set of checks to make sure hermes is ok",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		config.Init()
		_, err := nats.Client("testing")
		if err != nil {
			monitoring.Logger().Fatalf("error in connecting to nats: %v", err)
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
		con, err := subscription.Redis()
		if err != nil {
			monitoring.Logger().Fatalf("could not connect redis:%v", err)
		}
		_, err = con.Ping().Result()
		if err != nil {
			monitoring.Logger().Fatalf("could not ping redis:%v", err)
		}
		return
	},
}

func init() {
	rootCmd.AddCommand(areyouokCmd)
	areyouokCmd.Aliases = []string{"youok", "ok"}
}
