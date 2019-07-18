package eventHandler

import (
	"fmt"
	"git.raad.cloud/cloud/hermes/pkg/auth"
	"log"
	"net/http"
)

func Serve() {
	BaseHub = NewHub()
	go BaseHub.Run()
	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		ServeWs(BaseHub, w, r)
	})
	http.HandleFunc("/messages", func(w http.ResponseWriter, r *http.Request) {
		token := r.Header.Get("authorization")
		ident, err := auth.GetAuthentication(token, "")
		if err != nil {
			fmt.Fprintf(w, "cannot authenticate with given authorization header")
			return
		}
		_ = ident
		//get messages of this user from mongo
	})
	http.HandleFunc("/channels", func(w http.ResponseWriter, r *http.Request) {
		token := r.Header.Get("authorization")
		ident, err := auth.GetAuthentication(token, "")
		if err != nil {
			fmt.Fprintf(w, "cannot authenticate with given authorization header")
			return
		}
		_ = ident
		//get channels of this user from mongo
	})
	err := http.ListenAndServe(":9000", nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
