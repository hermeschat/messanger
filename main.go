package main

import (
	"git.raad.cloud/cloud/hermes/cmd"
	"github.com/sirupsen/logrus"
)

func main() {

	cmd.Launch("")
	defer func() {
		if err := recover(); err != nil {
			logrus.Errorf("panic happend:\n%v", err)

		}
	}()
}
