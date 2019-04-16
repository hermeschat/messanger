package main

// logrus
// pgk/errors
// gorilla/mux

import (
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
)

func main() {
	log.WithFields(log.Fields{
		"service": "hermes",
	}).Info("A walrus appears")

	errors.Wrap(errors.New("ok"), "ok")
}
