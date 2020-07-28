package main
//go:generate sqlboiler --wipe psql -p repository
import (
	"github.com/hermeschat/engine/cmd"
	"github.com/hermeschat/engine/config"
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
	err := config.Init()
	if err != nil {
	 	panic(err)
	}
	cmd.Execute()
}
