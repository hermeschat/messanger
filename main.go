package main

import (
	"os"
	"os/signal"
	"syscall"

	"hermes/cmd"
	"hermes/subscription"
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
