package api

import (
	"aibot-backend/model"
	"aibot-backend/modules/knowledge_base"
	"github.com/gin-gonic/gin"
)

func (s *server) handlerKbList(c *gin.Context) {
	resp, err := knowledge_base.GetKbList(c, s.io)
	if err != nil {
		ginAbortWithCode(c, 500, err)
		return
	}
	ginJson(c, resp)
}

func (s *server) handlerKbCreate(c *gin.Context) {
	req := &model.CreateKbReq{}
	if err := c.ShouldBindJSON(req); err != nil {
		ginAbortWithCode(c, 400, err)
		return
	}

	result, err := knowledge_base.CreateCategory(c, s.io, req)
	if err != nil {
		ginAbortWithCode(c, 500, err)
		return
	}
	ginJson(c, &model.CreateKbResp{
		Id: result,
	})
}

func (s *server) handlerKbDetail(c *gin.Context) {
	kbID := c.Query("kb_id")

	if kbID == "" {
		ginAbortWithCode(c, 400, nil)
		return
	}
	result, err := knowledge_base.GetKbDetail(c, s.io, kbID)
	if err != nil {
		ginAbortWithCode(c, 500, err)
		return
	}
	ginJson(c, result)
}

func (s *server) handlerKbUpdate(c *gin.Context) {
	req := &model.UpdateKbReq{}
	if err := c.ShouldBindJSON(req); err != nil {
		ginAbortWithCode(c, 400, err)
		return
	}

	err := knowledge_base.UpdateKb(c, s.io, req)
	if err != nil {
		ginAbortWithCode(c, 500, err)
		return
	}
	ginJson(c, nil)
}

func (s *server) handlerKbDelete(c *gin.Context) {
	kbID := c.Query("kb_id")
	err := knowledge_base.DeleteKb(c, s.io, kbID)
	if err != nil {
		ginAbortWithCode(c, 500, err)
		return
	}
	ginJson(c, nil)
}

func (s *server) handlerKbSearch(c *gin.Context) {
	req := &model.SearchKbReq{}
	if err := c.ShouldBindJSON(req); err != nil {
		ginAbortWithCode(c, 400, err)
		return
	}

	result, err := knowledge_base.SearchKb(c, s.io, req)
	if err != nil {
		ginAbortWithCode(c, 500, err)
		return
	}
	ginJson(c, result)
}

func (s *server) handlerKbUploadDoc(c *gin.Context) {
	req := &model.UploadDocReq{}
	if err := c.ShouldBindJSON(req); err != nil {
		ginAbortWithCode(c, 400, err)
		return
	}
	err := knowledge_base.UploadDoc(c, s.io, req)
	if err != nil {
		ginAbortWithCode(c, 500, err)
		return
	}
	ginJson(c, nil)
}
