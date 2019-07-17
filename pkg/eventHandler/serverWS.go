package eventHandler

import (
	"log"
	"net/http"
)

func Serve() {
	BaseHub = NewHub()
	go BaseHub.Run()
	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		ServeWs(BaseHub, w, r)
	})
	err := http.ListenAndServe(":9000", nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
