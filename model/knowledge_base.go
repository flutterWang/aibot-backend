package model

type GetKbResp struct {
	Kbs []*Kb `json:"kbs" bson:"kbs"`
}

type Kb struct {
	Id         string     `json:"id" bson:"id"`
	Name       string     `json:"name" bson:"name,omitempty"`
	Cover      string     `json:"cover" bson:"cover,omitempty"`
	Desc       string     `json:"desc" bson:"desc,omitempty"`
	CreateBy   *CreateBy  `json:"create_by" bson:"create_by,omitempty"`
	CreateAt   int64      `json:"create_at" bson:"create_at,omitempty"`
	UpdateAt   int64      `json:"update_at" bson:"update_at,omitempty"`
	CategoryID string     `json:"category_id" bson:"category_id,omitempty"`
	Labels     []string   `json:"labels" bson:"labels,omitempty"`
	DeleteAt   int64      `json:"delete_at" bson:"delete_at,omitempty"`
	Docs       []*DocInfo `json:"docs" bson:"docs,omitempty"`
}

type CreateBy struct {
	Id       int64  `json:"id" bson:"id"`
	UserName string `json:"username" bson:"user_name"`
}

type DocInfo struct {
	Name string `json:"name" bson:"name"`
	Url  string `json:"url" bson:"url"`
}

type CategoryInfo struct {
	Id   string `json:"id" bson:"id"`
	Name string `json:"name" bson:"name"`
}

type Label struct {
	Id   string `json:"id" bson:"id"`
	Name string `json:"name" bson:"name"`
}

type CreateKbReq struct {
	Name       string   `json:"name" bson:"name"`
	CategoryID string   `json:"category_id" bson:"category_id"`
	Labels     []string `json:"labels" bson:"labels"`
	CoverPath  string   `json:"cover_path" bson:"cover_path"`
	Desc       string   `json:"desc" bson:"desc"`
}

type CreateKbResp struct {
	Id string `json:"id" bson:"id"`
}

type GetKbDetailResp struct {
	Kb *Kb `json:"kb" bson:"kb"`
}

type UpdateKbReq struct {
	ID         string   `json:"id" bson:"id"`
	Name       string   `json:"name" bson:"name"`
	CategoryID string   `json:"category_id" bson:"category_id"`
	Labels     []string `json:"labels" bson:"labels"`
	CoverPath  string   `json:"cover_path" bson:"cover_path"`
	Desc       string   `json:"desc" bson:"desc"`
}

type GetKbLabelsListResp struct {
	Labels []*Label `json:"labels" bson:"labels"`
}

type CreateLabelReq struct {
	Name string `json:"name" bson:"name"`
}

type CreateLabelResp struct {
	Id string `json:"id"`
}

type GetLabelKbsResp struct {
	Kbs []*Kb `json:"kbs" bson:"kbs"`
}

type SearchKbReq struct {
	Keyword string `json:"keyword"`
}

type SearchKbResp struct {
	Kbs []*Kb `json:"kbs" bson:"kbs"`
}

type GetCategoryListResp struct {
	Categories []*CategoryInfo `json:"categories" bson:"categories"`
}

type CreateCategoryReq struct {
	Name string `json:"name" bson:"name"`
}

type CreateCategoryResp struct {
	Id string `json:"id" bson:"id"`
}

type GetCategoryKbsResp struct {
	Kbs []*Kb `json:"kbs" bson:"kbs"`
}

type UploadDocReq struct {
	Doc  *DocInfo `json:"doc" bson:"doc"`
	KbID string   `json:"kb_id" bson:"kb_id"`
}
