package bot_io

import (
	"aibot-backend/model"
	"github.com/google/uuid"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
)

func (io *Io) GetKbList(c *gin.Context) ([]*model.Kb, error) {
	mongoClient := io.GetMongoClient()

	result, err := mongoClient.Collection("kb").Find(c, bson.D{})
	if err != nil {
		return nil, err
	}

	var kbs []*model.Kb
	if err = result.All(c, &kbs); err != nil {
		return nil, err
	}
	return kbs, nil
}

func (io *Io) CreateKb(c *gin.Context, req *model.CreateKbReq) (string, error) {
	mongoClient := io.GetMongoClient()

	id := uuid.New().String()

	item := &model.Kb{
		Id:         id,
		Name:       req.Name,
		CategoryID: req.CategoryID,
		Labels:     req.Labels,
		Cover:      req.CoverPath,
		CreateAt:   time.Now().Unix(),
		Desc:       req.Desc,
		Docs:       make([]*model.DocInfo, 0),
	}

	_, err := mongoClient.Collection("kb").InsertOne(c, item)
	if err != nil {
		return "", err
	}

	return id, nil
}

func (io *Io) GetKbDetail(c *gin.Context, kbID string) (*model.Kb, error) {
	mongoClient := io.GetMongoClient()

	var kb model.Kb
	err := mongoClient.Collection("kb").FindOne(c, bson.M{"id": kbID}).Decode(&kb)
	if err != nil {
		return nil, err
	}

	return &kb, nil
}

func (io *Io) UpdateKb(c *gin.Context, kb *model.Kb) error {
	mongoClient := io.GetMongoClient()

	_, err := mongoClient.Collection("kb").UpdateOne(c, bson.M{"id": kb.Id}, bson.M{"$set": kb})
	if err != nil {
		return err
	}

	return nil
}

func (io *Io) DeleteKb(c *gin.Context, kbID string) error {
	mongoClient := io.GetMongoClient()

	_, err := mongoClient.Collection("kb").UpdateOne(c, bson.M{"_id": kbID}, bson.M{"$set": bson.M{"delete_at": time.Now().Unix()}})
	if err != nil {
		return err
	}

	return nil

}

func (io *Io) SearchKb(c *gin.Context, req *model.SearchKbReq) ([]*model.Kb, error) {
	mongoClient := io.GetMongoClient()

	result, err := mongoClient.Collection("kb").Find(c, bson.M{"name": bson.M{"$regex": req.Keyword}})
	if err != nil {
		return nil, err
	}

	var kbs []*model.Kb
	if err = result.All(c, &kbs); err != nil {
		return nil, err
	}

	return kbs, nil
}
