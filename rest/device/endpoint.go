package device

import (
	"github.com/gin-gonic/gin"

	"github/demo/rest/content"
	"github/demo/service"
)

type Device struct {
	Id         string `form:"id"`
	Model      string `form:"model"`
	Color      string `form:"color"`
	Version    string `form:"version"`
	CreateTime int64  `form:"create_time"`
	UpdateTime int64  `form:"update_time"`
}

func (d *Device) serviceType() *service.Device {
	return &service.Device{
		Id:      d.Id,
		Model:   d.Model,
		Color:   d.Color,
		Version: d.Version,
	}
}

func (d *Device) Assemble(s *service.Device) {
	d.Id = s.Id
	d.Model = s.Model
	d.Color = s.Color
	d.Version = s.Version
	d.CreateTime = s.CreateTime
	d.UpdateTime = s.UpdateTime
}

func FindDevice(c *gin.Context) {
	page := &service.Page{}
	c.ShouldBind(page)
	device := &Device{}
	c.ShouldBind(device)

	rows, code := service.DeviceService.Find(device.serviceType(), page)
	resp := content.NewContent()
	if code == service.ErrorCodeSuccess {
		data := make(map[string]interface{})
		if page != nil {
			data["page"] = page
		}
		var re []*Device
		for _, v := range rows {
			srv := &Device{}
			srv.Assemble(v)
			re = append(re, srv)
		}
		data["datas"] = re
		resp.Data(data)
	}

	resp.Code(code.Int()).Msg(service.ErrorMsg(code))
	c.JSON(service.ErrorStatusCode(code), resp)
}

func RegisterDevice(c *gin.Context) {
	device := &Device{}
	c.ShouldBind(device)
	m, code := service.DeviceService.Register(device.serviceType())
	resp := content.NewContent()

	var re *Device
	if code == service.ErrorCodeSuccess {
		re = &Device{}
		re.Assemble(m)
	}
	resp.Data(re)

	resp.Code(code.Int()).Msg(service.ErrorMsg(code))
	c.JSON(service.ErrorStatusCode(code), resp)
}

func DeleteDevice(c *gin.Context) {
	d := c.Param("id")
	code := service.DeviceService.Delete(d)
	resp := content.NewContent()
	resp.Code(code.Int()).Msg(service.ErrorMsg(code))
	c.JSON(service.ErrorStatusCode(code), resp)
}

func UpdateDevice(c *gin.Context) {
	device := &Device{}
	c.ShouldBind(device)
	affect, code := service.DeviceService.Update(device.serviceType())
	resp := content.NewContent()
	m := map[string]int64{
		"affect": affect,
	}
	resp.Data(m)
	resp.Code(code.Int()).Msg(service.ErrorMsg(code))
	c.JSON(service.ErrorStatusCode(code), resp)
}
