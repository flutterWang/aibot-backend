package labels

import (
	"aibot-backend/io"
	"aibot-backend/model"
	"github.com/gin-gonic/gin"
)

func GetKbLabelsList(c *gin.Context, io *bot_io.Io) (*model.GetKbLabelsListResp, error) {
	result, err := io.GetKbLabelsList(c)
	if err != nil {
		return nil, err
	}

	return &model.GetKbLabelsListResp{
		Labels: result,
	}, nil

}

func CreateLabel(c *gin.Context, io *bot_io.Io, req *model.CreateLabelReq) (string, error) {
	result, err := io.CreateLabel(c, req)
	if err != nil {
		return "", err
	}

	return result, nil
}

func DeleteLabel(c *gin.Context, io *bot_io.Io, labelID string) error {
	return io.DeleteLabel(c, labelID)
}

func GetLabelKbs(c *gin.Context, io *bot_io.Io, labelID int64) (*model.GetLabelKbsResp, error) {
	result, err := io.GetLabelKbs(c, labelID)
	if err != nil {
		return nil, err
	}

	return &model.GetLabelKbsResp{
		Kbs: result,
	}, nil

}
