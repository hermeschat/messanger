package core

import (
	"database/sql"
	"github.com/nats-io/stan.go"
)

type ChatService interface{

}
type chatService struct {
	nc stan.Conn
	db *sql.DB//TODO: should fix
	ps []Pusher
}
type Pusher interface {
	Push(data []byte) error
}
func NewChatService(pushers ...Pusher) (ChatService, error) {
	return ChatService(nil), nil
}