package api

import (
	"aibot-backend/io"
	"aibot-backend/model"

	"go.uber.org/zap"
)

type server struct {
	conf        *model.Config
	logger      *zap.Logger
	io          *bot_io.Io
	chatService *chatService
}

type chatService struct {
	chats map[string]*model.ChatServerInfo
}

func newServer() *server {
	return &server{
		conf:   &model.Config{},
		logger: NewLogger("logs/ai.log"),
		chatService: &chatService{
			chats: make(map[string]*model.ChatServerInfo),
		},
	}
}
