package cmd

import (
	"context"
	"time"

	"hermes/config"
	"hermes/pkg/db"
	"hermes/pkg/drivers/redis"
	"hermes/pkg/grpcserver"

	stan "github.com/nats-io/go-nats-streaming"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

var appContext = context.Background()

//Launch initalize needed things, Checks health of service by checking nats and db, and runs grpc server
func Launch() {
	logrus.Info("Loading C")
	logrus.Info("Processing Env")
	err := config.FromEnv(nil)
	if err != nil {
		logrus.Fatalf("error in loading config map from env: %v", err)
	}
	logrus.Info("Initiating DB package")
	db.Init()
	logrus.Info("Checking health")
	healthCheck()
	logrus.Info("Health check passed")
	grpcserver.CreateGRPCServer(appContext)
	logrus.Info("Initializing Hermes")

}

func healthCheck() {
	_, err := stan.Connect("test-cluster", "hermes-itself")
	if err != nil {
		logrus.Fatalf("Health Check failed : %v", err)
	}
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	logrus.Infof("Database URI is %v", config.C().Get("mongo_uri"))
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
}
