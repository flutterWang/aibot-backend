package bot_io

import (
	"aibot-backend/model"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
)

func (io *Io) GetCategoryList(c *gin.Context) ([]*model.CategoryInfo, error) {
	var categories []*model.CategoryInfo

	result, err := io.MongoClient.Collection("category").Find(c, bson.D{})
	if err != nil {
		return nil, err
	}

	err = result.All(c, &categories)
	if err != nil {
		return nil, err
	}

	return categories, err
}

func (s *Io) CreateCategory(c *gin.Context, req *model.CreateCategoryReq) (string, error) {
	id := uuid.New().String()

	item := &model.CategoryInfo{
		Name: req.Name,
		Id:   id,
	}

	_, err := s.MongoClient.Collection("category").InsertOne(c, item)
	if err != nil {
		return "", err
	}

	return id, nil
}

func (s *Io) DeleteCategory(c *gin.Context, categoryID string) error {
	_, err := s.MongoClient.Collection("category").DeleteOne(c, bson.M{"id": categoryID})
	return err
}

func (s *Io) GetCategoryKbs(c *gin.Context, categoryID int64) ([]*model.Kb, error) {
	var kbs []*model.Kb

	// 根据 categoryID 查找 kb, categoryID 在 kb 的 category 结构体的 id 中
	result, err := s.MongoClient.Collection("kb").Find(c, bson.M{"category.id": categoryID})
	if err != nil {
		return nil, err
	}

	err = result.All(c, &kbs)
	if err != nil {
		return nil, err
	}

	return kbs, nil
}
