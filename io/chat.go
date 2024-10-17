package bot_io

import (
	"aibot-backend/model"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func (io *Io) GetChatConversations(c *gin.Context, sessionID string) (*model.GetChatConversationsResp, error) {
	var chat model.BotChat

	err := io.MongoClient.Collection("chat").FindOne(c, bson.D{
		{Key: "session_id", Value: sessionID},
		{Key: "status", Value: model.NormalChatStatus},
	}).Decode(&chat)
	if err != nil {
		return nil, err
	}

	resp := &model.GetChatConversationsResp{
		ConversationList: chat.ConversationList,
		SessionID:        chat.SessionID,
	}

	return resp, nil
}

func (io *Io) GetChatList(c *gin.Context) ([]*model.BotChat, error) {
	var chat []*model.BotChat

	opts := options.Find().SetSort(bson.D{{Key: "last_update_time", Value: -1}})
	result, err := io.MongoClient.Collection("chat").Find(c, bson.D{
		{Key: "status", Value: model.NormalChatStatus},
	}, opts)
	if err != nil {
		return nil, err
	}

	err = result.All(c, &chat)
	if err != nil {
		return nil, err
	}

	return chat, err
}

func (io *Io) GetChat(c *gin.Context, sessionID string) (*model.BotChat, error) {
	var chat model.BotChat
	err := io.MongoClient.Collection("chat").FindOne(c, bson.D{
		{Key: "session_id", Value: sessionID},
	}).Decode(&chat)
	if err != nil {
		return nil, err
	} else {
		return &chat, nil
	}
}

func (s *Io) CreateChat(c *gin.Context, chat *model.BotChat) (string, error) {
	chat.SessionID = uuid.New().String()
	_, err := s.MongoClient.Collection("chat").InsertOne(c, chat)
	if err != nil {
		return "", err
	}

	return chat.SessionID, nil
}

func (ai_Io *Io) UpdateChat(c *gin.Context, chat *model.BotChat) error {
	mongoClient := ai_Io.GetMongoClient()
	_, err := mongoClient.Collection("chat").UpdateOne(c, bson.M{
		"session_id": chat.SessionID,
	}, bson.M{"$set": chat})
	if err != nil {
		return err
	}
	return nil
}

func (ai_Io *Io) DeleteChat(c *gin.Context, sessionID string) error {
	_, err := ai_Io.MongoClient.Collection("chat").DeleteOne(c, bson.M{"session_id": sessionID})
	return err
}
