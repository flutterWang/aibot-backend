package bot_io

import (
	"aibot-backend/model"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
)

func (io *Io) GetKbLabelsList(c *gin.Context) ([]*model.Label, error) {
	result, err := io.MongoClient.Collection("labels").Find(c, bson.D{})
	if err != nil {
		return nil, err
	}

	var labels []*model.Label
	if err = result.All(c, &labels); err != nil {
		return nil, err
	}

	return labels, nil
}

func (io *Io) CreateLabel(c *gin.Context, req *model.CreateLabelReq) (string, error) {
	mongoClient := io.GetMongoClient()
	ID := uuid.New().String()

	item := &model.Label{
		Id:   ID,
		Name: req.Name,
	}

	_, err := mongoClient.Collection("labels").InsertOne(c, item)
	if err != nil {
		return "", err
	}

	return ID, nil
}

func (io *Io) DeleteLabel(c *gin.Context, labelID string) error {
	mongoClient := io.GetMongoClient()
	_, err := mongoClient.Collection("labels").DeleteOne(c, bson.M{"id": labelID})
	if err != nil {
		return err
	}

	return nil
}

func (io *Io) GetLabelKbs(c *gin.Context, labelID int64) ([]*model.Kb, error) {
	mongoClient := io.GetMongoClient()

	// 根据 label ID 获取满足条件的 Kb，Kb 的 labels 中包含了 label ID
	result, err := mongoClient.Collection("kb").Find(c, bson.M{"labels": bson.M{"$in": []int64{labelID}}})
	if err != nil {
		return nil, err
	}

	var kbs []*model.Kb
	if err = result.All(c, &kbs); err != nil {
		return nil, err
	}

	return kbs, nil
}
