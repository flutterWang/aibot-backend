package bot_io

import (
	"aibot-backend/model"
	"sort"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
)

func (io *Io) GetDeviceList(c *gin.Context) ([]*model.DeviceInfo, error) {
	var deviceInfos []*model.DeviceInfo

	result, err := io.MongoClient.Collection("device").Find(c, bson.D{})
	if err != nil {
		return nil, err
	}

	err = result.All(c, &deviceInfos)
	if err != nil {
		return nil, err
	}

	// 把 status == "online" 的排在最前面
	customSort(deviceInfos)

	return deviceInfos, err
}

func customSort(data []*model.DeviceInfo) {
	sort.Slice(data, func(i, j int) bool {
		order := map[string]int{
			"online":     1,
			"offline":    2,
			"brokendown": 3,
		}

		return order[data[i].State] < order[data[j].State]
	})
}

func (io *Io) GetDeviceInfo(c *gin.Context, name string) (*model.DeviceInfo, error) {
	mongoClient := io.GetMongoClient()

	var deviceInfo model.DeviceInfo
	err := mongoClient.Collection("device").FindOne(c, bson.M{"name": name}).Decode(&deviceInfo)
	if err != nil {
		return nil, err
	}

	return &deviceInfo, nil
}

func (s *Io) CreateDevice(c *gin.Context, req *model.DeviceInfo) error {
	_, err := s.MongoClient.Collection("device").InsertOne(c, req)
	if err != nil {
		return err
	}

	return nil
}

func (s *Io) CreateDevices(c *gin.Context, devices []*model.DeviceInfo) error {
	datas := make([]interface{}, len(devices))
	for i := range devices {
		datas[i] = devices[i]
	}
	_, err := s.MongoClient.Collection("device").InsertMany(c, datas)
	if err != nil {
		return err
	}

	return nil
}

func (io *Io) UpdateDevice(c *gin.Context, device *model.DeviceInfo) error {
	mongoClient := io.GetMongoClient()
	_, err := mongoClient.Collection("device").UpdateOne(c, bson.M{"id": device.ID}, bson.M{"$set": device})
	if err != nil {
		return err
	}
	return nil
}
