package device

import (
	"github/demo/model"

	"github.com/jinzhu/gorm"
)

type UUID string

func (u UUID) String() string { return string(u) }

type Device struct {
	Id         UUID   `gorm:"column:id;unique;type:uuid;primary_key" mapKey:"ignore"`
	Model      string `gorm:"column:model;not null" mapKey:"model,omitempty"`
	Color      string `gorm:"column:color;not null" mapKey:"color,omitempty"`
	Version    string `gorm:"column:version;not null" mapKey:"version,omitempty"`
	CreateTime int64  `gorm:"column:create_time;not null" mapKey:"ignore"`
	UpdateTime int64  `gorm:"column:update_time;not null" mapKey:"update_time"`
}

func (Device) TableName() string {
	return "device"
}

type Repository interface {
	Get(id UUID) (*Device, error)
	Create(d *Device) (*Device, error)
	Update(d *Device) (*Device, int64, error)
	Delete(id UUID) (int64, error)
	List(d *Device) ([]*Device, error)
	Find(d *Device, p *model.Page) ([]*Device, error)
	Query(query interface{}, args ...interface{}) *gorm.DB
}
