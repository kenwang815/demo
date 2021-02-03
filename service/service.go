package service

import (
	"github/demo/config"
	"github/demo/daos"
	"github/demo/model/device"
	"github/demo/repository"
	"github/demo/utils/log"
)

var (
	// === Repository ===
	DeviceRepo device.Repository

	// === Service ===
	DeviceService IDeviceService
)

func Init(cf *config.Config, engine *repository.Engine) error {
	// === Repository ===
	DeviceRepo = daos.NewDeviceRepo(engine.GormDB)

	// === Service ===
	DeviceService = NewDeviceService(DeviceRepo)

	log.Info("Create service success")
	return nil
}
