package bot_io

import (
	"aibot-backend/model"
	"errors"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

var (
	ErrorUserNotExist = errors.New("user does not exist")
)

func (io *Io) GetUserInfo(c *gin.Context, openID string) (*model.UserInfo, error) {
	var user model.UserInfo

	err := io.MongoClient.Collection("user").FindOne(c, bson.M{"open_id": openID}).Decode(&user)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			err = ErrorUserNotExist
		}

		return nil, err
	} else {
		return &user, nil
	}
}

func (io *Io) CreateUser(c *gin.Context, user *model.UserInfo) error {
	_, err := io.MongoClient.Collection("user").InsertOne(c, user)
	if err != nil {
		return err
	}

	return nil
}

func (io *Io) UpdateUser(c *gin.Context, user *model.UserInfo) error {
	mongoClient := io.GetMongoClient()
	_, err := mongoClient.Collection("user").UpdateOne(c, bson.M{"open_id": user.OpenID}, bson.M{"$set": user})
	if err != nil {
		return err
	}
	return nil
}
