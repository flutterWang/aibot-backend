package api

import (
	"aibot-backend/model"
	"aibot-backend/modules/labels"
	"github.com/gin-gonic/gin"
	"strconv"
)

func (s *server) handlerKbLabelsList(c *gin.Context) {
	resp, err := labels.GetKbLabelsList(c, s.io)
	if err != nil {
		ginAbortWithCode(c, 500, err)
		return
	}
	ginJson(c, resp)
}

func (s *server) handlerLabelCreate(c *gin.Context) {
	req := &model.CreateLabelReq{}
	if err := c.ShouldBindJSON(req); err != nil {
		ginAbortWithCode(c, 400, err)
		return
	}

	result, err := labels.CreateLabel(c, s.io, req)
	if err != nil {
		ginAbortWithCode(c, 500, err)
		return
	}
	ginJson(c, &model.CreateLabelResp{
		Id: result,
	})
}

func (s *server) handlerLabelDelete(c *gin.Context) {
	labelID := c.Query("label_id")

	if labelID == "" {
		ginAbortWithCode(c, 400, nil)
		return
	}

	err := labels.DeleteLabel(c, s.io, labelID)
	if err != nil {
		ginAbortWithCode(c, 500, err)
		return
	}
	ginJson(c, nil)
}

func (s *server) handlerLabelKbs(c *gin.Context) {
	labelID := c.Query("label_id")

	if labelID == "" {
		ginAbortWithCode(c, 400, nil)
		return
	}

	labelIDInt, err := strconv.ParseInt(labelID, 10, 64)
	if err != nil {
		ginAbortWithCode(c, 400, err)
		return
	}

	resp, err := labels.GetLabelKbs(c, s.io, labelIDInt)
	if err != nil {
		ginAbortWithCode(c, 500, err)
		return
	}
	ginJson(c, resp)
}
