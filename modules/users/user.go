package users

import (
	"aibot-backend/io"
	"aibot-backend/model"
	"aibot-backend/utils"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"io"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

var (
	errorNoAppInfo    = errors.New("no app id or secret")
	errorNoJWTConfig  = errors.New("no jwt config")
	errorNoUserOpenID = errors.New("no user open id")
)

const code2SessionURL = "https://api.weixin.qq.com/sns/jscode2session?appid=%s&secret=%s&js_code=%s&grant_type=authorization_code"

func getCode2Session(appID string, secret string, code string) (*model.Code2SessionResp, error) {
	url := fmt.Sprintf(code2SessionURL, appID, secret, code)

	res, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	resp := &model.Code2SessionResp{}
	err = json.Unmarshal(body, resp)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func GetUserInfo(c *gin.Context, io *bot_io.Io) (*model.UserInfo, error) {
	value, ok := c.Get("open_id")
	if !ok {
		return nil, errorNoUserOpenID
	}

	open_id, ok := value.(string)
	if !ok {
		return nil, errorNoUserOpenID
	}

	user, err := io.GetUserInfo(c, open_id)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func UserLogin(c *gin.Context, io *bot_io.Io, req *model.UserLoginReq) (*model.UserLoginResp, error) {
	appID, exists := c.Get("app_id")
	if !exists {
		return nil, errorNoAppInfo
	}

	appSecret, exists := c.Get("app_secret")
	if !exists {
		return nil, errorNoAppInfo
	}

	jwtSecret, exists := c.Get("jwt_secret")
	if !exists {
		return nil, errorNoJWTConfig
	}

	jwtExpiration, exists := c.Get("jwt_expiration")
	if !exists {
		return nil, errorNoJWTConfig
	}

	resp, err := getCode2Session(appID.(string), appSecret.(string), req.Code)
	if err != nil {
		return nil, err
	}

	var loginResp = &model.UserLoginResp{}
	user, err := io.GetUserInfo(c, resp.OpenID)
	if err != nil && !errors.Is(err, bot_io.ErrorUserNotExist) {
		return nil, err
	} else if err == nil {
		user.SessionKey = resp.SessionKey

		err := io.UpdateUser(c, user)
		if err != nil {
			return nil, err
		}
	} else {
		user = &model.UserInfo{
			OpenID:     resp.OpenID,
			SessionKey: resp.SessionKey,
			UnionID:    resp.UnionID,
		}
		err := io.CreateUser(c, &model.UserInfo{
			OpenID:     resp.OpenID,
			SessionKey: resp.SessionKey,
			UnionID:    resp.UnionID,
		})
		if err != nil {
			return nil, err
		}
	}

	if user.AvatarPath == "" || user.NickName == "" {
		loginResp.New = true
	} else {
		loginResp.New = false
	}

	claims := jwt.MapClaims{}
	claims["open_id"] = resp.OpenID
	token, err := utils.GenerateToken(claims, jwtSecret.(string), time.Duration(jwtExpiration.(int64))*time.Second)
	if err != nil {
		return nil, err
	}

	loginResp.Token = token
	loginResp.GetUserInfoResp = model.GetUserInfoResp{
		AvatarPath: user.AvatarPath,
		NickName:   user.NickName,
	}

	return loginResp, nil
}

func UserUpload(c *gin.Context, io *bot_io.Io, req *model.UserUploadReq) (*model.UserUploadResp, error) {
	user, err := GetUserInfo(c, io)
	if err != nil {
		return nil, err
	}

	user.AvatarPath = req.AvatarPath
	user.NickName = req.NickName

	err = io.UpdateUser(c, user)
	if err != nil {
		return nil, err
	}

	return &model.UserUploadResp{}, nil
}
