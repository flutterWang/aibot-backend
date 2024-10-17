package model

type DeviceInfo struct {
	ID               int         `json:"id" bson:"id"`
	Name             string      `json:"name" bson:"name"`
	NameCN           string      `json:"name_cn" bson:"name_cn"`
	Model            string      `json:"model" bson:"model"`
	Manufacturer     string      `json:"manufacturer" bson:"manufacturer"`
	Category         string      `json:"category" bson:"category"`
	Workshop         string      `json:"workshop" bson:"workshop"`
	InstallationDate string      `json:"installation_date" bson:"installation_date"`
	Image            string      `json:"image" bson:"image"`
	State            string      `json:"state" bson:"state"`
	Collected        bool        `json:"collected" bson:"collected"`
	Description      string      `json:"description" bson:"description"`
	FileInfos        []*FileInfo `json:"file_infos" bson:"file_infos"`
}

type FileInfo struct {
	Url  string `json:"url" bson:"url"`
	Name string `json:"name" bson:"name"`
	Size string `json:"size" bson:"size"`
}

type DeviceCollectReq struct {
	Name      string `json:"name" bson:"name"`
	Collected bool   `json:"collected" bson:"collected"`
}

type CreateDevicesReq struct {
	Devices []*DeviceInfo `json:"devices"`
}

type UpdateDevicesReq struct {
	Device *DeviceInfo `json:"device"`
}
