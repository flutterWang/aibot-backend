package model

type AsyncCheckContent struct {
	DocID       string `json:"doc_id" bson:"doc_id"`
	DocName     string `json:"doc_name" bson:"doc_name"`
	DocFilePath string `json:"doc_path" bson:"doc_path"`
	AssistantID string `json:"assistant_id" bson:"assistant_id"`
	Template    string `json:"template" bson:"template"`
}

type ActionToBack struct {
	Action  string `json:"action"`
	Message []byte `json:"message"`
}

type BackToFront struct {
	Action  string `json:"action"`
	Message string `json:"message"`
}
