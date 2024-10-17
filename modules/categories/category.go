package categories

import (
	"aibot-backend/io"
	"aibot-backend/model"

	"github.com/gin-gonic/gin"
)

func GetCategoryList(c *gin.Context, io *bot_io.Io) (*model.GetCategoryListResp, error) {
	result, err := io.GetCategoryList(c)
	if err != nil {
		return nil, err
	}

	return &model.GetCategoryListResp{
		Categories: result,
	}, nil

}

func CreateCategory(c *gin.Context, io *bot_io.Io, req *model.CreateCategoryReq) (string, error) {
	result, err := io.CreateCategory(c, req)
	if err != nil {
		return "", err
	}

	return result, nil

}

func DeleteCategory(c *gin.Context, io *bot_io.Io, categoryID string) error {
	return io.DeleteCategory(c, categoryID)
}

func GetCategoryKbs(c *gin.Context, io *bot_io.Io, categoryID int64) (*model.GetCategoryKbsResp, error) {
	result, err := io.GetCategoryKbs(c, categoryID)
	if err != nil {
		return nil, err
	}

	return &model.GetCategoryKbsResp{
		Kbs: result,
	}, nil

}
