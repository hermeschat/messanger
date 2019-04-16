package cmd

import (
	"context"

	"github.com/sirupsen/logrus"
)

var AppContext = context.Background()

func Launch(configPath string) {
	logrus.Infof(" Initializing Hermes ...")

}
