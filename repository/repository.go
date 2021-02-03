package repository

import (
	"github.com/jinzhu/gorm"

	"github/demo/config"
	"github/demo/database"
	"github/demo/utils/log"
)

type Engine struct {
	Database database.IDatabase
	GormDB   *gorm.DB
}

func NewEngine(c *config.Config) (*Engine, error) {
	db, err := database.NewDatabase(c.Database)
	if err != nil {
		return nil, err
	}

	e := &Engine{
		Database: db,
		GormDB:   db.GetDB(),
	}

	log.Info("Create engine success")
	return e, nil
}
