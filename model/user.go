package model

type Code2SessionResp struct {
	OpenID     string `json:"openid"`
	SessionKey string `json:"session_key"`
	UnionID    string `json:"unionid"`

	ErrCode int64  `json:"errcode"`
	ErrMsg  string `json:"errmsg"`
}

type UserLoginReq struct {
	Code string `json:"code"`
}

type UserLoginResp struct {
	New   bool   `json:"new"` // 是否新用户
	Token string `json:"token"`
	GetUserInfoResp
}

type UserInfo struct {
	OpenID     string `json:"open_id" bson:"open_id"`
	SessionKey string `json:"session_key" bson:"session_key"`
	UnionID    string `json:"union_id" bson:"union_id"`
	AvatarPath string `json:"avatar_path" bson:"avatar_path"`
	NickName   string `json:"nick_name" bson:"nick_name"`
}

type GetUserInfoResp struct {
	AvatarPath string `json:"avatar_path" bson:"avatar_path"`
	NickName   string `json:"nick_name" bson:"nick_name"`
}

type UserUploadReq struct {
	AvatarPath string `json:"avatar_path" bson:"avatar_path"`
	NickName   string `json:"nick_name" bson:"nick_name"`
}

type UserUploadResp struct {
	AvatarPath string `json:"avatar_path"`
}
