package knowledge_base

import (
	"aibot-backend/io"
	"log"

	"aibot-backend/model"
	"github.com/gin-gonic/gin"
)

func GetKbList(c *gin.Context, io *bot_io.Io) (*model.GetKbResp, error) {
	result, err := io.GetKbList(c)
	if err != nil {
		return nil, err
	}

	return &model.GetKbResp{
		Kbs: result,
	}, nil
}

func CreateCategory(c *gin.Context, io *bot_io.Io, req *model.CreateKbReq) (string, error) {
	result, err := io.CreateKb(c, req)
	if err != nil {
		return "", err
	}

	return result, nil
}

func GetKbDetail(c *gin.Context, io *bot_io.Io, kbID string) (*model.GetKbDetailResp, error) {
	result, err := io.GetKbDetail(c, kbID)
	if err != nil {
		return nil, err
	}

	return &model.GetKbDetailResp{
		Kb: result,
	}, nil
}

func UpdateKb(c *gin.Context, io *bot_io.Io, req *model.UpdateKbReq) error {
	kb := &model.Kb{
		Id:         req.ID,
		Name:       req.Name,
		CategoryID: req.CategoryID,
		Labels:     req.Labels,
		Desc:       req.Desc,
		Cover:      req.CoverPath,
	}

	return io.UpdateKb(c, kb)
}

func DeleteKb(c *gin.Context, io *bot_io.Io, kbID string) error {
	return io.DeleteKb(c, kbID)
}

func SearchKb(c *gin.Context, io *bot_io.Io, req *model.SearchKbReq) (*model.SearchKbResp, error) {
	result, err := io.SearchKb(c, req)
	if err != nil {
		return nil, err
	}

	return &model.SearchKbResp{
		Kbs: result,
	}, nil
}

func UploadDoc(c *gin.Context, io *bot_io.Io, req *model.UploadDocReq) error {
	kbData, err := io.GetKbDetail(c, req.KbID)
	if err != nil {
		return err
	}

	log.Println("kbData", kbData)

	for _, doc := range kbData.Docs {
		if doc.Name == req.Doc.Name {
			return nil
		}
	}

	kbData.Docs = append(kbData.Docs, req.Doc)
	err = io.UpdateKb(c, kbData)
	return nil
}
