package api

import (
	"aibot-backend/model"
	"aibot-backend/modules/categories"
	"github.com/gin-gonic/gin"
	"strconv"
)

func (s *server) handlerKbCategoriesList(c *gin.Context) {
	resp, err := categories.GetCategoryList(c, s.io)
	if err != nil {
		ginAbortWithCode(c, 500, err)
		return
	}
	ginJson(c, resp)
}

func (s *server) handlerCategoryCreate(c *gin.Context) {
	req := &model.CreateCategoryReq{}
	if err := c.ShouldBindJSON(req); err != nil {
		ginAbortWithCode(c, 400, err)
		return
	}

	result, err := categories.CreateCategory(c, s.io, req)
	if err != nil {
		ginAbortWithCode(c, 500, err)
		return
	}
	ginJson(c, &model.CreateCategoryResp{
		Id: result,
	})

}

func (s *server) handlerCategoryDelete(c *gin.Context) {
	categoryID := c.Query("category_id")

	if categoryID == "" {
		ginAbortWithCode(c, 400, nil)
		return
	}

	err := categories.DeleteCategory(c, s.io, categoryID)
	if err != nil {
		ginAbortWithCode(c, 500, err)
		return
	}
	ginJson(c, nil)
}

func (s *server) handlerCategoryKbs(c *gin.Context) {
	categoryID := c.Query("category_id")

	if categoryID == "" {
		ginAbortWithCode(c, 400, nil)
		return
	}

	categoryIDInt, err := strconv.ParseInt(categoryID, 10, 64)
	if err != nil {
		ginAbortWithCode(c, 400, err)
		return
	}

	result, err := categories.GetCategoryKbs(c, s.io, categoryIDInt)
	if err != nil {
		ginAbortWithCode(c, 500, err)
		return
	}
	ginJson(c, result)

}
