package main

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/hermeschat/engine/cmd"
	"github.com/hermeschat/engine/subscription"
)

func main() {
	go func() {
		sigs := make(chan os.Signal)
		signal.Notify(sigs, syscall.SIGTERM, syscall.SIGINT)
		for range sigs {
			subscription.Clean()
			os.Exit(0)
		}
	}()
	cmd.Execute()
}
