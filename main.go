package main

import (
	"github.com/hermeschat/engine/cmd"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	go func() {
		sigs := make(chan os.Signal)
		signal.Notify(sigs, syscall.SIGTERM, syscall.SIGINT)
		for range sigs {
			os.Exit(0)
		}
	}()
	cmd.Execute()
}
