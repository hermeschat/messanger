package core

import (
	"context"
	"database/sql"
	"github.com/hermeschat/engine/db"
	"github.com/hermeschat/engine/transport/nats"
	"github.com/hermeschat/proto"
	"github.com/nats-io/stan.go"
)

type ChatService interface{
	NewMessage(ctx context.Context, msg *proto.Message) error
}
type chatService struct {
	nc stan.Conn
	db *sql.DB//TODO: should fix
	ps []Pusher
}
type Pusher interface {
	Push(to string, message *proto.Message) error
}
func NewChatService(pushers ...Pusher) (ChatService, error) {
	prv, err := db.NewSQLProvider()
	if err != nil {
	 	return nil, err
	}
	_db, err := prv.DB()
	if err != nil {
	 	return nil, err
	}
	nc, err := nats.Client()
	if err != nil {
		return nil, err
	}
	return &chatService{db: _db, nc: nc, ps: pushers}, nil
}