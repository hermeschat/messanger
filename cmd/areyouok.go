package cmd

import (
	"context"
	"time"

	"github.com/amirrezaask/config"
	stan "github.com/nats-io/go-nats-streaming"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"hermes/pkg/drivers/redis"
)

// areyouokCmd represents the areyouok command
var areyouokCmd = &cobra.Command{
	Use:   "areyouok",
	Short: "areyouok runs set of checks to make sure hermes is ok",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		_, err := stan.Connect("test-cluster", "hermes-itself")
		if err != nil {
			logrus.Fatalf("Health Check failed : %v", err)
		}
		ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
		client, err := mongo.Connect(ctx, options.Client().ApplyURI(config.C().Get("mongo_uri")))
		if err != nil {
			logrus.Fatalf(errors.Wrap(err, "can't connect to mongodb FUCK").Error())
		}
		err = client.Ping(ctx, readpref.Primary())
		if err != nil {
			logrus.Fatalf("could not ping db")
		}
		con, err := redis.ConnectRedis()
		if err != nil {
			logrus.Fatalf("could not connect redis:%v", err)
		}
		_, err = con.Ping().Result()
		if err != nil {
			logrus.Fatalf("could not ping redis:%v", err)
		}
		return
	},
}

func init() {
	rootCmd.AddCommand(areyouokCmd)
	areyouokCmd.Aliases = []string{"youok", "ok"}
}
