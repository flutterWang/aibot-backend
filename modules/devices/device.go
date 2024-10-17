package devices

import (
	"aibot-backend/io"
	"aibot-backend/model"

	"github.com/gin-gonic/gin"
)

func GetDeviceList(c *gin.Context, io *bot_io.Io) ([]*model.DeviceInfo, error) {
	result, err := io.GetDeviceList(c)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func DeviceCollect(c *gin.Context, io *bot_io.Io, req *model.DeviceCollectReq) error {
	device, err := io.GetDeviceInfo(c, req.Name)
	if err != nil {
		return err
	}

	device.Collected = req.Collected

	err = io.UpdateDevice(c, device)
	if err != nil {
		return err
	}

	return nil
}

func CreateDevices(c *gin.Context, io *bot_io.Io, devices []*model.DeviceInfo) error {
	return io.CreateDevices(c, devices)
}

func UpdateDevice(c *gin.Context, io *bot_io.Io, device *model.DeviceInfo) error {
	return io.UpdateDevice(c, device)
}

func GetDevice(c *gin.Context, io *bot_io.Io, name string) (*model.DeviceInfo, error) {
	return io.GetDeviceInfo(c, name)
}
