package api

import (
	"aibot-backend/model"
	"aibot-backend/modules/chats"
	"errors"

	"github.com/gin-gonic/gin"
)

func (s *server) handlerChatList(c *gin.Context) {
	resp, err := chats.GetChatList(c, s.io)
	if err != nil {
		ginAbortWithCode(c, 500, err)
		return
	}
	ginJson(c, resp)
}

func (s *server) handlerChatConversations(c *gin.Context) {
	sessionID := c.Query("session_id")
	if sessionID == "" {
		ginAbortWithCode(c, 400, errors.New("session_id is empty string"))
		return
	}

	resp, err := chats.GetChatConversations(c, s.io, sessionID)
	if err != nil {
		ginAbortWithCode(c, 500, err)
		return
	}

	ginJson(c, resp)
}

func (s *server) handlerCreateChat(c *gin.Context) {
	req := &model.CreateChatReq{}
	if err := c.ShouldBindJSON(req); err != nil {
		ginAbortWithCode(c, 400, err)
		return
	}

	result, err := chats.CreateChat(c, s.io, req)
	if err != nil {
		ginAbortWithCode(c, 500, err)
		return
	}
	ginJson(c, result)
}

func (s *server) handlerSaveChatConversations(c *gin.Context) {
	req := &model.SaveChatConversationsReq{}
	if err := c.ShouldBindJSON(req); err != nil {
		ginAbortWithCode(c, 400, err)
		return
	}

	resp, err := chats.SaveChatConversations(c, s.io, req)
	if err != nil {
		ginAbortWithCode(c, 500, err)
		return
	}
	ginJson(c, resp)
}

func (s *server) handlerChatDelete(c *gin.Context) {
	req := &model.DeleteChatReq{}
	if err := c.ShouldBindJSON(req); err != nil {
		ginAbortWithCode(c, 400, err)
		return
	}

	err := chats.DeleteChat(c, s.io, req)
	if err != nil {
		ginAbortWithCode(c, 500, err)
		return
	}
	ginJson(c, nil)
}
