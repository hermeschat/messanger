package main

// logrus
// pgk/errors
// gorilla/mux

import (
	"time"

	"git.raad.cloud/cloud/hermes/pkg/drivers/nats"
)

func main() {
	nats.Run()
	time.Sleep(time.Second * 20)
	// cmd.Launch("harchi")
}
