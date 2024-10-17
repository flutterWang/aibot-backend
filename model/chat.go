package model

const (
	// chat status
	NormalChatStatus  = 1
	DeletedChatStatus = 2

	// conversation type
	TextConversationType     = 1
	ImageConversationType    = 2
	AdditionConversationType = 3

	UserRole = "user"
	BotRole  = "bot"
)

// Mongo 数据库结构
type BotConversation struct {
	Role       string         `json:"role" bson:"role"` // user / bot
	Content    string         `json:"content" bson:"content"`
	Type       int64          `json:"type" bson:"type"` // 1 文本类型， 2 图片类型 3 附加消息
	Options    string         `json:"options" bson:"options"`
	CreateTime int64          `json:"create_time" bson:"create_time"`
	Conclusion string         `json:"conclusion" bson:"conclusion"`
	Additions  []AdditionInfo `json:"additions" bson:"additions"`
}

type AdditionInfo struct {
	DocumentName string `json:"document_name" bson:"document_name"`
	DocumentID   string `json:"document_id" bson:"document_id"`
	DatasetName  string `json:"dataset_name" bson:"dataset_name"`
	DatasetID    string `json:"dataset_id" bson:"dataset_id"`
	Position     int    `json:"position" bson:"position"`
	Content      string `json:"content" bson:"content"`
}

// Mongo 数据库结构
type BotChat struct {
	BotName          string             `json:"bot_name" bson:"bot_name"`                 // 机器名称
	NameCN           string             `json:"name_cn" bson:"name_cn"`                   // 机器中文名称
	Image            string             `json:"image" bson:"image"`                       // 图片
	LastUpdateTime   int64              `json:"last_update_time" bson:"last_update_time"` // 最后更新时间
	FirstQuestion    string             `json:"first_question" bson:"first_question"`     // 最开始的问题
	SessionID        string             `json:"session_id" bson:"session_id"`             // 对话唯一 ID
	Status           int32              `json:"status" bson:"status"`                     // 1 正常，2 已删除
	ConversationList []*BotConversation `json:"conversation_list" bson:"conversation_list"`
}

type GetChatConversationsResp struct {
	SessionID        string             `json:"session_id"`
	ConversationList []*BotConversation `json:"conversation_list"`
}

// 请求回复
type GetChatResp struct {
	ChatList []*BotChat `json:"chat_list"`
}

type CreateChatReq struct {
	BotName string `json:"bot_name" bson:"bot_name"` // 机器名称
}

type CreateChatResp struct {
	SessionID string `json:"session_id" bson:"session_id"`
}

type SaveChatConversationsReq struct {
	SessionID        string             `json:"session_id" bson:"session_id"`
	ConversationList []*BotConversation `json:"conversation_list" bson:"conversation_list"`
}

type SaveChatConversationsResp struct {
	BotName string `json:"bot_name" bson:"bot_name"`
}

type DeleteChatReq struct {
	SessionID string `json:"session_id" bson:"session_id"`
}
