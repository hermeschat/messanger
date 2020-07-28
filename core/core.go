package core

import (
	"database/sql"
)

type ChatService interface{

}
type chatService struct {
	db *sql.DB//TODO: should fix
}
func NewChatService() (ChatService, error) {
	return ChatService(nil), nil
}