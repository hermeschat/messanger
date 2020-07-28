package core

type ChatService interface{

}

func NewChatService() (ChatService, error) {
	return ChatService(nil), nil
}