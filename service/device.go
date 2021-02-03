package service

import (
	"time"

	"github.com/gofrs/uuid"

	"github/demo/model"
	"github/demo/model/device"
	"github/demo/utils/log"
)

type Device struct {
	Id         string `json:"id"`
	Model      string `json:"model"`
	Color      string `json:"color"`
	Version    string `json:"version"`
	CreateTime int64  `json:"create_time"`
	UpdateTime int64  `json:"update_time"`
}

func (d *Device) repoType() *device.Device {
	return &device.Device{
		Id:      device.UUID(d.Id),
		Model:   d.Model,
		Color:   d.Color,
		Version: d.Version,
	}
}

func (d *Device) Assemble(r *device.Device) {
	d.Id = r.Id.String()
	d.Model = r.Model
	d.Color = r.Color
	d.Version = r.Version
	d.CreateTime = r.CreateTime
	d.UpdateTime = r.UpdateTime
}

type deviceService struct {
	deviceRepo device.Repository
}

func (s *deviceService) Find(d *Device, page *Page) ([]*Device, ErrorCode) {
	if d == nil {
		return nil, ErrorCodeBadRequest
	}

	p := &model.Page{}
	if page != nil {
		p.Limit = page.Number
		p.Offset = page.Number * (page.Page - 1)
	}

	rows, err := s.deviceRepo.Find(d.repoType(), p)
	if err != nil {
		return nil, ErrorCodeDeviceDBFindFail
	}

	if len(rows) == 0 {
		return nil, ErrorCodeSuccessButNotFound
	}

	devices := []*Device{}
	for _, v := range rows {
		srv := &Device{}
		srv.Assemble(v)
		devices = append(devices, srv)
	}

	return devices, ErrorCodeSuccess
}

func (s *deviceService) Register(d *Device) (*Device, ErrorCode) {
	if d == nil {
		return nil, ErrorCodeBadRequest
	}

	now := time.Now().UnixNano() / int64(time.Millisecond)
	d.CreateTime = now
	d.UpdateTime = now
	d.Id = uuid.Must(uuid.NewV4()).String()
	x, err := s.deviceRepo.Create(d.repoType())
	if err != nil {
		return nil, ErrorCodeDeviceDBCreateFail
	}

	re := &Device{}
	re.Assemble(x)
	return re, ErrorCodeSuccess
}

func (s *deviceService) Update(d *Device) (int64, ErrorCode) {
	iformat := uuid.FromStringOrNil(d.Id)
	if iformat == uuid.Nil {
		return 0, ErrorCodeParseUUIDFail
	}

	_, a, err := s.deviceRepo.Update(d.repoType())
	if err != nil {
		return 0, ErrorCodeDeviceDBUpdateFail
	}

	return a, ErrorCodeSuccess
}

func (s *deviceService) Delete(i string) ErrorCode {
	iformat := uuid.FromStringOrNil(i)
	if iformat == uuid.Nil {
		return ErrorCodeParseUUIDFail
	}

	affect, err := s.deviceRepo.Delete(device.UUID(i))
	if err != nil {
		return ErrorCodeDeviceDBDeleteFail
	}
	if affect <= 0 {
		log.Error("File does not exist: ", i)
		return ErrorCodeSuccessButNotFound
	}

	return ErrorCodeSuccess
}

func NewDeviceService(dr device.Repository) IDeviceService {
	return &deviceService{
		deviceRepo: dr,
	}
}
