package api

import (
	"aibot-backend/model"
	"aibot-backend/modules/devices"

	"github.com/gin-gonic/gin"
)

func (s *server) handlerDeviceList(c *gin.Context) {
	resp, err := devices.GetDeviceList(c, s.io)
	if err != nil {
		ginAbortWithCode(c, 500, err)
		return
	}
	ginJson(c, resp)
}

func (s *server) handlerDeviceCollect(c *gin.Context) {
	req := &model.DeviceCollectReq{}
	if err := c.ShouldBindJSON(req); err != nil {
		ginAbortWithCode(c, 400, err)
		return
	}

	err := devices.DeviceCollect(c, s.io, req)
	if err != nil {
		ginAbortWithCode(c, 500, err)
		return
	}

	ginJson(c, nil)
}

func (s *server) handlerCreateDevices(c *gin.Context) {
	req := &model.CreateDevicesReq{}
	if err := c.ShouldBindJSON(req); err != nil {
		ginAbortWithCode(c, 400, err)
		return
	}

	err := devices.CreateDevices(c, s.io, req.Devices)
	if err != nil {
		ginAbortWithCode(c, 500, err)
		return
	}

	ginJson(c, nil)
}

func (s *server) handlerGetDeviceInfo(c *gin.Context) {
	name := c.Query("name")
	if name == "" {
		ginAbortWithCode(c, 400, nil)
		return
	}

	result, err := devices.GetDevice(c, s.io, name)
	if err != nil {
		ginAbortWithCode(c, 500, err)
		return
	}

	ginJson(c, result)
}

func (s *server) handlerUpdateDevice(c *gin.Context) {
	req := &model.UpdateDevicesReq{}
	if err := c.ShouldBindJSON(req); err != nil {
		ginAbortWithCode(c, 400, err)
		return
	}

	err := devices.UpdateDevice(c, s.io, req.Device)
	if err != nil {
		ginAbortWithCode(c, 500, err)
		return
	}

	ginJson(c, nil)
}
