package model

type AssistantInfo struct {
	Id         string `json:"id" bson:"id"`
	Name       string `json:"name" bson:"name"`
	Avatar     string `json:"avatar" bson:"avatar"`
	Desc       string `json:"desc" bson:"desc"`
	CategoryID string `json:"category_id" bson:"category_id"`
	Type       string `json:"type" bson:"type"`
	CreateAt   int64  `json:"create_at" bson:"create_at"`

	// 计量助手
	Docs []*CacleDocInfo `json:"docs" bson:"docs"`

	// 问答助手
	Chats []*ChatInfo `json:"chats" bson:"chats"`
}

type ChatServerInfo struct {
	Id           string          `json:"id" bson:"id"`
	Name         string          `json:"name" bson:"name"`
	Conversation []*Conversation `json:"conversation" bson:"conversation"`
	CreateAt     int64           `json:"create_at" bson:"create_at"`
	AssistantID  string          `json:"assistant_id" bson:"assistant_id"`
}

type ChatInfo struct {
	Id           string          `json:"id" bson:"id"`
	Name         string          `json:"name" bson:"name"`
	Conversation []*Conversation `json:"conversation" bson:"conversation"`
	CreateAt     int64           `json:"create_at" bson:"create_at"`
	DeleteAt     int64           `json:"delete_at" bson:"delete_at"`
}

type ChatDone struct {
	BotName      string        `json:"botname" bson:"botname"`
	Conversation *Conversation `json:"conversation" bson:"conversation"`
	Details      string        `json:"details" bson:"details"`
}

type Conversation struct {
	Role            string        `json:"role,omitempty" bson:"role"`
	Content         string        `json:"content,omitempty" bson:"content"`
	FullPrompt      string        `json:"full_prompt,omitempty" bson:"full_prompt"`
	SearchResultStr string        `json:"search_result_str,omitempty" bson:"search_result_str"`
	SearchResult    []*SearchItem `json:"search_result,omitempty" bson:"search_result"`
	Time            string        `json:"time,omitempty" bson:"time"`
}

type SearchItem struct {
	KbName string      `json:"kbname" bson:"kbname"`
	Result *ResultInfo `json:"result" bson:"result"`
}

type ResultInfo struct {
	Id       [][]string `json:"id" bson:"id"`
	Question [][]string `json:"question" bson:"question"`
	Answer   [][]string `json:"answer" bson:"answer"`
}

type CacleDocInfo struct {
	Name     string `json:"name" bson:"name"`
	Id       string `json:"id" bson:"id"`
	FilePath string `json:"file_path" bson:"file_path"`
	Status   string `json:"status" bson:"status"` // 0: 正在处理 1: 通过 2: 未通过
	Result   string `json:"result" bson:"result"` // 0: 未通过，1: 通过
}

const (
	CacleDocInfoStatusOk = "ok"
	CacleDocInfoStatusNo = "no"
	CacleDocInfoStatusUn = "ing"
)

type CreateAssistantReq struct {
	Name       string `json:"name" bson:"name"`
	Avatar     string `json:"avatar" bson:"avatar"`
	CategoryID string `json:"category_id" bson:"category_id"`
	Desc       string `json:"desc" bson:"desc"`
	Type       string `json:"type" bson:"type"`
}

type CheckAssistantDocsReq struct {
	DocID       string `json:"doc_id" bson:"doc_id"`
	DocName     string `json:"doc_name" bson:"doc_name"`
	DocFilePath string `json:"doc_path" bson:"doc_path"`
	AssistantID string `json:"assistant_id" bson:"assistant_id"`
}

type AsyncCheckAssistantDocsReq struct {
	DocID       string `json:"doc_id" bson:"doc_id"`
	DocName     string `json:"doc_name" bson:"doc_name"`
	DocFilePath string `json:"doc_path" bson:"doc_path"`
	AssistantID string `json:"assistant_id" bson:"assistant_id"`
	Template    string `json:"template" bson:"template"`
}

type CheckAssistantDocsResp struct {
	Status string `json:"status"`
	DocID  string `json:"doc_id"`
}

type ApiCheckDocReq struct {
	DocID    string `json:"docid"`
	DocPath  string `json:"docpath"`
	Template string `json:"template"`
}

type ApiCheckDocResp struct {
	Status string `json:"status"`
}

type AsyncCheckDoc struct {
	StatusCode int    `json:"status_code"`
	Docid      string `json:"docid"`
	Content    string `json:"content"`
}

type StreamChatResp struct {
	StatusCode int    `json:"status_code"`
	BotName    string `json:"botname"`
	Content    string `json:"content"`
}

type UploadAssistantDocsReq struct {
	AssistantID string `json:"assistant_id" bson:"assistant_id"`
	DocName     string `json:"doc_name" bson:"doc_name"`
	DocFilePath string `json:"doc_path" bson:"doc_path"`
}

type AssistantUpdateDocsStatusReq struct {
	AssistantID string `json:"assistant_id" bson:"assistant_id"`
	DocID       string `json:"doc_id" bson:"doc_id"`
}

type UploadAssistantDocsResp struct {
	DocID string `json:"doc_id"`
}

type AssistantChatListReq struct {
	AssistantID string `json:"assistant_id" bson:"assistant_id"`
}

type AssistantChatListResp struct {
	Chats []*ChatInfo `json:"chats"`
}

type AssistantChatSendReq struct {
	AssistantID string `json:"assistant_id" bson:"assistant_id"`
	ChatID      string `json:"chat_id" bson:"chat_id"`
	Content     string `json:"content" bson:"content"`
}

type ApiSendContentToAIReq struct {
	BotName      string          `json:"botname"`
	Conversation []*Conversation `json:"conversation"`
	Details      string          `json:"details"`
}

type ApiSendContentToAIResp struct {
	BotName      string          `json:"botname"`
	Conversation []*Conversation `json:"conversation"`
	Details      string          `json:"details"`
}

type AssistantChatCreateReq struct {
	AssistantID string `json:"assistant_id" bson:"assistant_id"`
	ChatName    string `json:"chat_name" bson:"chat_name"`
}

type AssistantChatCreateResp struct {
	ChatID string `json:"chat_id"`
}

type AssistantChatDeleteReq struct {
	AssistantID string `json:"assistant_id" bson:"assistant_id"`
	ChatID      string `json:"chat_id" bson:"chat_id"`
}

type AssistantChatDeleteResp struct {
	ChatID string `json:"chat_id"`
}

type AssistantChatUpdateReq struct {
	AssistantID string `json:"assistant_id" bson:"assistant_id"`
	ChatID      string `json:"chat_id" bson:"chat_id"`
	ChatName    string `json:"chat_name" bson:"chat_name"`
}

type AssistantChatSaveConversationReq struct {
	AssistantID  string          `json:"assistant_id" bson:"assistant_id"`
	ChatID       string          `json:"chat_id" bson:"chat_id"`
	Conversation []*Conversation `json:"conversation"`
}

type CheckAssistantDocResultReq struct {
	AssistantID string `json:"assistant_id" bson:"assistant_id"`
	DocID       string `json:"doc_id" bson:"doc_id"`
}

type CheckAssistantDocClearReq struct {
	AssistantID string `json:"assistant_id" bson:"assistant_id"`
}

type CheckAssistantDocResultResp struct {
	DocID       string         `json:"docid"`
	CheckResult []*CheckResult `json:"results"`
}

type CheckResult struct {
	Doc_check_uuid string            `json:"doc_check_uuid"`
	Id             int               `json:"id"`
	Check_type     int               `json:"check_type"`
	Reason         string            `json:"reason"`
	Created_at     string            `json:"created_at"`
	Check_task     string            `json:"check_task"`
	Checkername    string            `json:"checkername"`
	Result         int               `json:"result"`
	Doc_page       int               `json:"doc_page"`
	Last_modified  string            `json:"last_modified"`
	Standard       *UtilsFunStandard `json:"standard"`
	Ledger         *UtilsFunLedger   `json:"ledger"`
}

type UtilsFunStandard struct {
	Standard_id   int    `json:"standard_id"`
	Standard_path string `json:"standard_path"`
	Standard_page int    `json:"standard_page"`
}

type UtilsFunLedger struct {
	Ledger_id    int    `json:"ledger_id"`
	Ledger_value string `json:"ledger_value"`
}
