package api

import (
	"aibot-backend/model"
	"aibot-backend/modules/users"
	"path"

	"github.com/gin-gonic/gin"
)

func (s *server) handlerGetUserInfo(c *gin.Context) {
	user, err := users.GetUserInfo(c, s.io)
	if err != nil {
		ginAbortWithCode(c, 500, err)
		return
	}

	result := model.GetUserInfoResp{
		AvatarPath: user.AvatarPath,
		NickName:   user.NickName,
	}
	ginJson(c, result)
}

func (s *server) handlerUserLogin(c *gin.Context) {
	req := &model.UserLoginReq{}
	if err := c.ShouldBindJSON(req); err != nil {
		ginAbortWithCode(c, 400, err)
		return
	}

	c.Set("jwt_secret", s.conf.JWTSecret)
	c.Set("jwt_expiration", s.conf.JWTExpiration)
	c.Set("app_id", s.conf.AppID)
	c.Set("app_secret", s.conf.AppSecret)

	result, err := users.UserLogin(c, s.io, req)
	if err != nil {
		ginAbortWithCode(c, 500, err)
		return
	}
	ginJson(c, result)
}

func (s *server) handlerUserUpload(c *gin.Context) {
	req := &model.UserUploadReq{}

	file, err := c.FormFile("file")
	if err != nil {
		ginAbortWithCode(c, 400, err)
		return
	}

	req.AvatarPath = file.Filename
	req.NickName = c.PostForm("nick_name")

	err = c.SaveUploadedFile(file, path.Join("./static", file.Filename))
	if err != nil {
		ginAbortWithCode(c, 500, err)
		return
	}

	result, err := users.UserUpload(c, s.io, req)
	if err != nil {
		ginAbortWithCode(c, 500, err)
		return
	}

	ginJson(c, result)
}
