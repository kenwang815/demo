package daos

import (
	"time"

	"github/demo/model"
	"github/demo/model/device"
	"github/demo/utils"
	"github/demo/utils/log"

	"github.com/jinzhu/gorm"
)

type deviceRepo struct {
	copy *gorm.DB
	db   *gorm.DB
}

func (r *deviceRepo) Get(id device.UUID) (*device.Device, error) {
	var d device.Device

	if err := r.db.Where("id = ?", id).Find(&d).Error; err != nil {
		log.Errorf("deviceRepository Get fail => %+v", err)
		return nil, err
	}

	return &d, nil
}

func (r *deviceRepo) Create(d *device.Device) (*device.Device, error) {
	d.CreateTime = time.Now().UnixNano() / int64(time.Millisecond)

	md := r.db.Create(d)
	if err := md.Error; err != nil {
		log.Errorf("deviceRepository Create fail => %+v", err)
		return nil, err
	}
	x := md.Value.(*device.Device)
	return x, nil
}

func (r *deviceRepo) Update(d *device.Device) (*device.Device, int64, error) {
	if d == nil {
		return nil, 0, nil
	}

	c := device.Device{
		Id: d.Id,
	}

	if *d == c {
		return nil, 0, nil
	}

	d.UpdateTime = time.Now().UnixNano() / int64(time.Millisecond)
	umap := utils.Map(d)

	re := &device.Device{}
	x := r.Query("id = ?", d.Id).Model(re).Updates(umap)
	if err := x.Error; err != nil {
		log.Errorf("[DB][device] update Info error: %+v", err)
		return nil, 0, err
	}
	affectRow := x.RowsAffected

	if affectRow <= 0 {
		return re, affectRow, x.Error
	}

	return re, affectRow, nil
}

func (r *deviceRepo) Delete(id device.UUID) (int64, error) {
	delete := r.db.Where("id = ?", id).Delete(&device.Device{})
	row := delete.RowsAffected
	var err error
	if err = delete.Error; err != nil {
		log.Errorf("deviceRepository Delete fail => %+v", err)
	}
	return row, err
}

func (r *deviceRepo) List(d *device.Device) ([]*device.Device, error) {
	var devices []*device.Device
	if err := r.db.Where(d).Find(&devices).Error; err != nil {
		log.Errorf("deviceRepository List fail => %+v", err)
		return nil, err
	}
	return devices, nil
}

func (r *deviceRepo) Find(d *device.Device, p *model.Page) ([]*device.Device, error) {
	var devices []*device.Device
	if err := r.db.Where(d).Limit(p.Limit).Offset(p.Offset).Find(&devices).Error; err != nil {
		log.Errorf("deviceRepository Find fail => %+v", err)
		return nil, err
	}
	return devices, nil
}

func (r *deviceRepo) Query(query interface{}, args ...interface{}) *gorm.DB {
	return r.db.Where(query, args...)
}

func (r *deviceRepo) NewTransactions() {
	r.db = r.db.Begin()
}

func (r *deviceRepo) TransactionsRollback() {
	r.db.Rollback()
	r.db = r.copy
}

func (r *deviceRepo) TransactionsCommit() {
	r.db.Commit()
	r.db = r.copy
}

func NewDeviceRepo(db *gorm.DB) device.Repository {
	return &deviceRepo{
		copy: db,
		db:   db,
	}
}
