package chats

import (
	"aibot-backend/io"
	"aibot-backend/model"
	"time"

	"github.com/gin-gonic/gin"
)

func GetChatList(c *gin.Context, io *bot_io.Io) (*model.GetChatResp, error) {
	chatList, err := io.GetChatList(c)
	if err != nil {
		return nil, err
	}

	result := &model.GetChatResp{
		ChatList: chatList,
	}

	return result, nil
}

func CreateChat(c *gin.Context, io *bot_io.Io, req *model.CreateChatReq) (*model.CreateChatResp, error) {
	device, err := io.GetDeviceInfo(c, req.BotName)
	if err != nil {
		return nil, err
	}

	session, err := io.CreateChat(c, &model.BotChat{
		BotName:        req.BotName,
		NameCN:         device.NameCN,
		Image:          device.Image,
		Status:         model.NormalChatStatus,
		LastUpdateTime: time.Now().Unix(),
	})
	if err != nil {
		return nil, err
	}

	return &model.CreateChatResp{
		SessionID: session,
	}, nil
}

func GetChatConversations(c *gin.Context, io *bot_io.Io, session string) (*model.GetChatConversationsResp, error) {
	result, err := io.GetChatConversations(c, session)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func SaveChatConversations(c *gin.Context, io *bot_io.Io, req *model.SaveChatConversationsReq) (*model.SaveChatConversationsResp, error) {
	chat, err := io.GetChat(c, req.SessionID)
	if err != nil {
		return nil, err
	}

	var firstQuestion string
	for _, conversation := range req.ConversationList {
		if conversation.Role == model.UserRole {
			firstQuestion = conversation.Content
			break
		}
	}

	if len(chat.ConversationList) == 0 {
		chat.FirstQuestion = firstQuestion
	}

	chat.ConversationList = append(chat.ConversationList, req.ConversationList...)
	chat.LastUpdateTime = time.Now().Unix()

	err = io.UpdateChat(c, chat)
	if err != nil {
		return nil, err
	}

	return &model.SaveChatConversationsResp{BotName: chat.BotName}, nil
}

func DeleteChat(c *gin.Context, io *bot_io.Io, req *model.DeleteChatReq) error {
	return io.DeleteChat(c, req.SessionID)
}
