package api

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"log"
	"net/http"
)

type baseRsp struct {
	Meta *Meta       `json:"meta"`
	Data interface{} `json:"data"`
}

type Meta struct {
	Code  int    `json:"code"`
	Msg   string `json:"msg"`
	Error string `json:"error"`
}

func ginJson(c *gin.Context, data interface{}) {
	c.JSON(200, &baseRsp{
		Meta: &Meta{Code: 0, Msg: "success"},
		Data: data,
	})

	addZapJsonField(c, "ginJson:", data)
}

func ginAbortWithCode(c *gin.Context, code int, err error) {
	c.Abort()
	log.Printf("gin abort, path:%s, code:%d, err:%s", c.Request.URL.Path, code, err)
	c.JSON(http.StatusOK, &baseRsp{
		Meta: &Meta{Code: code, Msg: "fail", Error: err.Error()},
		Data: nil,
	})

	c.Set("error", err)
	addZapField(c, zap.Error(err))
}

func addZapJsonField(c *gin.Context, key string, data interface{}) {
	if body, err := json.Marshal(data); err != nil {
		log.Printf("json marshal error:%s", err.Error())
	} else {
		value := string(body)
		addZapField(c, zap.String(key, value))
		// log.Printf("%s : %s", key, value)
	}
}

func addZapField(c *gin.Context, field zap.Field) {
	var zapFields []zap.Field
	zapFieldsData, ok := c.Get("zap_fields")
	if !ok {
		zapFields = make([]zap.Field, 0)
	} else {
		zapFields = zapFieldsData.([]zap.Field)
	}

	zapFields = append(zapFields, field)
	c.Set("zap_fields", zapFields)
}
