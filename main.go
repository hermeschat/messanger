package main

// pgk/errors
// gorilla/mux

import (
	"git.raad.cloud/cloud/hermes/cmd"
	"github.com/sirupsen/logrus"
)

func main() {
	// nats.Run()
	logrus.Info("Hi. Im hermes")

	cmd.Launch("")
}
